package mongoid

import (
	"mongoid/log"
	"mongoid/util"

	"reflect"
	"strconv"

	"github.com/iancoleman/strcase"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/primitive"
	// "strings"
	// "time"
)

// ToBson converts the document model object into a BsonDocument.
// This makes the object easily serializable and comparable.
func (d *Base) ToBson() BsonDocument {
	log.Trace("Base.ToBson()")
	bsonOut := structToBsonM(d.DocumentBase())
	return bsonOut
}

// ToUpdateBson converts the document model object into a BsonDocument
func (d *Base) ToUpdateBson() BsonDocument {
	log.Trace("Base.ToUpdateBson()")
	updateBson := bson.M{
		"$set": d.Changes(),
	}
	// TODO - add an $unset operator to appropriately remove unset fields (instead of just setting them to null, like it currently does)
	return updateBson
}

func structToBsonM(rawStructPtr interface{}) bson.M {
	// log.Trace("structToBsonM(<detecting>)")
	retMap := make(bson.M)
	handleType := reflect.TypeOf(rawStructPtr)
	// ensure handleType is always a pointer to a struct, otherwise bad stuff might happen later
	if handleType.Kind() != reflect.Ptr {
		handleType = reflect.PtrTo(handleType)
	}

	handleStructType := handleType.Elem()
	// log.Tracef("structToBsonM(%+v)", handleStructType)

	handleValue := reflect.Indirect(reflect.ValueOf(rawStructPtr))
	// log.Printf("structToBsonM INPUT: %+v", handleValue) // TODO CLEANUP

	for i := 0; i < handleStructType.NumField(); i++ {
		field := handleStructType.Field(i) // Get the field type, returns https://golang.org/pkg/reflect/#StructField
		fieldValue := handleValue.Field(i) // Get the field value, returns https://golang.org/pkg/reflect/#Value
		// log.Tracef("Field Name: %s", field.Name) // TODO CLEANUP
		bsonValue := structFieldToBsonM(field, fieldValue)
		// log.Errorf("Field BSON Value: %+v", bsonValue) // TODO CLEANUP
		// log.Errorf("Field BSON Value(len): %d", len(bsonValue)) // TODO CLEANUP

		// append results
		if len(bsonValue) > 0 {
			for k, v := range bsonValue {
				retMap[k] = v
			}
		}
	}
	if len(retMap) == 0 {
		return nil
	}
	return retMap
}

/*
omitempty - Only include the field if it's not set to the zero-value for the type,
null - Set a field value to "null" if it's set to zero-value for the type and for empty slices or maps.
inline - Inlines the field, which must be a struct, causing all of its fields to be processed as if they were part of the outer struct
*/

// builds an appropriate bson.M from a StructField (definition) and Value (content) pair
// returns an empty bson.M if the given field should be omitted
// includes exported fields (capitalized) by default unless bson name is unset via `bson:"-"` tag
// unexported fields cannot be converted to bson, even if given an explicit bson field name (embedded anonymous structs are the exception, if inlined, even then only the exported fields)
// exported struct-type fields are embedded unless explicitly inlined via `bson:",inline"`
// anonymous structs are omitted unless explicitly inlined via `bson:",inline"`
// zero value fields that would normally be included may be omitted via `bson:",omitempty"`
// nil pointers (zero value) result in a "null" value unless omitted via `bson:",omitempty"`

