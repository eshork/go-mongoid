package mongoid

import (
	"context"
	mongoidError "mongoid/errors"
	"mongoid/log"

	"go.mongodb.org/mongo-driver/mongo"
	// "strconv"
	"runtime"

	"go.mongodb.org/mongo-driver/bson"
)

// Result provides convenient access to the data produced by database operations
type Result struct {
	cursor  *mongo.Cursor   // the mongo driver cursor for the query
	context context.Context // context to pass to any future driver calls
	model   *ModelType      // the ModelType associated with the query

	// First() IDocumentBase
	// Last() IDocumentBase
	// Count() int
	// At(index int) IDocumentBase
	//
	// ToA() error
	// ForEach(f func () error) error
	lookback    []bson.M // cache of records, to support random access via At(), First(), Last(), etc
	cursorIndex uint     // the current index of the driver cursor
	streaming   bool     // track streaming access
	closed      bool     // track closed state
}

func makeResult(ctx context.Context, cursor *mongo.Cursor, model *ModelType) *Result {
	// finalizer to ensure that given cursor is properly Close()'d, even if caller later forgets
	runtime.SetFinalizer(cursor, func(c *mongo.Cursor) {
		c.Close(context.Background())
	})

	return &Result{
		cursor:      cursor,
		context:     ctx,
		model:       model,
		lookback:    make([]bson.M, 0),
		cursorIndex: 0,
		streaming:   false,
		closed:      false,
	}
}

// Close the Result, indicating you are done accessing the object and are ready to free the related resources.
func (res *Result) Close() {
	if !res.closed {
		res.cursor.Close(res.context)
		res.closed = true
	}
}

// Streaming disables the lookback cache of the Result, which can be useful whenever memory usage is a concern (such as working with very large result sets).
// Without a lookback cache, the Result is only readable once via either Result.ForEach() or Result.ToAry().
// Attempting to read from a Streaming Result more than once will Panic - (TODO maybe we should change this into a re-execution of the query).
// The call to enable Streaming() should be the first operation performed on the Result, otherwise the behavior is undefined.
// Calling certain methods on a Result after Streaming is declared can result in a panic -- affected methods indicate such within their documentation.
func (res *Result) Streaming() *Result {
	if !res.streaming && res.cursorIndex > 0 {
		log.Panic(&mongoidError.InvalidOperation{
			MethodName: "Result.Streaming",
			Reason:     "Cannot enable Streaming after random access",
		})
	}
	res.streaming = true
	return res
}

// IsStreaming returns true after Streaming has been enabled on the Result
func (res *Result) IsStreaming() bool {
	return res.streaming
}

// First returns an interface to the first document in the Result set.
// This method will panic if Streaming() was enabled.
func (res *Result) First() IDocumentBase {
	log.Debug("Result.First()")
	if res.streaming {
		log.Panic(&mongoidError.InvalidOperation{
			MethodName: "Result.First",
			Reason:     "Cannot perform random access when Result.IsStreaming()",
		})
	}
	return res.at(0)
}

// Last returns an interface to the last document in the result set
// This method will panic if Streaming() was enabled.
func (res *Result) Last() IDocumentBase {
	log.Debug("Result.Last()")
	if res.streaming {
		log.Panic(&mongoidError.InvalidOperation{
			MethodName: "Result.Last",
			Reason:     "Cannot perform random access when Result.IsStreaming()",
		})
	}
	return res.at(res.Count() - 1) // this is the simplest implementation - res.Count will cause the full result to be read from the driver
}

// At returns an interface to the Document in the result set at the given index (range is 0 to count-1)
// This method will panic if Streaming() was enabled.
func (res *Result) At(index uint) IDocumentBase {
	log.Debug("Result.At()")
	if res.streaming {
		log.Panic(&mongoidError.InvalidOperation{
			MethodName: "Result.At",
			Reason:     "Cannot perform random access when Result.IsStreaming()",
		})
	}
	return res.at(index)
}

// retrieves the record at the given index, reading additional records from the db driver as needed
func (res *Result) at(index uint) IDocumentBase {
	// read responses until we have the record we need in lookback cache
	for res.cursorIndex <= index {
		more := res.readNextToLookback()
		if !more {
			log.Panic(&mongoidError.IndexOutOfBounds{})
		}
	}
	result := res.lookback[index]
	retAsIDocumentBase := makeDocument(res.model, result)
	return retAsIDocumentBase
}

// Count returns the number of documents in the Result.
// This method will panic if Streaming() was enabled.
// The current implementation will read all remaining result items into memory, making this a poor choice for queries with large result sets.
// TODO - Make a ModelType.Count(filter_query) to perform server-side document count queries
func (res *Result) Count() uint {
	log.Debug("Result.Count()")
	if res.streaming {
		log.Panic(&mongoidError.InvalidOperation{
			MethodName: "Result.Count",
			Reason:     "Cannot perform random access when Result.IsStreaming()",
		})
	}

	for res.readNextToLookback() {
		// loop until there's nothing left to read into lookback
	}

	// then the count is simply the length of the lookback
	return uint(len(res.lookback))
}

// OneAndClose is a convenience function to retrieve a single document and then Close the Result, as a combined step.
// This is most useful for queries that only yield a single record.
// Ex. `myDocObj := MyDocuments.Find_("myDocumentID").OneAndClose().(*MyDocument)`
func (res *Result) OneAndClose() IDocumentBase {
	log.Debug("Result.OneAndClose()")
	if res.streaming {
		log.Panic(&mongoidError.InvalidOperation{
			MethodName: "Result.OneAndClose",
			Reason:     "Cannot perform random access when Result.IsStreaming()",
		})
	}

	// For the moment, we can use this as the defacto method for fetching data within tests.
	// We'll want to reimplement this later, but for the moment it serves as a great place to encapsulate a bunch of single-item fetch logic.

	defer res.cursor.Close(res.context)
	for res.cursor.Next(res.context) {
		var result bson.M
		err := res.cursor.Decode(&result)
		if err != nil {
			log.Panic(err)
		}
		retAsIDocumentBase := makeDocument(res.model, result)
		return retAsIDocumentBase
	}
	if err := res.cursor.Err(); err != nil {
		log.Panic(err)
	}
	log.Panic(&mongoidError.NotFound{})
	return nil
}

// read the next result from the query cursor, and append it to the lookback cache
func (res *Result) readNextToLookback() bool {
	log.Trace("Result.readNextToLookback()")
	more := res.cursor.Next(res.context)
	// check for driver errors
	if err := res.cursor.Err(); err != nil {
		log.Panic(err)
	}
	if more { // process a new record if we found one
		var result bson.M
		err := res.cursor.Decode(&result)
		if err != nil {
			log.Panic(err)
		}
		res.lookback = append(res.lookback, result)
		res.cursorIndex++
	} else { // self close when we know there is no more data
		res.Close()
	}
	return more
}
