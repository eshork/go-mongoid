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
	// Count() uint
	// At(index int) IDocumentBase

	// ToAry() error
	// ForEach(f func (v IDocumentBase) error) error

	lookback    []bson.M // cache of records, to support random access via At(), First(), Last(), etc
	cursorIndex uint     // the current index of the driver cursor, what the next read will yield
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

// readNext reads the next result from the query cursor, and writes it into the bson.M at the given pointer.
// This will advance res.cursorIndex, but does not save the result within lookback cache.
// This method will panic if the Result is not Streaming().
func (res *Result) readNext(v *bson.M) bool {
	log.Trace("Result.readNext()")
	if !res.streaming {
		log.Panic(&mongoidError.InvalidOperation{
			MethodName: "Result.readNext",
			Reason:     "Expected Result.IsStreaming()",
		})
	}
	more := res.cursor.Next(res.context)
	// check for driver errors
	if err := res.cursor.Err(); err != nil {
		log.Panic(err)
	}
	if more { // process a new record if we found one
		err := res.cursor.Decode(v)
		if err != nil {
			log.Panic(err)
		}
		res.cursorIndex++
	}
	return more
}

// ForEach will call the given function once for each result, in the order they were returned by the server.
// The given function "fn" should accept an IDocumentBase as the only parameter.
// The given function "fn" may return an non-nil error value to halt further iterations - that value will be passed upward and returned by ForEach.
//
// Example:
//   ret := myResult.ForEach(func(v IDocumentBase) error {
//   	// do something with v
//   	// change it and save it, just read it, or delete it
//   	return nil // or return some error value
//   })
//
func (res *Result) ForEach(fn func(IDocumentBase) error) error {
	// the heavy lifting is within ForEachBson
	return res.ForEachBson(func(v bson.M) error {
		asIDocumentBase := makeDocument(res.model, v)
		return fn(asIDocumentBase)
	})
}

// ForEachBson is similar to ForEach, but provides the raw bson.M instead of an IDocumentBase
func (res *Result) ForEachBson(fn func(bson.M) error) error {
	if !res.streaming { // non-streaming implementation (records are stored to lookback cache as they are read)
		for i := uint(0); res.readNextToLookback(); i++ {
			// loop until there's nothing left to read into lookback
			result := res.lookback[i] // fetch current value from lookback
			r := fn(result)           // run the given fn
			// if fn had a non-nil return, then we should stop and bubble that value upward
			if r != nil {
				return r
			}
		}
		return nil
	}
	// streaming implementation (records are not stored to lookback cache)
	more := true
	for more {
		var result bson.M
		more = res.readNext(&result)
		if more {
			r := fn(result) // run the given fn
			// if fn had a non-nil return, then we should stop and bubble that value upward
			if r != nil {
				return r
			}
		}
	}
	return nil
}

// ToAry returns the results as a slice of []IDocumentBase
func (res *Result) ToAry() []IDocumentBase {
	resultAry := make([]IDocumentBase, 0)
	res.ForEach(func(v IDocumentBase) error {
		resultAry = append(resultAry, v)
		return nil
	})
	return resultAry
}

// ToAryBson returns the results as a slice of []bson.M
func (res *Result) ToAryBson() []bson.M {
	return nil
}
