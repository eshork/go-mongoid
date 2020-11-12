package mongoid

import (
	"mongoid/log"
	// "strconv"
	// "go.mongodb.org/mongo-driver/bson"
)

// THIS SHOULD BE THE FILE WHERE INITAIL CRITERIA->RESULT IMPLEMENTATION GOES

// X will force eXecution of the criteria query, caching results of the Criteria
func (criteria *criteriaStruct) X() *Result {
	log.Panicf("NYI - Criteria.X()")
	return nil
}
