package mongoid

import (
	"reflect"
)

// returns true if both left and right are samesame, but different, but still same (both are 0pointers to same type, but maybe not to same concrete object)
// if you look closely, this is at least partly stolen from: https://github.com/stretchr/testify/blob/34c6fa2dc70986bccbbffcc6130f6920a924b075/assert/assertions.go#L359
func verifyBothAreSameSame(left interface{}, right interface{}) bool {
	leftPtr, rightPtr := reflect.ValueOf(left), reflect.ValueOf(right)
	if leftPtr.Kind() != reflect.Ptr || rightPtr.Kind() != reflect.Ptr {
		return false // both arguments must be pointers
	}
	leftType, rightType := reflect.TypeOf(left), reflect.TypeOf(right)
	if leftType != rightType {
		return false // type mismatch
	}
	return true
}
