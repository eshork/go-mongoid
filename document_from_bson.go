package mongoid

import (
	"mongoid/log"
	"mongoid/util"
	"reflect"

	"github.com/iancoleman/strcase"
	"go.mongodb.org/mongo-driver/bson"
)

// "time"

var reflectTypeObjectID = reflect.TypeOf(ZeroObjectID)

// Apply matching values to the given struct (passed by pointer) from the given bsonM.
// If a value for a field is not found within the given bsonM, it will be skipped without error.
// Returns true if any value in the struct was written, even if the new value was equal to the existing value (existing values are not examined).
func structValuesFromBsonM(rawStructPtr interface{}, bsonM bson.M) (updated bool) {
	// zero value is updated=false, no need to explicitly set
	if rawStructPtr == nil {
		log.Panic("rawStructPtr interface{} cannot be nil")
	}
	//
	// if len(bsonM) == 0 { // an empty bson.M is okay, but nothing will change
	// 	return false // so we can return early
	// }
	log.Trace("structValuesFromBsonM(<detecting>)")
	handleType := reflect.TypeOf(rawStructPtr)
	if handleType.Kind() != reflect.Ptr { // ensure handleType is always a pointer to a struct
		log.Panic("rawStructPtr interface{} type.Kind() must be reflect.Ptr, but found: ", handleType.Kind())
	}
	handleStructType := handleType.Elem()
	log.Tracef("structValuesFromBsonM(%s) = %v", handleStructType, bsonM)
	handleValue := reflect.Indirect(reflect.ValueOf(rawStructPtr)) // TODO: optimize everything above this; currently unsure if this is efficient (expecting likely not)

	for i := 0; i < handleStructType.NumField(); i++ {
		// walk each field
		// set the value according to bsonM where possible (recursion may occur)
		field := handleStructType.Field(i) // Get the field type - https://golang.org/pkg/reflect/#StructField
		newFieldValue, found := structFieldValueFromBsonM(field, bsonM)
		if found { // apply the given value
			curFieldValue := handleValue.Field(i) // Get the current field value as reflect.Value - https://golang.org/pkg/reflect/#Value
			util.SetValueByInterfacePtr(curFieldValue.Addr().Interface(), newFieldValue.Interface())
			updated = true // record that a field was updated
		}
	}
	// log.Errorf("structValuesFromBsonM() result[%t]: %+v", updated, handleValue) // TODO: remove this - debug output
	return
}

// Retrieves the value for a struct field from the given bsonM, following the direction of struct field tags if present.
// Will return found=true when a matching value was available within the given bsonM, otherwise found=false.
// This will not assign a matching value to the struct field -- you will need to do that yourself.
func structFieldValueFromBsonM(field reflect.StructField, bsonM bson.M) (retValue reflect.Value, found bool) {
	// skip private fields
	if field.PkgPath != "" {
		return retValue, false
	}

	// extract relevant struct tag info, if any
	tagFieldName, _, _, tagInline := getBsonStructTagOpts(field)
	// tagFieldName, tagOmitempty, tagNull, tagInline := getBsonStructTagOpts(field)
	// TODO: remove this - debug output
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
		return retValue, false // field is explicitly omitted from bson, so we're done
	default: // anything else replaces the field name with the given name
		fieldName = tagFieldName
	}
	fieldType := field.Type
	fieldTypeKind := fieldType.Kind()
	// TODO: remove this - debug output
	// log.Debugf("fieldName: %s", fieldName)
	// log.Debugf("field type : %s", fieldType)
	// log.Debugf("field kind : %s", fieldTypeKind)

	// fetch the bsonM value (if present)
	bsonMfieldValue, bsonMhasKey := bsonM[fieldName]
	// log.Warnf("bsonMfieldValue: %+v", bsonMfieldValue)

	// unroll ptr vs value return type
	if fieldTypeKind == reflect.Ptr {
		if !tagInline && bsonMhasKey == true && bsonMfieldValue == nil { // if a non-inlined PTR kind and a null value is present, we can return early
			retValue = reflect.New(fieldType).Elem()
			return retValue, true
		}
		fieldType = fieldType.Elem() // get actual type (instead of a ptr to that type)
	}

	// if this is an inlined struct or inlined ptr to a struct, then do something special, otherwise continue per-normal
	if tagInline {
		// validate field type
		if !(fieldTypeKind == reflect.Struct ||
			(fieldTypeKind == reflect.Ptr && fieldType.Kind() == reflect.Struct)) {
			log.Panicf("invalid inlined field type (%s) - must be struct or *struct", fieldType)
			return retValue, false
		}

		// make an instance of the struct and try to fill it
		valuePtrValue := reflect.New(fieldType)
		retValueInterface := valuePtrValue.Interface()
		// structValuesFromBsonM(retValueInterface, bsonM)
		updated := structValuesFromBsonM(retValueInterface, bsonM)
		if !updated { // only continue if the current bsonM was found to contain at least one inlined value entry
			// log.Warnf("inlined struct field had no updates")
			return retValue, false
		}

		// if any fields were filled, return our new entity for assignment
		retValue = valuePtrValue.Elem()
		if fieldTypeKind == reflect.Ptr { // If the field type kind is a pointer, then we need to return a pointer as the new value
			retValue = retValue.Addr()
		}
		return retValue, true
	}

	// if this is not an inlined field and a value does not exist within the bsonM, we can abort
	if !bsonMhasKey {
		// log.Warnf("bsonMhasKey == FALSE")
		return retValue, false
	}

	// TODO should look into handling custom FromBSON marshallables here...
	if util.IsIfaceBsonMarshalSafe(bsonMfieldValue) {
		// log.Fatal("util.IsBsonMarshalSafe() is not a real thing -- that was more than 10 lies")
	}

	switch fieldType.Kind() {
	case reflect.Struct:
		valuePtrValue := reflect.New(fieldType)
		retValueInterface := valuePtrValue.Interface()
		structValuesFromBsonM(retValueInterface, bsonMfieldValue.(bson.M))
		retValue = valuePtrValue.Elem()
	case reflect.Slice:
		retValue = sliceValueFromUnknownAry(fieldType, bsonMfieldValue)
	default:
		// for default (concrete) types, we need to make a new copy on the heap (otherwise they sometimes fail CanAddr())
		retValue = reflect.ValueOf(util.MarshalFromDB(fieldType, bsonMfieldValue))
		valuePtrValue := reflect.New(retValue.Type())
		valuePtrValue.Elem().Set(retValue)
		retValue = valuePtrValue.Elem()
	}

	// If the field type kind is a pointer, then we need to return a pointer as the new value
	if fieldTypeKind == reflect.Ptr {
		retValue = retValue.Addr()
	}
	return retValue, true
}

