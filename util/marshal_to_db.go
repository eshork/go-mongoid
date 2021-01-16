package util

import (
	"mongoid/log"
	"reflect"
	"strconv"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MarshalToDB converts the given fromValue into the appropriate database storage type, returning an interface to the database-ready value.
// It is capable of marshalling all built-in types, plus primitive.ObjectID, as well as pointers to those types.
// Some types of fromValue may result in an altogether different return type; ex: uint64 is converted into a string.
// If the type of fromValue does not require a conversion, the given fromValue interface may be returned directly, but this behavior is not guaranteed.
// If the given fromValue is a type that MarshalToDB is unable to marshal, the given fromValue will be returned directly - but first we will panic
//   because I'm not yet convinced this is the right behavior.
func MarshalToDB(fromValue interface{}) interface{} {
	switch fromValue.(type) {
	case primitive.ObjectID:
		return fromValue
	case bool:
		return fromValue
	case string:
		return fromValue
	case error:
		return fromValue.(error).Error()
	case int:
		return int32(fromValue.(int))
	case int8:
		return int32(fromValue.(int8))
	case int16:
		return int32(fromValue.(int16))
	case int32: // also covers rune
		return int32(fromValue.(int32))
	case int64:
		return int64(fromValue.(int64))
	case uint:
		return int64(fromValue.(uint))
	case uint8: // also covers byte
		return int32(fromValue.(uint8))
	case uint16:
		return int32(fromValue.(uint16))
	case uint32:
		return int64(fromValue.(uint32))
	case uint64:
		val := fromValue.(uint64)
		return strconv.FormatUint(val, 10)
	case float32:
		return float64(fromValue.(float32))
	case float64:
		return float64(fromValue.(float64))
	case *bool:
		if fromValue != nil {
			return MarshalToDB(*(fromValue.(*bool)))
		}
		return nil
	case *string:
		if fromValue != nil {
			return MarshalToDB(*(fromValue.(*string)))
		}
		return nil
	case *error:
		if fromValue != nil {
			return MarshalToDB(*(fromValue.(*error)))
		}
		return nil
	case *int:
		if fromValue != nil {
			return MarshalToDB(*(fromValue.(*int)))
		}
		return nil
	case *int8:
		if fromValue != nil {
			return MarshalToDB(*(fromValue.(*int8)))
		}
		return nil
	case *int16:
		if fromValue != nil {
			return MarshalToDB(*(fromValue.(*int16)))
		}
		return nil
	case *int32:
		if fromValue != nil {
			return MarshalToDB(*(fromValue.(*int32)))
		}
		return nil
	case *int64:
		if fromValue != nil {
			return MarshalToDB(*(fromValue.(*int64)))
		}
		return nil
	case *uint:
		if fromValue != nil {
			return MarshalToDB(*(fromValue.(*uint)))
		}
		return nil
	case *uint8:
		if fromValue != nil {
			return MarshalToDB(*(fromValue.(*uint8)))
		}
		return nil
	case *uint16:
		if fromValue != nil {
			return MarshalToDB(*(fromValue.(*uint16)))
		}
		return nil
	case *uint32:
		if fromValue != nil {
			return MarshalToDB(*(fromValue.(*uint32)))
		}
		return nil
	case *uint64:
		if fromValue != nil {
			return MarshalToDB(*(fromValue.(*uint64)))
		}
		return nil
	case *float32:
		if fromValue != nil {
			return MarshalToDB(*(fromValue.(*float32)))
		}
		return nil
	case *float64:
		if fromValue != nil {
			return MarshalToDB(*(fromValue.(*float64)))
		}
		return nil
	default:
		log.Panicf("default marshalToDB: %v ", reflect.TypeOf(fromValue))
		return fromValue
	}
}
