package mongoid

import (
	mongoidErr "mongoid/errors"
	"mongoid/log"
	"reflect"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/jinzhu/inflection"
)

// Model returns a ModelType for the given document
func Model(documentType interface{}) ModelType {
	// function sig uses interface{} param instead of IDocumentBase so that we can accept either a value or reference

	// if not passed as &MyModelStruct{} (that implements IDocumentBase), we need to make our own pointer
	if _, ok := documentType.(IDocumentBase); !ok {
		// ensure object kind is struct
		if documentTypeValue := reflect.ValueOf(documentType); documentTypeValue.Kind() != reflect.Struct {
			log.Panic(mongoidErr.InvalidOperation{
				MethodName: "Model",
				Reason:     "",
			})
		}
		// if passed as MyModelStruct{} (and implements IDocumentBase) ...
		newDupeVP := reflect.New(reflect.TypeOf(documentType))
		if v, ok := newDupeVP.Interface().(IDocumentBase); ok {
			newDupeVP.Elem().Set(reflect.ValueOf(documentType)) // value assignment documentType => (*newDupeVP) to preserve any given default state
			documentType = v
		}
	}
	docType, ok := documentType.(IDocumentBase)
	if !ok {
		log.Panic(mongoidErr.InvalidOperation{
			MethodName: "Model",
			Reason:     "Given struct fails requirements. Must implement the IDocumentBase interface",
		})
	}
	modelType := generateModelTypeFromDocument(docType)
	return modelType
}

func generateModelTypeFromDocument(documentType IDocumentBase) ModelType {
	// start with the defaults
	docTypeNameStr := getDocumentTypeStructName(documentType)
	docTypeFullNameStr := getDocumentTypeFullStructName(documentType)
	modelType := ModelType{
		rootTypeRef:    documentType,
		modelName:      docTypeNameStr,
		modelFullName:  docTypeFullNameStr,
		collectionName: strcase.ToSnake(inflection.Plural(docTypeNameStr)),
		defaultValue:   structToBsonM(documentType), // build a default value record from the original type reference
	}

	// update attributes where overridden by struct tags
	tagOpts := getDocumentTypeOptions(documentType)
	log.Tracef("%s - detected struct tags options: %+v", docTypeFullNameStr, tagOpts)
	if tagOpts.modelName != "" {
		modelType.modelName = tagOpts.modelName
	}
	if tagOpts.collectionName != "" {
		modelType.collectionName = tagOpts.collectionName
	}
	if tagOpts.databaseName != "" {
		modelType.databaseName = tagOpts.databaseName
	}
	if tagOpts.clientName != "" {
		modelType.clientName = tagOpts.clientName
	}
	return modelType
}

// return the original full struct name
func getDocumentTypeFullStructName(documentType IDocumentBase) string {
	handleType := reflect.TypeOf(documentType)
	handleTypeStr := handleType.String()
	if handleTypeStr[:1] == "*" { //drop leading * when present
		handleTypeStr = handleTypeStr[1:]
	}
	return handleTypeStr
}

// return the original struct name
func getDocumentTypeStructName(documentType IDocumentBase) string {
	handleType := reflect.TypeOf(documentType)
	handleTypeStr := handleType.String()
	dotIndex := strings.Index(handleTypeStr, ".")
	if dotIndex > 0 {
		handleTypeStr = handleTypeStr[dotIndex:]
	}
	if handleTypeStr[:1] == "*" { //drop leading * when present
		handleTypeStr = handleTypeStr[1:]
	}
	if handleTypeStr[:1] == "." { //drop leading . when present
		handleTypeStr = handleTypeStr[1:]
	}
	return handleTypeStr
}

// this is the struct tag key name used for mongoid options
const structTagName = "mongoid"

// extract struct tag options for documentType, returning a fully compiled list
func getDocumentTypeOptions(documentType IDocumentBase) modelTypeTagOpts {
	tagOpts := modelTypeTagOpts{}
	// get a handleStructType that always represents the top struct definition
	handleType := reflect.TypeOf(documentType)
	handleStructType := reflect.TypeOf(documentType)
	if handleStructType.Kind() == reflect.Ptr {
		handleStructType = handleType.Elem()
	}
	// walk the struct fields, checking all tags in order, first hit wins
	for i := 0; i < handleStructType.NumField(); i++ {
		field := handleStructType.Field(i) // Get the field, returns https://golang.org/pkg/reflect/#StructField
		if tag, ok := field.Tag.Lookup(structTagName); ok {
			tags := modelTypeTagOptsFromString(tag)
			return tags // first hit wins, this is it!
		}
	}
	return tagOpts
}
