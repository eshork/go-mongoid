package mongoid

import (
	"mongoid/log"
	// "fmt"
	"go.mongodb.org/mongo-driver/bson"
	"strings"
)

// Where adds criteria that must be matched in order to return results
func (model *ModelType) Where(where ...Query) Criteria {
	log.Debug("ModelType.Where ", where)
	return criteriaWhere(model, nil, where...)
}

// Where adds criteria that must be matched in order to return results
func (criteria *criteriaStruct) Where(where ...Query) Criteria {
	log.Debug("Criteria.Where ", where)
	return criteriaWhere(nil, criteria, where...)
}

func criteriaWhere(srcModel *ModelType, prevCriteria *criteriaStruct, where ...Query) Criteria {
	curPrevCriteria := prevCriteria
	for _, thisWhereQuery := range where {
		log.Traceln("New Criteria.Where ", thisWhereQuery)
		newCriteria := criteriaStruct{
			sourceModel:  srcModel,
			criteriaType: whereCriteria,
			prevCriteria: curPrevCriteria,
			thisQuery:    thisWhereQuery,
		}
		newCriteria.thisQueryBsonD = criteriaWhereToBsonD(&newCriteria)
		curPrevCriteria = &newCriteria

		// log.Tracef("newCriteria %+v\n", newCriteria)
		// log.Printf("curPrevCriteria %+v\n", curPrevCriteria)
	}
	return curPrevCriteria
}

func criteriaWhereToBsonD(where *criteriaStruct) bson.D {
	if where == nil {
		return bson.D{}
	}
	var bsonD bson.D
	switch len(where.thisQuery) {
	case 0:
		bsonD = bson.D{}
	default:
		bsonD = bson.D{}
		for k, v := range where.thisQuery {
			element := normalizeQueryElement(k, v)
			bsonE := bson.E{
				Key:   element[0].Key,
				Value: element[0].Value,
			}
			bsonD = append(bsonD, bsonE)
		}
	}
	// log.Tracef("bsonD %+v\n", bsonD)
	return bsonD

}

func normalizeQueryElement(fieldKey string, fieldValue interface{}) bson.D {
	// find operator
	operatorStart := strings.LastIndex(fieldKey, ".$")
	if operatorStart == -1 {
		// if no operator is found, early return without changes
		return bson.D{{Key: fieldKey, Value: fieldValue}}
	}

	operator := fieldKey[operatorStart:]
	switch operator {
	case ".$eq":
		fallthrough
	case ".$gt":
		fallthrough
	case ".$gte":
		fallthrough
	case ".$in":
		fallthrough
	case ".$lt":
		fallthrough
	case ".$lte":
		fallthrough
	case ".$ne":
		fallthrough
	case ".$nin":
		fallthrough
	case ".$exists":
		// extract the operator (less the . separator)
		operator = operator[1:]
		// remove the operator from the field name
		fieldKey = fieldKey[:operatorStart]
		// reposition the fieldValue into a map
		fieldValue = bson.M{operator: fieldValue}
	}

	// return the current state
	return bson.D{{Key: fieldKey, Value: fieldValue}}
}

// https://docs.mongodb.com/manual/tutorial/query-documents/
/*
// Equality
https://docs.mongodb.com/manual/reference/operator/query/eq/
db.inventory.find( { status: "D" } )
{ <field>: { $eq: <value> } }


$eq	Matches values that are equal to a specified value.
$gt	Matches values that are greater than a specified value.
$gte	Matches values that are greater than or equal to a specified value.
$in	Matches any of the values specified in an array.
$lt	Matches values that are less than a specified value.
$lte	Matches values that are less than or equal to a specified value.
$ne,$neq	Matches all values that are not equal to a specified value.
$nin	Matches none of the values specified in an array.



db.inventory.find( { status: { $in: [ "A", "D" ] } } )
*/
