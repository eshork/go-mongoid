package errors

/*
	Error objects relating to document methods
*/

// import (
// 	"mongoid/log"
// 	"reflect"

// 	"github.com/iancoleman/strcase"
// )

// FieldNotFound -
type FieldNotFound struct {
	FieldName string
}

func (err *FieldNotFound) Error() string {
	if err.FieldName != "" {
		return "FieldNotFound: " + err.FieldName
	}
	return "FieldNotFound"
}
