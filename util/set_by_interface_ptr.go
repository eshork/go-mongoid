package util

import (
	"mongoid/log"
	"reflect"
)

// SetValueByInterfacePtr sets the value pointed to by dstPtr to the value held by srcValue, without having to know the types beforehand.
// The given dstPtr must be either a pointer to a builtin type or struct, or a pointer to a pointer (ex **int).
// If dstPtr is given a pointer to a builtin type or struct (ex *int), the value pointer to by dstPtr will be set to the
// value held by srcValue via assignment operator.
// If dstPtr is a pointer to a pointer, a new object of dstPtr's type will be created to recieve the given srcValue,
// unless srcValue itself is nil or another pointer, in which case the pointer that dstPtr points to will be directly
// set to nil or the given pointer value (without a new object creation).
func SetValueByInterfacePtr(dstPtr interface{}, srcValue interface{}) {
	// This has been implemented a few different ways over time, some being heavily dependent on the reflect package.
	// While quite long, using a large type-switch with obvious type casts and object creation seems to be the most understandable (and hopefully maintainable).
	// Overall, the logic for handling a pointer to a value and the logic for handling a pointer to a pointer to a value is the same for each basic type.
	// So there is, obviously, a lot of repetition.
	switch dstPtr.(type) {
	case string:
		log.Panic("dstPtr was string, expected *string or **string")
	case *string:
		*(dstPtr.(*string)) = srcValue.(string)
	case **string:
		if srcValue == nil {
			*(dstPtr.(**string)) = nil
		} else if _, ok := srcValue.(*string); ok {
			*(dstPtr.(**string)) = srcValue.(*string)
		} else {
			newValue := new(string)
			*newValue = srcValue.(string)
			*(dstPtr.(**string)) = newValue
		}
	case bool:
		log.Panic("dstPtr was bool, expected *bool or **bool")
	case *bool:
		*(dstPtr.(*bool)) = srcValue.(bool)
	case **bool:
		if srcValue == nil {
			*(dstPtr.(**bool)) = nil
		} else if _, ok := srcValue.(*bool); ok {
			*(dstPtr.(**bool)) = srcValue.(*bool)
		} else {
			newValue := new(bool)
			*newValue = srcValue.(bool)
			*(dstPtr.(**bool)) = newValue
		}
	case int:
		log.Panic("dstPtr was int, expected *int or **int")
	case *int:
		*(dstPtr.(*int)) = srcValue.(int)
	case **int:
		if srcValue == nil {
			*(dstPtr.(**int)) = nil
		} else if _, ok := srcValue.(*int); ok {
			*(dstPtr.(**int)) = srcValue.(*int)
		} else {
			newValue := new(int)
			*newValue = srcValue.(int)
			*(dstPtr.(**int)) = newValue
		}
	case int8:
		log.Panic("dstPtr was int8, expected *int8 or **int8")
	case *int8:
		*(dstPtr.(*int8)) = srcValue.(int8)
	case **int8:
		if srcValue == nil {
			*(dstPtr.(**int8)) = nil
		} else if _, ok := srcValue.(*int8); ok {
			*(dstPtr.(**int8)) = srcValue.(*int8)
		} else {
			newValue := new(int8)
			*newValue = srcValue.(int8)
			*(dstPtr.(**int8)) = newValue
		}
	case int16:
		log.Panic("dstPtr was int8, expected *int8 or **int8")
	case *int16:
		*(dstPtr.(*int16)) = srcValue.(int16)
	case **int16:
		if srcValue == nil {
			*(dstPtr.(**int16)) = nil
		} else if _, ok := srcValue.(*int16); ok {
			*(dstPtr.(**int16)) = srcValue.(*int16)
		} else {
			newValue := new(int16)
			*newValue = srcValue.(int16)
			*(dstPtr.(**int16)) = newValue
		}
	case int32:
		log.Panic("dstPtr was int8, expected *int8 or **int8")
	case *int32:
		*(dstPtr.(*int32)) = srcValue.(int32)
	case **int32:
		if srcValue == nil {
			*(dstPtr.(**int32)) = nil
		} else if _, ok := srcValue.(*int32); ok {
			*(dstPtr.(**int32)) = srcValue.(*int32)
		} else {
			newValue := new(int32)
			*newValue = srcValue.(int32)
			*(dstPtr.(**int32)) = newValue
		}
	case int64:
		log.Panic("dstPtr was int8, expected *int8 or **int8")
	case *int64:
		*(dstPtr.(*int64)) = srcValue.(int64)
	case **int64:
		if srcValue == nil {
			*(dstPtr.(**int64)) = nil
		} else if _, ok := srcValue.(*int64); ok {
			*(dstPtr.(**int64)) = srcValue.(*int64)
		} else {
			newValue := new(int64)
			*newValue = srcValue.(int64)
			*(dstPtr.(**int64)) = newValue
		}
	case uint:
		log.Panic("dstPtr was uint, expected *uint or **uint")
	case *uint:
		*(dstPtr.(*uint)) = srcValue.(uint)
	case **uint:
		if srcValue == nil {
			*(dstPtr.(**uint)) = nil
		} else if _, ok := srcValue.(*uint); ok {
			*(dstPtr.(**uint)) = srcValue.(*uint)
		} else {
			newValue := new(uint)
			*newValue = srcValue.(uint)
			*(dstPtr.(**uint)) = newValue
		}
	case uint8:
		log.Panic("dstPtr was uint8, expected *uint8 or **uint8")
	case *uint8:
		*(dstPtr.(*uint8)) = srcValue.(uint8)
	case **uint8:
		if srcValue == nil {
			*(dstPtr.(**uint8)) = nil
		} else if _, ok := srcValue.(*uint8); ok {
			*(dstPtr.(**uint8)) = srcValue.(*uint8)
		} else {
			newValue := new(uint8)
			*newValue = srcValue.(uint8)
			*(dstPtr.(**uint8)) = newValue
		}
	case uint16:
		log.Panic("dstPtr was uint16, expected *uint16 or **uint16")
	case *uint16:
		*(dstPtr.(*uint16)) = srcValue.(uint16)
	case **uint16:
		if srcValue == nil {
			*(dstPtr.(**uint16)) = nil
		} else if _, ok := srcValue.(*uint16); ok {
			*(dstPtr.(**uint16)) = srcValue.(*uint16)
		} else {
			newValue := new(uint16)
			*newValue = srcValue.(uint16)
			*(dstPtr.(**uint16)) = newValue
		}
	case uint32:
		log.Panic("dstPtr was uint32, expected *uint32 or **uint32")
	case *uint32:
		*(dstPtr.(*uint32)) = srcValue.(uint32)
	case **uint32:
		if srcValue == nil {
			*(dstPtr.(**uint32)) = nil
		} else if _, ok := srcValue.(*uint32); ok {
			*(dstPtr.(**uint32)) = srcValue.(*uint32)
		} else {
			newValue := new(uint32)
			*newValue = srcValue.(uint32)
			*(dstPtr.(**uint32)) = newValue
		}
	case uint64:
		log.Panic("dstPtr was uint64, expected *uint64 or **uint64")
	case *uint64:
		*(dstPtr.(*uint64)) = srcValue.(uint64)
	case **uint64:
		if srcValue == nil {
			*(dstPtr.(**uint64)) = nil
		} else if _, ok := srcValue.(*uint64); ok {
			*(dstPtr.(**uint64)) = srcValue.(*uint64)
		} else {
			newValue := new(uint64)
			*newValue = srcValue.(uint64)
			*(dstPtr.(**uint64)) = newValue
		}
	case float32:
		log.Panic("dstPtr was float32, expected *float32 or **float32")
	case *float32:
		*(dstPtr.(*float32)) = srcValue.(float32)
	case **float32:
		if srcValue == nil {
			*(dstPtr.(**float32)) = nil
		} else if _, ok := srcValue.(*float32); ok {
			*(dstPtr.(**float32)) = srcValue.(*float32)
		} else {
			newValue := new(float32)
			*newValue = srcValue.(float32)
			*(dstPtr.(**float32)) = newValue
		}
	case float64:
		log.Panic("dstPtr was float64, expected *float64 or **float64")
	case *float64:
		*(dstPtr.(*float64)) = srcValue.(float64)
	case **float64:
		if srcValue == nil {
			*(dstPtr.(**float64)) = nil
		} else if _, ok := srcValue.(*float64); ok {
			*(dstPtr.(**float64)) = srcValue.(*float64)
		} else {
			newValue := new(float64)
			*newValue = srcValue.(float64)
			*(dstPtr.(**float64)) = newValue
		}
	case complex64:
		log.Panic("dstPtr was complex64, expected *complex64 or **complex64")
	case *complex64:
		*(dstPtr.(*complex64)) = srcValue.(complex64)
	case **complex64:
		if srcValue == nil {
			*(dstPtr.(**complex64)) = nil
		} else if _, ok := srcValue.(*complex64); ok {
			*(dstPtr.(**complex64)) = srcValue.(*complex64)
		} else {
			newValue := new(complex64)
			*newValue = srcValue.(complex64)
			*(dstPtr.(**complex64)) = newValue
		}
	case complex128:
		log.Panic("dstPtr was complex128, expected *complex128 or **complex128")
	case *complex128:
		*(dstPtr.(*complex128)) = srcValue.(complex128)
	case **complex128:
		if srcValue == nil {
			*(dstPtr.(**complex128)) = nil
		} else if _, ok := srcValue.(*complex128); ok {
			*(dstPtr.(**complex128)) = srcValue.(*complex128)
		} else {
			newValue := new(complex128)
			*newValue = srcValue.(complex128)
			*(dstPtr.(**complex128)) = newValue
		}
	default:
		// we've exhausted the statically identifiable types, so we'll use reflect to handle custom Struct, Map, and Array (Slice) types
		dst := reflect.Indirect(reflect.ValueOf(dstPtr))
		switch dst.Kind() {
		case reflect.Struct:
			fallthrough
		case reflect.Array:
			fallthrough
		case reflect.Slice:
			fallthrough
		case reflect.Map:
			dst.Set(reflect.ValueOf(srcValue))
		case reflect.Ptr:
			if srcValue == nil {
				dst.Set(reflect.Zero(dst.Type())) // set the pointer to an appropriate version of nil
				return
			}

			// if srcValue is also a pointer, we can simply do direct assignment
			if srcValueRef := reflect.ValueOf(srcValue); srcValueRef.Kind() == reflect.Ptr {
				dst.Set(srcValueRef)
				return
			}

			// otherwise make a new object and copy over the value and we'll store a ptr to that new object
			newValuePtr := reflect.New(reflect.TypeOf(srcValue))
			reflect.Indirect(newValuePtr).Set(reflect.ValueOf(srcValue))
			dst.Set(newValuePtr)

		default:
			log.Panic("setByInterfacePtr() - dstPtr is an unhandled type")
		}
	}
}
