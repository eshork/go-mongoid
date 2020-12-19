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

// Result provides access to the response data produced by database query operations.
// Result initially expects to be accessed randomly, and supports multiple access.
// To support unpredictable access patterns, many access operations read the entire result set from the Mongo driver into memory before returning.
// To avoid aggressive caching behavior, call the Streaming() method to declare the Result is for one-time sequential only access.
type Result struct {
	cursor  *mongo.Cursor   // the mongo driver cursor for the query
	context context.Context // context to pass to any future driver calls
	model   *ModelType      // the ModelType associated with the query

	// First() IDocumentBase
	// Last() IDocumentBase
	// Count() uint
	// At(index int) IDocumentBase

	// ToAry() []IDocumentBase
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
		context:     ctx, // yes, this is an anti-pattern
		model:       model,
		lookback:    make([]bson.M, 0),
		cursorIndex: 0,
		streaming:   false,
		closed:      false,
	}
}

func (res *Result) close() {
	if !res.closed {
		res.closed = true
		res.cursor.Close(res.context)
		res.cursor = nil
	}
}

// Streaming disables the lookback cache of the Result, which can be useful whenever memory usage is a concern (such as working with very large result sets).
// Without a lookback cache, the Result is only readable once via either Result.ForEach() or Result.ToAry().
// Attempting to read from a Streaming Result more than once will Panic - (TODO maybe we should change this into a re-execution of the query).
// The call to invoke Streaming should be the first operation performed on the Result, otherwise the behavior is undefined.
// Streaming returns a pointer back to the originating Result struct, so that it may be included within a naturally reading method call chain.
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

// First returns an interface to the first document in the Result set, or nil if the Result contains no records.
// This method will panic if Streaming() was enabled.
func (res *Result) First() IDocumentBase {
	log.Debug("Result.First()")
	if res.streaming {
		log.Panic(&mongoidError.InvalidOperation{
			MethodName: "Result.First",
			Reason:     "Cannot perform random access when Result.IsStreaming()",
		})
	}
	res.Count() // cheap way to load all results and close the db cursor
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
	return res.at(res.Count() - 1) // res.Count will cause the full result set to be read from the driver
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
	res.Count() // cheap way to load all results and close the db cursor
	return res.at(index)
}

// retrieves the record at the given index, reading additional records from the db driver as needed
func (res *Result) at(index uint) IDocumentBase {
	result := res.atBson(index)
	retAsIDocumentBase := makeDocument(res.model, result)
	return retAsIDocumentBase
}

// retrieves the record bson at the given index, reading additional records from the db driver as needed
func (res *Result) atBson(index uint) bson.M {
	// read responses until we have the record we need in lookback cache
	for res.cursorIndex <= index {
		more := res.readNextToLookback()
		if !more {
			log.Panic(&mongoidError.IndexOutOfBounds{})
		}
	}
	result := res.lookback[index]
	return result
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

// One returns a single document from the Result, also ensuring that only exactly one record was available to be read.
// One will panic if the Result contains more than one record or zero records.
// This method will panic if Streaming() was enabled.
func (res *Result) One() IDocumentBase {
	log.Debug("Result.One()")
	if res.streaming {
		log.Panic(&mongoidError.InvalidOperation{
			MethodName: "Result.One",
			Reason:     "Cannot perform random access when Result.IsStreaming()",
		})
	}

	// Count() will read all records into memory and cause the the cursor to be closed
	if count := res.Count(); count <= 0 || count > 1 {
		// panic here, because we have the wrong number of results
		if count == 0 {
			log.Panic(&mongoidError.NotFound{})
		} else {
			log.Panic(&mongoidError.MongoidError{}) // this could be a better type, maybe we need a new one
		}
	}
	return res.at(0)
}

// read the next result from the query cursor, and append it to the lookback cache
func (res *Result) readNextToLookback() bool {
	log.Trace("Result.readNextToLookback()")
	if res.closed {
		return false // closed cursor means there is definitely no more to read
	}
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
	} else { // close db cursor when we know there is no more data
		res.close()
	}
	return more
}

// readNext reads the next result from the query cursor, if one is present it will write the value into the bson.M at the given pointer and return true.
// If there are no more records available, it will return false.
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
	if res.closed {
		return false // closed cursor means there is definitely no more to read
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
	} else { // close db cursor when we know there is no more data
		res.close()
	}
	return more
}

// ForEach will call the given function once for each result, in the order they were returned by the server.
// The given function "fn" should accept an IDocumentBase as the only parameter.
// The given function "fn" may return a non-nil error value to halt further iterations - the return value is passed upward and returned by ForEach.
//
// Example:
//    ret := myResult.ForEach(func(v IDocumentBase) error {
//    	// do something with v
//    	v.Attribute = "New Value" // change it...
//    	v.Save() = "New Value" // save it...
//    	v.Delete() // or delete it...
//
//    	// the return value signals whather processing should halt or continue
//    	return nil // returning nil indicates we wish to continue - this function will be called again with the next record if there is one
//    	// return error.New("done") // or we could return some non-nil error value to signal that we would like to halt, skipping any remaining records
//    })
//
func (res *Result) ForEach(fn func(IDocumentBase) error) error {
	// the heavy lifting is within ForEachBson
	return res.ForEachBson(func(v bson.M) error {
		asIDocumentBase := makeDocument(res.model, v)
		return fn(asIDocumentBase)
	})
}

// ForEachBson is similar to ForEach, but provides the raw bson.M instead of an IDocumentBase object
func (res *Result) ForEachBson(fn func(bson.M) error) error {
	if !res.streaming { // non-streaming implementation (records are stored to lookback cache as they are read)
		count := res.Count() // this will read all records and close the mongo driver cursor for us
		for i := uint(0); i < count; i++ {
			result := res.atBson(i)
			r := fn(result) // run the given fn
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

// ToBsonAry returns the results as a slice of []bson.M
func (res *Result) ToBsonAry() []bson.M {
	resultAry := make([]bson.M, 0)
	res.ForEachBson(func(v bson.M) error {
		resultAry = append(resultAry, v)
		return nil
	})
	return resultAry
}