// accepts unknownAry input and destination sliceType (ie []int, []string, []any) and returns a new reflect.Value containing the data
func sliceValueFromUnknownAry(sliceType reflect.Type, unknownAry interface{}) (reflectValue reflect.Value) {
	log.Trace("sliceValueFromUnknownAry")

	if sliceType.Kind() != reflect.Slice {
		log.Panicf("reflectValueFromBsonA expects sliceType(%s).Kind() == reflect.Slice", sliceType)
	}

	sliceElemType := sliceType.Elem()
	sliceElemTypeKind := sliceElemType.Kind()

	// log.Warn("bsonA")
	// log.Warn("sliceType = ", sliceType)
	// log.Warn("sliceElemType = ", sliceElemType)
	// log.Warn("sliceElemTypeKind = ", sliceElemTypeKind)

	lenBsonA := reflect.ValueOf(unknownAry).Len() // len(bsonA)

	if lenBsonA <= 0 {
		return reflect.MakeSlice(sliceType, 0, 0)
	}

	// make a new slice of the sliceType
	sliceValue := reflect.MakeSlice(sliceType, lenBsonA, lenBsonA)

	if sliceElemTypeKind == reflect.Ptr { // pointer-based types
		// assign each slice index by ptr
		sliceElemBaseType := sliceElemType.Elem()
		if sliceElemBaseType.Kind() == reflect.Struct {
			// ptrs to custom struct type
			for i := 0; i < lenBsonA; i++ {
				// if bsonA[i] != nil { // only process non-nil values
				if !reflect.ValueOf(unknownAry).Index(i).IsNil() { // only process non-nil values
					// bsonM, ok := bsonA[i].(bson.M)
					bsonM, ok := reflect.ValueOf(unknownAry).Index(i).Interface().(bson.M)
					if !ok {
						log.Panicf("sliceValueFromBsonA - value at index %d must be valid bson.M when slice kind is struct of type: %s", i, sliceElemBaseType)
					}
					newStructPtr := reflect.New(sliceElemBaseType)
					structValuesFromBsonM(newStructPtr.Interface(), bsonM)
					sliceIndexValue := sliceValue.Index(i)
					sliceIndexValue.Set(newStructPtr)
				}
			}
		} else {
			// ptrs to other built-in type
			for i := 0; i < lenBsonA; i++ {
				// if bsonA[i] != nil { // only process non-nil values
				if !reflect.ValueOf(unknownAry).Index(i).IsNil() { // only process non-nil values
					bsonIndexValue := reflect.ValueOf(unknownAry).Index(i)
					valuePtrValue := reflect.New(bsonIndexValue.Type())
					valuePtrValue.Elem().Set(bsonIndexValue)
					sliceIndexValue := sliceValue.Index(i)
					sliceIndexValue.Set(valuePtrValue)
				}
			}
		}
	} else if sliceElemTypeKind == reflect.Interface { // interface type
		for i := 0; i < lenBsonA; i++ {
			bsonIndexValue := reflect.ValueOf(unknownAry).Index(i)
			sliceIndexValue := sliceValue.Index(i)
			sliceIndexValue.Set(bsonIndexValue)
		}
	} else { // non-pointer-based types (concrete types)
		if sliceElemTypeKind == reflect.Struct {
			// custom struct type
			for i := 0; i < lenBsonA; i++ {
				bsonM, ok := reflect.ValueOf(unknownAry).Index(i).Interface().(bson.M)
				if !ok {
					log.Panicf("sliceValueFromBsonA - value at index %d must be valid bson.M when slice kind is struct of type: %s", i, sliceElemType)
				}
				newStructPtr := reflect.New(sliceElemType)
				structValuesFromBsonM(newStructPtr.Interface(), bsonM)
				sliceIndexValue := sliceValue.Index(i)
				sliceIndexValue.Set(newStructPtr.Elem())
			}
		} else {
			// other built-in type
			for i := 0; i < lenBsonA; i++ {
				bsonIndexValue := reflect.ValueOf(unknownAry).Index(i)
				sliceIndexValue := sliceValue.Index(i)
				indexValue := reflect.ValueOf(util.MarshalFromDB(sliceElemType, bsonIndexValue.Interface()))
				sliceIndexValue.Set(indexValue)
			}
		}
	}

	return sliceValue
}
