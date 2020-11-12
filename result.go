package mongoid

import (
	"mongoid/log"
	// "strconv"
	// "go.mongodb.org/mongo-driver/bson"
)

// Result provides access to the data produced by database operations
type Result struct {
	// First() IDocumentBase
	// Last() IDocumentBase
	// Count() int
	// At(index int) IDocumentBase

	//
	// ToA() error
	// ForEach(f func () error) error
}

// First returns an interface to the first Document in the result set
func (res *Result) First() IDocumentBase {
	log.Panicf("NYI - Result.First()")
	return nil
}

// Last returns an interface to the last Document in the result set
func (res *Result) Last() IDocumentBase {
	log.Panicf("NYI - Result.Last()")
	return nil
}

// Count returns the number of Document in the result set
func (res *Result) Count() int {
	log.Panicf("NYI - Result.Count()")
	return 0
}

// At returns an interface to the Document in the result set at the given index (range is 0 to count-1)
func (res *Result) At(index int) IDocumentBase {
	log.Panicf("NYI - Result.Last()")
	return nil
}