// zero value fields  `bson:",omitempty"`
func structFieldToBsonM(field reflect.StructField, fieldValue reflect.Value) bson.M {
	// log.Trace("structFieldToBsonM")

	// unexported struct fields are automatically skipped - no need to do any other work on them
	if field.PkgPath != "" {
		return bson.M{}
	}

	// extract relevant struct tag info, if any
	tagFieldName, tagOmitempty, tagNull, tagInline := getBsonStructTagOpts(field)
	// log.Debugf("tagFieldName %s", tagFieldName)
	// log.Debugf("tagOmitempty %t", tagOmitempty)
	// log.Debugf("tagNull %t", tagNull)
	// log.Debugf("tagInline %t", tagInline)

	// process field name
	fieldName := field.Name // default bson field name is the struct field name, converted to snake_case
	switch tagFieldName {
	case "": // empty tagFieldName means no explicit field name substitution was given, so make due with the existing struct field name
		fieldName = strcase.ToSnake(fieldName)
	case "-": // "-" indicates this field is explicitly omitted
		return bson.M{} // ... so we're done
	default: // anything else replaces the field name with the given name
		fieldName = tagFieldName
	}

	fieldType := field.Type
	fieldTypeKind := fieldType.Kind()

	// if ptr and value is nil, this is really easy to solve
	if fieldTypeKind == reflect.Ptr {
		if fieldValue.IsNil() {
			if tagOmitempty == true {
				return bson.M{} // field has no entry if it cannot represent and is omit on empty
			}
			return bson.M{fieldName: nil} // nil ptr means nil value
		}
		fieldValue = fieldValue.Elem() // reset fieldValue to concrete value (instead of a ptr to that value)
	}

	fieldValueKind := fieldValue.Kind() // this is used more than once below, so store it on the stack

	if tagInline { // bad inline opt protection; validates field type supports inline (ie, struct or *struct) else panics
		if fieldValueKind != reflect.Struct { // must be fieldValueKind == reflect.Struct
			log.Panicf("invalid type for 'inline' field (%s) - must be struct or *struct", fieldValueKind)
		}
	}

	// TODO should look into handling custom ToBSON marshallables here...
	// if util.IsIfaceBsonMarshalSafe(fieldValue.Interface()) {
	// 	log.Fatal("util.IsBsonMarshalSafe() is not a real thing -- that was more than 10 lies")
	// }
	if !util.IsIfaceBsonMarshalSafe(fieldValue.Interface()) {
		switch fieldValueKind {
		case reflect.Struct:
			// omit the field if struct is anonymous and not inlined
			if tagInline != true && field.Anonymous == true {
				return bson.M{} // anonymous non-inline structs are omitted
			} // TODO: one day, wouldn't it be nice to be able to reach into the raw bson of an object, so we might get a look at values that weren't necessarily assigned into struct fields?
			// get the struct converted into a bson map
			structBsonM := structToBsonM(fieldValue.Interface())

			// if inline, pop all the elements up to the top level (drops the existing struct field name)
			if tagInline == true {
				retMap := make(bson.M)
				for k, v := range structBsonM {
					retMap[k] = v
				}
				return retMap
			}
			return bson.M{fieldName: structBsonM}

		case reflect.Slice:
			retBsonA := indexableValueToBsonA(reflect.Indirect(fieldValue))
			return bson.M{fieldName: retBsonA}
		case reflect.Bool:
			return bson.M{fieldName: marshalToDB(fieldValue.Interface())}
		case reflect.String:
			return bson.M{fieldName: marshalToDB(fieldValue.Interface())}
		case reflect.Int:
			return bson.M{fieldName: marshalToDB(fieldValue.Interface())}
		case reflect.Int8:
			return bson.M{fieldName: marshalToDB(fieldValue.Interface())}
		case reflect.Int16:
			return bson.M{fieldName: marshalToDB(fieldValue.Interface())}
		case reflect.Int32:
			return bson.M{fieldName: marshalToDB(fieldValue.Interface())}
		case reflect.Int64:
			return bson.M{fieldName: marshalToDB(fieldValue.Interface())}
		case reflect.Uint:
			return bson.M{fieldName: marshalToDB(fieldValue.Interface())}
		case reflect.Uint8:
			return bson.M{fieldName: marshalToDB(fieldValue.Interface())}
		case reflect.Uint16:
			return bson.M{fieldName: marshalToDB(fieldValue.Interface())}
		case reflect.Uint32:
			return bson.M{fieldName: marshalToDB(fieldValue.Interface())}
		case reflect.Uint64:
			return bson.M{fieldName: marshalToDB(fieldValue.Interface())}
		default:
			log.Panicf("unhandled builtin type: %v", fieldValueKind)
			// 	return bson.M{fieldName: marshalToDB(fieldValue.Interface())}
		}
	}
	if oid, ok := fieldValue.Interface().(primitive.ObjectID); ok {
		return bson.M{fieldName: marshalToDB(oid)}
	}
	if marshaler, ok := reflect.Indirect(fieldValue).Interface().(bsoncodec.Marshaler); ok {
		log.Error("this code is untested (e17635b9)")
		marshaledValue, err := marshaler.MarshalBSON()
		if err != nil {
			log.Panic(err)
		}
		return bson.M{fieldName: marshalToDB(marshaledValue)}
	}

	// standard values types and unknowns
	// if the concrete value is the zero-value, we may have special handling
	if !fieldValue.IsValid() || reflect.Zero(fieldValue.Type()).Interface() == fieldValue.Interface() {
		if tagOmitempty {
			return bson.M{} // ... so we're done
		}
		if tagNull { // null zero values
			return bson.M{fieldName: nil} // replace zero-value with nil/null
		}
	}
	log.Error("default bson enc: ", fieldValueKind)
	return bson.M{
		fieldName: fieldValue.Interface(),
	}
}

