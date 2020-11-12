package mongoid

import (
	// "mongoid/log"
	"go.mongodb.org/mongo-driver/bson"
	// "strconv"
)

// Query ...
type Query map[string]interface{}

// Q is shorthand for Query
type Q = Query

// Criteria facilitate the query-building process
type Criteria interface {
	Find(ids ...ObjectID) Criteria
	Where(where ...Query) Criteria
	X() *Result // see: [criteria_x.go] func (criteria *criteriaStruct) X() *Result
	getPrevCriteria() Criteria
	toBsonD() bson.D
}

const (
	_ = iota
	findCriteria
	whereCriteria
)

type criteriaStruct struct {
	criteriaType   int
	sourceModel    *ModelType
	prevCriteria   *criteriaStruct
	thisQuery      Query
	thisQueryBsonD bson.D
}

func (criteria *criteriaStruct) getPrevCriteria() Criteria {
	return criteria.prevCriteria
}

func (criteria *criteriaStruct) toBsonD() bson.D {
	switch criteria.criteriaType {
	case whereCriteria:
		if criteria.thisQueryBsonD == nil {
			criteria.thisQueryBsonD = criteriaWhereToBsonD(criteria)
		}
		return criteria.thisQueryBsonD
	case findCriteria:
		panic("findCriteria.toBsonD") // TODO implement this
	}
	return bson.D{}
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// func (criteria *ModelType) First() {
// }
// func (criteria *Criteria) First() {
// }
// func (criteria *ModelType) Last() {
// }
// func (criteria *Criteria) Last() {
// }

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// func (criteria *Criteria) build() {
// 	log.Error("criteria.build")
// }
