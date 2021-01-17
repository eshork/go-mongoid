package util

import (
	"strconv"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MarshalToDB converts the given fromValue into the appropriate database storage type, returning an interface to the database-ready value
// along with a true boolean value to indicate success.
// It is capable of marshalling all built-in types, plus primitive.ObjectID, as well as pointers to those types.
// Some types of fromValue may result in an altogether different return type; ex: uint64 is converted into a string.
// If the type of fromValue does not require a conversion, the given fromValue interface may be returned directly, but this behavior is not guaranteed.
// If the given fromValue is a type that MarshalToDB is unable to marshal, the given fromValue will be returned directly along with a false boolean value
// to indicate failure.
func MarshalToDB(fromValue interface{}) (interface{}, bool) {
	switch fromValue.(type) {
	case primitive.ObjectID:
		return fromValue, true
	case bool:
		return fromValue, true
	case string:
		return fromValue, true
	case int:
		return int32(fromValue.(int)), true
	case int8:
		return int32(fromValue.(int8)), true
	case int16:
		return int32(fromValue.(int16)), true
	case int32: // also covers rune
		return int32(fromValue.(int32)), true
	case int64:
		return int64(fromValue.(int64)), true
	case uint:
		return int64(fromValue.(uint)), true
	case uint8: // also covers byte
		return int32(fromValue.(uint8)), true
	case uint16:
		return int32(fromValue.(uint16)), true
	case uint32:
		return int64(fromValue.(uint32)), true
	case uint64:
		val := fromValue.(uint64)
		return strconv.FormatUint(val, 10), true
	case float32:
		return float64(fromValue.(float32)), true
	case float64:
		return float64(fromValue.(float64)), true
	case complex64:
		return strconv.FormatComplex(complex128(fromValue.(complex64)), 'f', -1, 64), true
	case complex128:
		return strconv.FormatComplex(fromValue.(complex128), 'f', -1, 128), true
	case *bool:
		if fromValue != nil {
			return MarshalToDB(*(fromValue.(*bool)))
		}
		return nil, true
	case *string:
		if fromValue != nil {
			return MarshalToDB(*(fromValue.(*string)))
		}
		return nil, true
	case *int:
		if fromValue != nil {
			return MarshalToDB(*(fromValue.(*int)))
		}
		return nil, true
	case *int8:
		if fromValue != nil {
			return MarshalToDB(*(fromValue.(*int8)))
		}
		return nil, true
	case *int16:
		if fromValue != nil {
			return MarshalToDB(*(fromValue.(*int16)))
		}
		return nil, true
	case *int32:
		if fromValue != nil {
			return MarshalToDB(*(fromValue.(*int32)))
		}
		return nil, true
	case *int64:
		if fromValue != nil {
			return MarshalToDB(*(fromValue.(*int64)))
		}
		return nil, true
	case *uint:
		if fromValue != nil {
			return MarshalToDB(*(fromValue.(*uint)))
		}
		return nil, true
	case *uint8:
		if fromValue != nil {
			return MarshalToDB(*(fromValue.(*uint8)))
		}
		return nil, true
	case *uint16:
		if fromValue != nil {
			return MarshalToDB(*(fromValue.(*uint16)))
		}
		return nil, true
	case *uint32:
		if fromValue != nil {
			return MarshalToDB(*(fromValue.(*uint32)))
		}
		return nil, true
	case *uint64:
		if fromValue != nil {
			return MarshalToDB(*(fromValue.(*uint64)))
		}
		return nil, true
	case *float32:
		if fromValue != nil {
			return MarshalToDB(*(fromValue.(*float32)))
		}
		return nil, true
	case *float64:
		if fromValue != nil {
			return MarshalToDB(*(fromValue.(*float64)))
		}
		return nil, true
	case *complex64:
		if fromValue != nil {
			return MarshalToDB(*(fromValue.(*complex64)))
		}
		return nil, true
	case *complex128:
		if fromValue != nil {
			return MarshalToDB(*(fromValue.(*complex128)))
		}
		return nil, true
	default:
		return fromValue, false
	}
}
