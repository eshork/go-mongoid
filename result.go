package mongoid

import (
	"context"
	mongoidError "mongoid/errors"
	"mongoid/log"

	"go.mongodb.org/mongo-driver/mongo"
	// "strconv"
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
	lookback  []bson.M
	streaming bool // track streaming access
	adhoc     bool // track adhoc access
}

func makeResult(ctx context.Context, cursor *mongo.Cursor, model *ModelType) *Result {
	return &Result{
		context: ctx,
		cursor:  cursor,
		model:   model,
	}
}

// Close the Result, indicating you are done accessing the object and are ready to free the related resources.
func (res *Result) Close() {
	res.cursor.Close(res.context)
}

// Streaming disables the lookback cache of the Result, which has some performance and access method implications.
// Without a lookback cache, the Result is only readable once, from beginning to end via Result.ForEach(), or by Result.ToAry().
// Disabling the lookback cache can be useful for very large result sets, when memory consumption becomes a concern.
// Attempting to read from a Streaming Result more than once will Panic - (TODO maybe we should change this into a re-execution of the query).
// The call to enable Streaming() should be the first operation performed on the Result, otherwise the behavior is undefined.
// Calling certain methods on a Result after Streaming is declared can result in a panic -- affected methods indicate such within their documentation.
func (res *Result) Streaming() *Result {
	if res.adhoc {
		panic(mongoidError.InvalidMethodCall{
			MethodName: "Result.Streaming()",
			Reason:     "Cannot enable Streaming after an adhoc access was performed",
		})
	}
	res.streaming = true
	return res
}

// First returns an interface to the first document in the Result set.
// This method will panic if Streaming() was enabled.
func (res *Result) First() IDocumentBase {
	if res.streaming {
		panic(mongoidError.InvalidMethodCall{
			MethodName: "Result.First()",
			Reason:     "Cannot perform adhoc access when Streaming()",
		})
	}
	// res.adhoc = true // performed within Result.at
	return res.at(0)
}

// Last returns an interface to the last document in the result set
// This method will panic if Streaming() was enabled.
func (res *Result) Last() IDocumentBase {
	if res.streaming {
		panic(mongoidError.InvalidMethodCall{
			MethodName: "Result.Last()",
			Reason:     "Cannot perform adhoc access when Streaming()",
		})
	}
	// res.adhoc = true // performed within Result.at
	return res.at(res.Count() - 1) // this is the simplest implementation - res.Count will cause the full result to be read from the driver
}

// At returns an interface to the Document in the result set at the given index (range is 0 to count-1)
// This method will panic if Streaming() was enabled.
func (res *Result) At(index int) IDocumentBase {
	log.Debug("Result.At()")
	if res.streaming {
		panic(mongoidError.InvalidMethodCall{
			MethodName: "Result.At()",
			Reason:     "Cannot perform adhoc access when Streaming()",
		})
	}
	// res.adhoc = true // performed within Result.at
	return res.at(index)
}

func (res *Result) at(index int) IDocumentBase {
	log.Trace("Result.at()")
	res.adhoc = true
	log.Panicf("NYI - Result.at()")
	return nil
}

// Count returns the number of documents in the Result.
// This method will panic if Streaming() was enabled.
// To perform this action, we must read all result items into memory (to support future adhoc access), making this a poor choice for queries with large result sets.
// TODO - Make ModelType.Count(filter_query) to perform server-side document count queries
func (res *Result) Count() int {
	if res.streaming {
		panic(mongoidError.InvalidMethodCall{
			MethodName: "Result.Count()",
			Reason:     "Cannot perform adhoc access when Streaming()",
		})
	}
	res.adhoc = true

	for res.readNextToLookback() {
		// loop until there's nothing left to read into lookback
	}

	return len(res.lookback)
}

// OneAndClose is a convenience function to retrieve a single document and then Close the Result, as a combined step.
// This is most useful for queries that only yield a single record.
// Ex. `myDocObj := MyDocuments.Find_("myDocumentID").OneAndClose().(*MyDocument)`
func (res *Result) OneAndClose() IDocumentBase {
	log.Debug("Result.OneAndClose()")
	if res.streaming {
		panic(mongoidError.InvalidMethodCall{
			MethodName: "Result.At()",
			Reason:     "Cannot perform adhoc access when Streaming()",
		})
	}
	res.adhoc = true

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

	log.Panic("NYI - Result.OneAndClose()") // <- this should be a panic for ItemNotFound or EmptyResultSet or similar
	return nil
}

func (res *Result) readNextToLookback() bool {
	return false
}
