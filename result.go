package mongoid

import (
	"context"
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
	lookback []bson.M
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

// Streaming disables the lookback cache of the Result, which has some access implications.
// Without a lookback cache, the Result is only readable once, from beginning to end via Result.ForEach(), or by Result.ToAry().
// Disabling the lookback cache can be useful for very large result sets, when memory consumption becomes a concern.
// Attempting to read from a Streaming Result more than once will Panic.
// Calling certain methods on a Result after Streaming is declared can result in a Panic -- affected methods indicate such within their documentation.
func (res *Result) Streaming() *Result {
	log.Panicf("NYI - Result.Streaming()")
	return res
}

// First returns an interface to the first document in the Result set
func (res *Result) First() IDocumentBase {
	log.Panicf("NYI - Result.First()")
	return nil
}

// Last returns an interface to the last document in the result set
func (res *Result) Last() IDocumentBase {
	log.Panicf("NYI - Result.Last()")
	return nil
}

// Count returns the number of document in the Result
func (res *Result) Count() int {
	log.Panicf("NYI - Result.Count()")
	return 0
}

// At returns an interface to the Document in the result set at the given index (range is 0 to count-1)
func (res *Result) At(index int) IDocumentBase {
	log.Panicf("NYI - Result.Last()")
	return nil
}

// OneAndClose is a convenience function to retrieve a single document and then Close the Result, as a combined step.
// This is most useful for queries that only yield a single record.
// Ex. `myDocObj := MyDocuments.Find_("myDocumentID").OneAndClose().(*MyDocument)`
func (res *Result) OneAndClose() IDocumentBase {
	log.Debug("Result.OneAndClose()")
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
