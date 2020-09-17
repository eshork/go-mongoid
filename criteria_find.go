package mongoid

import (
	"mongoid/log"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
)

// Find a document or multiple documents by their ids
func (model *ModelType) Find(ids ...ObjectID) Criteria {
	log.Debug("ModelType.Find ", ids)
	return criteriaFind(nil, ids...)
}

// Find a document or multiple documents by their ids
func (criteria *criteriaStruct) Find(ids ...ObjectID) Criteria {
	log.Debug("Criteria.Find ", ids)
	return criteriaFind(criteria, ids...)
}

func criteriaFind(prevCriteria *criteriaStruct, ids ...ObjectID) Criteria {
	// log.Print(ids)

	thisQuery := Query{}
	for k, v := range ids {
		thisQuery[strconv.Itoa(k)] = v
	}
	log.Traceln("New Criteria.Find ", thisQuery)
	newCriteria := criteriaStruct{
		criteriaType: findCriteria,
		prevCriteria: prevCriteria,
		thisQuery:    thisQuery,
	}
	// log.Printf("%+v\n", newCriteria)
	return &newCriteria
}

func criteriaFindToBsonD(find *criteriaStruct) bson.D {
	var idsArray bson.A
	for _, k := range find.thisQuery {
		objectID := k.(ObjectID)
		idsArray = append(idsArray, objectID)
	}
	var bsonD bson.D
	switch len(idsArray) {
	case 0:
		bsonD = bson.D{}
	case 1:
		bsonD = bson.D{{Key: "_id", Value: idsArray[0]}}
	default:
		bsonA := bson.A{}
		for _, v := range idsArray {
			bsonA = append(bsonA, bson.M{"_id": v})
		}
		bsonD = bson.D{{Key: "$or", Value: bsonA}}
	}
	return bsonD
}
