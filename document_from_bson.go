package mongoid

import (
	"mongoid/log"
	"mongoid/util"

	"reflect"

	"github.com/iancoleman/strcase"
	"go.mongodb.org/mongo-driver/bson"
	// "strings"
	// "time"
)

var reflectTypeObjectID = reflect.TypeOf(ZeroObjectID)

// Apply values to the given struct from the given bsonM
func structValuesFromBsonM(rawStructPtr interface{}, bsonM bson.M) (updated bool) {
	updated = false // initially returns unupdated state (updated = false)
	if rawStructPtr == nil {
		log.Panic("rawStructPtr interface{} cannot be nil")
	}
	if len(bsonM) == 0 {
		log.Panic("bsonM bson.M cannot be empty")
	}
	log.Trace("structValuesFromBsonM(<detecting>)")
	handleType := reflect.TypeOf(rawStructPtr)
	if handleType.Kind() != reflect.Ptr { // ensure handleType is always a pointer to a struct
		log.Panic("rawStructPtr interface{} type.Kind() must be reflect.Ptr, found ", handleType.Kind())
	}
	handleStructType := handleType.Elem()
	log.Tracef("structValuesFromBsonM(%s) = %v", handleStructType, bsonM)
	handleValue := reflect.Indirect(reflect.ValueOf(rawStructPtr)) // TODO: optimize everything above this; currently unsure if this is efficient (expecting likely not)

	for i := 0; i < handleStructType.NumField(); i++ {
		// walk each field
		// set the value according to bsonM where possible (recursion may occur)
		field := handleStructType.Field(i) // Get the field type - https://golang.org/pkg/reflect/#StructField
		fieldValue, found := structFieldValueFromBsonM(field, bsonM)
		if found { // apply the given value
			curFieldValue := handleValue.Field(i) // Get the current field value as reflect.Value - https://golang.org/pkg/reflect/#Value
			// log.Info("updating fieldValue WAS: .", field.Name, " = ", curFieldValue)
			// log.Info("updating fieldValue Setting: .", field.Name, " = ", fieldValue)
			curFieldValue.Set(fieldValue)
			// log.Info("updated fieldValue NOW: .", field.Name, " = ", curFieldValue)
			updated = true // record that a field was updated
		}
	}
	// log.Errorf("structValuesFromBsonM() result[%t]: %+v", updated, handleValue) // TODO: remove this - debug output
	return
}

// If `found`, this will provide a fully formed value that can be Set upon the field that was given
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
		retValue = sliceValueFromBsonA(fieldType, bsonMfieldValue.(bson.A))
		// log.Fatal("reflect.Slice: ", retValue)

		// ptr to slice?

	default:
		// for default (concrete) types, we need to make a new copy on the heap (otherwise they sometimes fail CanAddr())
		// log.Warnf("default: retValue field value type: %s", fieldType)
		// log.Warnf("default: retValue =: %s", retValue)
		retValue = reflect.Indirect(reflect.ValueOf(bsonMfieldValue))
		valuePtrValue := reflect.New(retValue.Type())
		valuePtrValue.Elem().Set(retValue)
		retValue = valuePtrValue.Elem()
	}

	// If the field type kind is a pointer, then we need to return a pointer as the new value
	if fieldTypeKind == reflect.Ptr {
		retValue = retValue.Addr()
		// log.Warnf("retValue PTR field value: %v", retValue.Elem())
	}
	// log.Warnf("retValue field value: %v", retValue)
	return retValue, true
}

// accepts input values bson.A and a destination sliceType reflect.Type (ie []int, []string, []any) and returns a new reflect.Value containing the data
func sliceValueFromBsonA(sliceType reflect.Type, bsonA bson.A) (reflectValue reflect.Value) {
	log.Trace("sliceValueFromBsonA")

	if sliceType.Kind() != reflect.Slice {
		log.Panicf("reflectValueFromBsonA expects sliceType(%s).Kind() == reflect.Slice", sliceType)
	}

	sliceElemType := sliceType.Elem()
	sliceElemTypeKind := sliceElemType.Kind()

	// log.Warn(bsonA)
	// log.Warn("sliceType = ", sliceType)
	// log.Warn("sliceElemType = ", sliceElemType)
	// log.Warn("sliceElemTypeKind = ", sliceElemTypeKind)

	lenBsonA := len(bsonA)

	if lenBsonA <= 0 {
		return reflect.MakeSlice(sliceType, 0, 0)
	}

	// make a new slice of the sliceType
	sliceValue := reflect.MakeSlice(sliceType, lenBsonA, lenBsonA)

	if sliceElemTypeKind == reflect.Ptr { // pointer-based types
		// assign each slice index by ptr
		// log.Warn("PTRs!")

		sliceElemBaseType := sliceElemType.Elem()
		if sliceElemBaseType.Kind() == reflect.Struct {
			// ptrs to custom struct type
			// log.Warn("sliceElemBaseType = ", sliceElemBaseType)
			for i := 0; i < lenBsonA; i++ {
				if bsonA[i] != nil { // only process non-nil values
					bsonM, ok := bsonA[i].(bson.M)
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
				if bsonA[i] != nil { // only process non-nil values
					bsonIndexValue := reflect.ValueOf(bsonA[i])
					valuePtrValue := reflect.New(bsonIndexValue.Type())
					valuePtrValue.Elem().Set(bsonIndexValue)
					sliceIndexValue := sliceValue.Index(i)
					sliceIndexValue.Set(valuePtrValue)
				}
			}
		}
	} else { // non-pointer-based types (concrete types)
		if sliceElemTypeKind == reflect.Struct {
			// custom struct type
			// log.Warn("sliceElemType = ", sliceElemType)
			for i := 0; i < lenBsonA; i++ {
				// log.Error("bsonA[i]= ", bsonA[i])
				bsonM, ok := bsonA[i].(bson.M)
				if !ok {
					log.Panicf("sliceValueFromBsonA - value at index %d must be valid bson.M when slice kind is struct of type: %s", i, sliceElemType)
				}
				newStructPtr := reflect.New(sliceElemType)
				structValuesFromBsonM(newStructPtr.Interface(), bsonM)
				sliceIndexValue := sliceValue.Index(i)
				sliceIndexValue.Set(newStructPtr.Elem())
				// log.Fatal("sliceElemTypeKind == reflect.Struct")

			}
		} else {
			// other built-in type
			for i := 0; i < lenBsonA; i++ {
				bsonIndexValue := reflect.ValueOf(bsonA[i])
				sliceIndexValue := sliceValue.Index(i)
				sliceIndexValue.Set(bsonIndexValue)
			}
		}
	}

	return sliceValue
}
