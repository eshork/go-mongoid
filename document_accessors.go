package mongoid

/*
	IDocumentBase implementations relating to generic field accessors: GetField(), SetField(), GetFieldWas()
*/

import (
	mongoidError "mongoid/errors"
	"mongoid/log"
	"reflect"

	"github.com/iancoleman/strcase"
)

// GetFieldPrevious returns an interface to the previous value from the document found at the given fieldNamePath and a true boolean if the path was valid
func (d *Base) GetFieldPrevious(fieldNamePath string) (interface{}, bool) {
	log.Tracef("GetFieldPrevious(%s)", fieldNamePath)
	// this doesn't chase into sub-documents, will probably want to fix that one day
	if value, found := d.previousValue[fieldNamePath]; found {
		return value, true
	}
	return nil, false
}

// SetField sets a value on the document via bson field name path
func (d *Base) SetField(fieldNamePath string, newValue interface{}) error {
	log.Debugf("%v.SetField(%s)", d.Model().modelName, fieldNamePath)
	// get a Value handle to the field we want
	found, retVal, _ := getStructFieldValueRefByBsonPath(d.DocumentBase(), fieldNamePath)
	if found { // if we find the field, assign the value
		newValue := reflect.ValueOf(newValue)
		retVal.Set(newValue)
		return nil
	}
	return &mongoidError.DocumentFieldNotFound{FieldName: fieldNamePath}
}

// GetField returns an interface to a value from the document via the bson field name path
func (d *Base) GetField(fieldNamePath string) (interface{}, error) {
	log.Tracef("GetField(%s)", fieldNamePath)
	// get a Value handle to the field we want
	found, retVal, _ := getStructFieldValueRefByBsonPath(d.DocumentBase(), fieldNamePath)
	if found { // if we find the field, return an interface to the value
		return retVal.Interface(), nil
	}
	return nil, &mongoidError.DocumentFieldNotFound{FieldName: fieldNamePath}
}

func getStructFieldValueRefByBsonName(rawStructPtr interface{}, fieldName string) (found bool, retVal reflect.Value, retField reflect.StructField) {
	handleType := reflect.TypeOf(rawStructPtr)
	if handleType.Kind() != reflect.Ptr { // ensure handleType is always a pointer to a struct
		handleType = reflect.PtrTo(handleType)
	}

	handleStructType := handleType.Elem()
	// log.Printf("getStructFieldValueRefByBson(%+v)", handleStructType) // TODO CLEANUP

	handleValue := reflect.Indirect(reflect.ValueOf(rawStructPtr))
	// log.Printf("getStructFieldValueRefByBson INPUT: %+v", handleValue) // TODO CLEANUP

	// walk each struct field
	for i := 0; i < handleStructType.NumField(); i++ {
		structField := handleStructType.Field(i) // Get the field type, returns https://golang.org/pkg/reflect/#StructField
		if structField.PkgPath != "" {           // skip non-exported fields
			continue
		}

		// collect field information
		tagFieldName, _, _, tagInline := getBsonStructTagOpts(structField)

		if tagInline { // inlined child-struct fields need special handling (recursion following)
			log.Panic("CANNOT SIMPLY SKIP INLINED")
			continue
		}

		// process field name
		thisFieldName := structField.Name
		switch tagFieldName {
		case "": // empty tagFieldName means no explicit field name substitution was given, so snake_case the existing struct field name
			thisFieldName = strcase.ToSnake(thisFieldName)
		case "-": // "-" indicates this field is explicitly omitted, so we can skip this field
			continue
		default: // otherwise use tagFieldName as field name
			thisFieldName = tagFieldName
		}
		// log.Error(thisFieldName)
		// log.Error(structField.Name)
		// log.Error(structField.PkgPath)

		if thisFieldName == fieldName { // if this is the field we're looking for ...
			retVal = handleValue.Field(i) // Get the field value, returns https://golang.org/pkg/reflect/#Value
			return true, retVal, structField
		}
	}

	return
}

// like getStructFieldValueRefByBson but follows a bson path
func getStructFieldValueRefByBsonPath(rawStructPtr interface{}, fieldNamePath string) (found bool, retVal reflect.Value, retField reflect.StructField) {
	// walk each struct field
	return getStructFieldValueRefByBsonName(rawStructPtr, fieldNamePath) // TODO FINISH ME
	log.Fatalf("getStructFieldValueRefByBsonPath(%s)", fieldNamePath)
	return
}
