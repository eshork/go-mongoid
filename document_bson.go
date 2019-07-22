package mongoid

import (
	"mongoid/log"
	// "github.com/iancoleman/strcase"
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	// "time"
)

// BsonDocument is a convenience re-export of bson.M from go.mongodb.org/mongo-driver/bson
type BsonDocument = bson.M

func bsonDeepCopyElement(element interface{}) interface{} {
	switch element.(type) {
	case bson.M:
		// log.Error("BsonDocumentDeepCopy - bson.M")
		bsonM := element.(bson.M)
		retBsonM := make(bson.M)
		for k, v := range bsonM {
			retBsonM[k] = bsonDeepCopyElement(v)
		}
		return retBsonM
	case bson.A:
		// log.Error("BsonDocumentDeepCopy - bson.A")
		bsonA := element.(bson.A)
		retBsonA := make(bson.A, 0)
		for _, v := range bsonA {
			retBsonA = append(retBsonA, bsonDeepCopyElement(v))
		}
		return retBsonA
	default:
		// log.Error("BsonDocumentDeepCopy - default")
	}
	newElementValue := reflect.ValueOf(element)
	if newElementValue.IsValid() {
		return newElementValue.Interface()
	}
	return element
}

// BsonDocumentDeepCopy returns a full deep copy of the original bson document, making it safe to edit in any way without risk of altering the original source
func BsonDocumentDeepCopy(bsonDoc BsonDocument) BsonDocument {
	log.Trace("BsonDocumentDeepCopy()")
	retBson := bsonDeepCopyElement(bsonDoc).(BsonDocument)
	return retBson
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// parses StructField "bson" options (if any) and returns them in a more usable form
func getBsonStructTagOpts(structField reflect.StructField) (fieldName string, omitempty bool, null bool, inline bool) {
	structTag := strings.TrimSpace(structField.Tag.Get("bson"))
	if structTag == "" {
		return "", false, false, false // early exit when there's nothing to do
	}
	tagParts := strings.Split(structTag, ",") // find the first comma, everything up until there is the field name
	fieldName = strings.TrimSpace(tagParts[0])
	if len(tagParts) == 1 { // whole string is the field name
		return fieldName, false, false, false
	}
	for i := 1; i < len(tagParts); i++ {
		switch strings.TrimSpace(strings.ToLower(tagParts[i])) {
		case "omitempty":
			omitempty = true
		case "null":
			null = true
		case "inline":
			inline = true
		}
	}
	return fieldName, omitempty, null, inline
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func makeBsonDocumentDiff(oldBson, newBson BsonDocument) BsonDocument {
	log.Trace("makeBsonDocumentDiff()")
	return makeBsonMDiff(oldBson, newBson) // this is an unabashed direct handoff to a private reusable function
}

// builds a bson.M representation of pending changes, based on given inputs of oldBson & newBson
// Key/Value pairs are returned only for instances where a difference is found between the two given states, otherwise they are excluded.
// The key/value included pairs will always reflect the "new" state, according to the given args.
// Unset keys or otherwise missing values will have the value side of their key/value pair set to 'nil', to reflect the newly unset status.
// Consumers should interpret nil-value keys in accordance to their own particular data situations
func makeBsonMDiff(oldBson, newBson bson.M) bson.M {
	changedBson := make(bson.M) // return value
	touchedBsonKeys := make(map[string]bool)
	// find missing and changed entries
	for key, oldValue := range oldBson {
		touchedBsonKeys[key] = true
		// log.Trace(key, ":", oldValue)
		newValue, newOk := newBson[key]
		if newOk != true { // key/value is missing within the new state
			changedBson[key] = nil
		} else { // key/value is still within the new, but may be changed (yet unknown)
			wasChanged := !reflect.DeepEqual(oldValue, newValue)
			if wasChanged {
				changedBson[key] = newValue
			}
		}
	}
	// add in any new keys that weren't already in oldBson
	for key, newValue := range newBson {
		if _, touched := touchedBsonKeys[key]; !touched {
			changedBson[key] = newValue
		}
	}
	if len(changedBson) == 0 { // if there's nothing in our change list, just return nil
		return nil
	}
	return changedBson // otherwise hand back what was built
}

// TODO: evaluate this function -- it seems maybe unnecessary given reflect.DeepEqual functionality
// builds a bson.A representation of pending changes, based on given inputs of oldBson & newBson
func makeBsonADiff(oldBson, newBson bson.A) bson.A {
	log.Panic("makeBsonADiff() - NYI")
	changedBson := make(bson.A, 0)
	return changedBson
}
