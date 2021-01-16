package mongoid

import (
	"reflect"
)

// returns true if both left and right are pointers to objects of the same type, regardless of their values
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
