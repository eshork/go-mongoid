package util

import (
	"reflect"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// IsIfaceBsonMarshalSafe returns true if the given item is safe for both marshalling and unmarshalling to/from native bson data types
// -- ie; implements both bson.ValueMarshaler and bson.ValueUnmarshaler
func IsIfaceBsonMarshalSafe(obj interface{}) bool {
	// log.Trace("IsBsonMarshalSafe()")
	// log.Fatal("util.IsBsonMarshalSafe() -- you lie!")
	// return true
	// return false

	objValue := reflect.ValueOf(obj)
	objValueIndirect := reflect.Indirect(objValue)
	// log.Warn(objValue)
	// log.Warn(objValueIndirect)

	_, objectID := objValueIndirect.Interface().(primitive.ObjectID)
	_, marshaler := objValueIndirect.Interface().(bsoncodec.Marshaler)
	_, valueMarshaler := objValueIndirect.Interface().(bsoncodec.ValueMarshaler)
	if valueMarshaler || marshaler || objectID {
		// log.Trace("IsBsonMarshalSafe() - true")
		return true
	}
	// log.Trace("IsBsonMarshalSafe() - false")
	return false

	// implements mongo-go-driver/bson.ValueMKarshaler
	// implements mongo-go-driver/bson.ValueUnmarshaler

	// test ptrs to objects as well as concrete objs, favor the type we were given if both types exist

	// if is ptr to a type...
	// if true { // is ptr interface
	// } else { // direct interface
	// }

	// return false // make compiler happy
}

// objPtrInterfaceFromObjInterface := func(){}
// objInterfaceFromObjPtrInterface := func(){}


// ValidateBson returns a list of bson keys holding invalid content
func ValidateBson(bsonM bson.M) []string {
	return findInvalidBsonM(bsonM, "")
}


// returns a list of bsonM keys holding invalid content (ie, complex types such as non-nil pointers, structs, non-bsonA slices or arrays, funcs, channels, interfaces, etc)
func findInvalidBsonM(bsonM bson.M, prepend string) []string {
	// fmt.Println("findInvalidBsonM()")
	badKeys := make([]string, 0)

	// fmt.Println("bsonM:", bsonM)

	for key, value := range bsonM {
		if value == nil { // skip nil interface values
			continue
		}

		valueValue := reflect.ValueOf(value)

		// native BSON marshallable object types are given a pass
		if IsIfaceBsonMarshalSafe(valueValue.Interface()) {
			continue
		}

		// fmt.Println("key: ", key, " = ", valueValue.Type())

		if childBsonM, ok := value.(bson.M); ok {
			// fmt.Println("bsonM")
			badKeys = append(badKeys, findInvalidBsonM(childBsonM, prepend+"."+key)...)
			continue
		}

		if childBsonA, ok := value.(bson.A); ok {
			// fmt.Println("bsonA")
			badKeys = append(badKeys, findInvalidBsonA(childBsonA, prepend+"."+key)...)
			continue
		}

		switch valueValue.Kind() {
		case reflect.Ptr:
			fallthrough
		case reflect.Invalid:
			fallthrough
		case reflect.Array:
			fallthrough
		case reflect.Chan:
			fallthrough
		case reflect.Func:
			fallthrough
		case reflect.Interface:
			fallthrough
		case reflect.Slice:
			fallthrough
		case reflect.Map:
			fallthrough
		case reflect.Struct:
			fallthrough
		case reflect.UnsafePointer:
			badKeys = append(badKeys, fmt.Sprintf("%s.%s(%s)", prepend, key, valueValue.Kind()))
		}
	}

	// badKeys = append(badKeys, "trash")
	return badKeys
}

func findInvalidBsonA(bsonA bson.A, prepend string) []string {
	// fmt.Println("findInvalidBsonA()")
	badIndexes := make([]string, 0)

	// fmt.Println("bsonA:", bsonA)
	for index, value := range bsonA {
		if value == nil { // skip nil interface values
			continue
		}
		valueValue := reflect.ValueOf(value)
		// fmt.Println("index: ", index, " = ", valueValue.Type())

		// native BSON marshallable object types are given a pass
		if IsIfaceBsonMarshalSafe(valueValue.Interface()) {
			continue
		}

		if childBsonM, ok := value.(bson.M); ok {
			// fmt.Println("bsonM")
			badIndexes = append(badIndexes, findInvalidBsonM(childBsonM, fmt.Sprintf("%s[%d]", prepend, index))...)
			continue
		}

		if childBsonA, ok := value.(bson.A); ok {
			// fmt.Println("bsonA")
			badIndexes = append(badIndexes, findInvalidBsonA(childBsonA, fmt.Sprintf("%s[%d]", prepend, index))...)
			continue
		}

		switch valueValue.Kind() {
		case reflect.Ptr:
			fallthrough
		case reflect.Invalid:
			fallthrough
		case reflect.Array:
			fallthrough
		case reflect.Chan:
			fallthrough
		case reflect.Func:
			fallthrough
		case reflect.Interface:
			fallthrough
		case reflect.Slice:
			fallthrough
		case reflect.Map:
			fallthrough
		case reflect.Struct:
			fallthrough
		case reflect.UnsafePointer:
			badIndexes = append(badIndexes, fmt.Sprintf("%s[%d](%s)", prepend, index, valueValue.Kind()))
		}
	}

	return badIndexes
}
