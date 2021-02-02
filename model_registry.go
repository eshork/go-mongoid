package mongoid

import (
	"fmt"
	mongoidErr "mongoid/errors"
	"mongoid/log"
	"reflect"
	"strings"
	"sync"

	"github.com/iancoleman/strcase"
	"github.com/jinzhu/inflection"
)

// the Config data holders & mutex
var mongoidModelRegistry *modelRegistry  // the global modelRegistry
var mongoidModelRegistryMutex sync.Mutex // the global modelRegistry mutex, used to synchronize access

// type typeDocumentBaseMap map[string]IDocumentBase
type typeModelTypeMapByName map[string]ModelType

type modelRegistry struct {
	modelTypeMap typeModelTypeMapByName
}

func init() {
	mongoidModelRegistryMutex.Lock()
	defer mongoidModelRegistryMutex.Unlock()
	if mongoidModelRegistry == nil {
		mongoidModelRegistry = new(modelRegistry)
		mongoidModelRegistry.modelTypeMap = make(typeModelTypeMapByName)
	}
}

// retrieves a ModelType for a previously registered IDocumentBase via the model name
// returns nil if no match was found
func getRegisteredModelTypeByName(modelTypeName string) *ModelType {
	// TODO - make read lock vs write lock more efficient (ie exclusive and non-exclusive locks)
	mongoidModelRegistryMutex.Lock()
	defer mongoidModelRegistryMutex.Unlock()
	for i, v := range mongoidModelRegistry.modelTypeMap {
		if i == modelTypeName {
			return &v
		}
	}
	return nil
}

// retrieves a ModelType for a previously registered IDocumentBase via an example IDocumentBase implementing object ref
func getRegisteredModelTypeByDocRef(modelTypeDocRef IDocumentBase) *ModelType {
	// TODO - make read lock vs write lock more efficient (ie exclusive and non-exclusive locks)
	mongoidModelRegistryMutex.Lock()
	defer mongoidModelRegistryMutex.Unlock()
	for _, v := range mongoidModelRegistry.modelTypeMap {
		if verifyBothAreSameSame(modelTypeDocRef, v.rootTypeRef) == true {
			return &v
		}
	}
	return nil
}

// Register the current ModelType model under the current modelName, overwriting a previous registration of the same modelName if necessary.
func (model ModelType) Register() ModelType {
	return mongoidModelRegistry.upsertModel(model)
}

// Register a new documentType
func Register(documentType interface{}) ModelType {

	// special handle ModelType
	if _, ok := documentType.(ModelType); ok {
		log.Panic("DERP")
	}

	// passed as &MyModelStruct{} (and implements IDocumentBase)
	if v, ok := documentType.(IDocumentBase); ok {
		return mongoidModelRegistry.upsert(v)
	}

	// ensure kind is struct
	if documentTypeValue := reflect.ValueOf(documentType); documentTypeValue.Kind() == reflect.Struct {
		// if passed as MyModelStruct{} (and implements IDocumentBase) ...
		newDupeVP := reflect.New(reflect.TypeOf(documentType))
		if v, ok := newDupeVP.Interface().(IDocumentBase); ok {
			newDupeVP.Elem().Set(documentTypeValue)
			return mongoidModelRegistry.upsert(v) // send it onward as the expected struct pointer kind
		}
	}

	// not one of the two things we know how to handle, so have a nice error message
	log.Panic(mongoidErr.InvalidOperation{
		MethodName: "Register",
		Reason:     fmt.Sprintf("documentType must implement IDocumentBase. Found: %s", reflect.ValueOf(documentType).String()),
	})
	return ModelType{} // unreachable
}

func (registry *modelRegistry) upsertModel(modelType ModelType) ModelType {
	mongoidModelRegistryMutex.Lock()
	defer mongoidModelRegistryMutex.Unlock()
	log.Infof("Registered document model: %s", modelType)
	mongoidModelRegistry.modelTypeMap[modelType.modelName] = modelType
	return modelType
}

func (registry *modelRegistry) upsert(documentType IDocumentBase) ModelType {
	modelType := generateModelTypeFromDocument(documentType)
	mongoidModelRegistryMutex.Lock()
	defer mongoidModelRegistryMutex.Unlock()
	log.Infof("Registered document model: %s", modelType)
	mongoidModelRegistry.modelTypeMap[modelType.modelName] = modelType
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
