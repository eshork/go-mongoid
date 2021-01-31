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

// Register a new documentType
func Register(documentType interface{}) *ModelType {
	// passed as &MyModelStruct{} (and implements IDocumentBase)
	if v, ok := documentType.(IDocumentBase); ok {
		return mongoidModelRegistry.register(v)
	}

	// ensure kind is struct
	if documentTypeValue := reflect.ValueOf(documentType); documentTypeValue.Kind() == reflect.Struct {
		// if passed as MyModelStruct{} (and implements IDocumentBase) ...
		newDupeVP := reflect.New(reflect.TypeOf(documentType))
		if v, ok := newDupeVP.Interface().(IDocumentBase); ok {
			newDupeVP.Elem().Set(documentTypeValue)
			return mongoidModelRegistry.register(v) // send it onward as the expected struct pointer kind
		}
	}

	// not one of the two things we know how to handle, so have a nice error message
	log.Panic(mongoidErr.InvalidOperation{
		MethodName: "Register",
		Reason:     fmt.Sprintf("documentType must implement IDocumentBase. Found: %s", reflect.ValueOf(documentType).String()),
	})
	return nil // unreachable
}

func (registry *modelRegistry) register(documentType IDocumentBase) *ModelType {
	log.Trace("ModelRegistry.Register()")

	modelType := generateModelTypeFromDocument(documentType)

	// registry lock
	mongoidModelRegistryMutex.Lock()
	defer mongoidModelRegistryMutex.Unlock()

	if _, ok := mongoidModelRegistry.modelTypeMap[modelType.modelName]; ok {
		log.Fatalf("Cannot register duplicate named model: %s (%s)", modelType.modelName, modelType.modelFullName)
		return nil // unreachable; here to satisfy the compiler
	}

	modelType.defaultValue = structToBsonM(documentType) // build a default value record from the original type reference

	log.Infof("Registered new document model: %s", modelType)
	// log.Infof("Registered new document model: %+v", modelType.defaultValue)
	mongoidModelRegistry.modelTypeMap[modelType.modelName] = modelType
	return &modelType
}

// updates an existing IDocumentBase registration with a new ModelType definition, or fails hard
func (registry *modelRegistry) updateModelTypeRegistration(newModelType *ModelType) *ModelType {
	// warn if configured
	if Configured() == true {
		log.Warnf("Updating existing document model after Configured: %s (%s)", newModelType.modelName, newModelType.modelFullName)
	} else {
		log.Debugf("Updating existing document model: %s (%s)", newModelType.modelName, newModelType.modelFullName)
	}

	// registry lock
	mongoidModelRegistryMutex.Lock()
	defer mongoidModelRegistryMutex.Unlock()

	// find by documentType
	var existingModelType *ModelType
	var existingModelKey string
	for k, v := range mongoidModelRegistry.modelTypeMap {
		if newModelType.equalsRootType(&v) {
			existingModelKey = k
			existingModelType = &v
		}
	}

	// fatal if no type match
	if existingModelType == nil {
		log.Fatalf("No matching registration found for document model: %s (%s)", newModelType.modelName, newModelType.modelFullName)
		return nil
	}

	// drop the existing entry
	delete(mongoidModelRegistry.modelTypeMap, existingModelKey)

	// store the new entry
	log.Infof("Updated existing document model: %s", newModelType)
	mongoidModelRegistry.modelTypeMap[newModelType.modelName] = *newModelType
	return newModelType // return the newly installed entry
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

// TODO this is (likely) dead code. Delete it
func dumpFields(documentType IDocumentBase) {
	handleType := reflect.TypeOf(documentType)
	log.Printf("reflect.Kind: %s", handleType.Kind())
	log.Printf("reflect.Kind: %s", handleType.Elem().Kind())

	handleStructType := reflect.TypeOf(documentType)
	if handleStructType.Kind() == reflect.Ptr {
		handleStructType = handleType.Elem()
	}
	log.Printf("reflect.handleStructType.Kind: %s", handleStructType.Kind())

	for i := 0; i < handleStructType.NumField(); i++ {
		field := handleStructType.Field(i) // Get the field, returns https://golang.org/pkg/reflect/#StructField
		log.Printf("dumpFields Name: %s", field.Name)
		log.Printf("dumpFields Tags: %s", field.Tag)
	}
}