// Casts the given fromValue into a database suitable storage type, returning an interface to the newly cast value.
// If fromValue does not require conversion, it may be returned directly, but it is not guaranteed to do so.
// If a value conversion would result in loss of data or precision, this function will panic.
func marshalToDB(fromValue interface{}) interface{} {
	switch fromValue.(type) {
	case primitive.ObjectID:
		return fromValue
	case bool:
		return fromValue
	case string:
		return fromValue
	case int:
		return int32(fromValue.(int))
	case int8:
		return int32(fromValue.(int8))
	case int16:
		return int32(fromValue.(int16))
	case int32:
		return int32(fromValue.(int32))
	case int64:
		return int64(fromValue.(int64))
	case uint:
		return int64(fromValue.(uint))
	case uint8:
		return int32(fromValue.(uint8))
	case uint16:
		return int32(fromValue.(uint16))
	case uint32:
		return int64(fromValue.(uint32))
	case uint64:
		val := fromValue.(uint64)
		return strconv.FormatUint(val, 10)
	case *uint64:
		if fromValue != nil {
			return marshalToDB(*(fromValue.(*uint64)))
		}
		return nil
	default:
		log.Panicf("default marfshalToDB: %v ", reflect.TypeOf(fromValue))
		return fromValue
	}
}

// accepts a reflect.Value of an indexable (slice or array) and returns bson.A
func indexableValueToBsonA(indexableValue reflect.Value) bson.A {
	// log.Trace("indexableValueToBsonA")
	indexableLen := indexableValue.Len()
	retBsonA := make(bson.A, indexableLen)
	for i := 0; i < indexableLen; i++ {
		indexableValueAtIndex := indexableValue.Index(i)
		if indexableValueAtIndex.Kind() == reflect.Ptr {
			if !indexableValueAtIndex.IsNil() { // only store the non-nil values for pointers
				indirectValueAtIndex := reflect.Indirect(indexableValueAtIndex)
				if indirectValueAtIndex.Kind() == reflect.Struct {
					retBsonA[i] = structToBsonM(indirectValueAtIndex.Interface())
				} else {
					retBsonA[i] = indirectValueAtIndex.Interface()
				}
			}
		} else {
			indirectValueAtIndex := reflect.Indirect(indexableValueAtIndex)
			if indirectValueAtIndex.Kind() == reflect.Struct {
				retBsonA[i] = structToBsonM(indirectValueAtIndex.Interface())
			} else {
				retBsonA[i] = indirectValueAtIndex.Interface()
			}
		}
	}
	return retBsonA
}
