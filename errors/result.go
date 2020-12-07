package errors

/*
	Errors relating to Result objects
*/

// import (
// 	"mongoid/log"
// 	"reflect"

// 	"github.com/iancoleman/strcase"
// )

// // FieldNotFound -
// type FieldNotFound struct {
// 	FieldName string
// }

// func (err *FieldNotFound) Error() string {
// 	if err.FieldName != "" {
// 		return "FieldNotFound: " + err.FieldName
// 	}
// 	return "FieldNotFound"
// }

// InvalidMethodCall can be raised when a method is called in an unexpected manner (out of operations order, etc)
type InvalidMethodCall struct {
	MongoidError
	MethodName string
	Reason     string
}
