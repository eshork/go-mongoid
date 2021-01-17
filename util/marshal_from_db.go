package util

import (
	"reflect"
	"strconv"

	"mongoid/log"
)

// MarshalFromDB casts the given fromValue into the given intoType according to expected DB value conversions, returning an interface to the newly cast value.
// If fromValue is already the type of intoType, it may be returned directly, but it is not guaranteed to do so.
// If a value conversion would result in loss of data or precision, this function will panic.
func MarshalFromDB(intoType reflect.Type, fromValue interface{}) interface{} {
	if reflect.TypeOf(fromValue) == intoType {
		return fromValue
	}

	switch intoType.Kind() {
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		fallthrough
	case reflect.Int:
		dstPtr := reflect.New(intoType)
		dst := reflect.Indirect(dstPtr)
		src := reflect.ValueOf(fromValue)
		if dst.OverflowInt(src.Int()) {
			log.Panicf("Overflow detected while storing %v within %v", src.Type(), dst.Type())
		}
		dst.SetInt(src.Int())
		return dst.Interface()
	case reflect.Uint8:
		fallthrough
	case reflect.Uint16:
		fallthrough
	case reflect.Uint32:
		fallthrough
	case reflect.Uint64:
		fallthrough
	case reflect.Uint:
		dstPtr := reflect.New(intoType)
		dst := reflect.Indirect(dstPtr)
		var srcStr string
		switch fromValue.(type) {
		case int64:
			srcStr = strconv.FormatInt(fromValue.(int64), 10)
		case int32:
			srcStr = strconv.FormatInt(int64(fromValue.(int32)), 10)
		case string:
			srcStr = fromValue.(string)
		}
		srcUint64, srcUint64Err := strconv.ParseUint(srcStr, 10, 64)
		if srcUint64Err != nil {
			log.Panicf("Error detected while storing %v within %v: %v", reflect.TypeOf(fromValue), intoType, srcUint64Err)
		}
		if dst.OverflowUint(srcUint64) {
			log.Panicf("Overflow detected while storing %v within %v", reflect.TypeOf(fromValue), intoType)
		}
		dst.SetUint(srcUint64)
		return dst.Interface()
	case reflect.Complex64:
		fallthrough
	case reflect.Complex128:
		dstPtr := reflect.New(intoType)
		dst := reflect.Indirect(dstPtr)
		var srcStr string
		switch fromValue.(type) {
		case complex64:
			srcStr = strconv.FormatComplex(complex128(fromValue.(complex64)), 'f', -1, 64)
		case complex128:
			srcStr = strconv.FormatComplex(fromValue.(complex128), 'f', -1, 128)
		case string:
			srcStr = fromValue.(string)
		}
		var dstBits int
		switch intoType.Kind() {
		case reflect.Complex64:
			dstBits = 64
		case reflect.Complex128:
			dstBits = 128
		}
		srcComplex128, srcComplex128Err := strconv.ParseComplex(srcStr, dstBits)
		if srcComplex128Err != nil {
			log.Panicf("Error detected while storing %v within %v: %v", reflect.TypeOf(fromValue), intoType, srcComplex128Err)
		}
		if dst.OverflowComplex(srcComplex128) {
			log.Panicf("Overflow detected while storing %v within %v", reflect.TypeOf(fromValue), intoType)
		}
		dst.SetComplex(srcComplex128)
		return dst.Interface()
	}
	log.Panicf("Unhandled kind: %v", intoType.Kind())
	return nil
}
