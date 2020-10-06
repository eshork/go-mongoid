package mongoid

import (
	"fmt"
	"mongoid/log"
	"reflect"
	"strings"
)

// ModelType represents a mongoid model/document type and provides methods to interact with the collection
type ModelType struct {
	fmt.Stringer                 // implements Stringer interface
	rootTypeRef    IDocumentBase // a reference to an object instance of the document/model type given during registration (for future type sanity)
	modelName      string
	modelFullName  string
	collectionName string
	databaseName   string
	clientName     string
	defaultValue   BsonDocument // bson representation of default values to be applied during creation of brand new document/model instances
}

func (model ModelType) String() string {
	// return "pew pew pew~!" // ModelType
	extras := ""

	// [clientName, databaseName, collectionName]
	if model.collectionName != "" {
		extras += "collection:" + model.collectionName + ","
	}
	if model.databaseName != "" {
		extras += "database:" + model.databaseName + ","
	}
	if model.clientName != "" {
		extras += "client:" + model.clientName + ","
	}
	if len(extras) > 0 {
		extras = " [" + extras[:len(extras)-1] + "]"
	}
	return fmt.Sprintf("%s (%s)%s", model.modelName, model.modelFullName, extras)
}

// reflect-ive rootTypeRef type equality check
func (model *ModelType) equalsRootType(comparisonModel *ModelType) bool {
	if model != nil && comparisonModel != nil && model.rootTypeRef != nil && comparisonModel.rootTypeRef != nil {
		if reflect.TypeOf(model.rootTypeRef) == reflect.TypeOf(comparisonModel.rootTypeRef) {
			if reflect.TypeOf(model.rootTypeRef).Kind() == reflect.Ptr {
				if reflect.TypeOf(model.rootTypeRef).Elem() == reflect.TypeOf(comparisonModel.rootTypeRef).Elem() {
					return true
				}
				return false
			}
			return true
		}
	}
	return false
}

// SetModelName changes the model name used by the ModelType
func (model *ModelType) SetModelName(newModelName string) *ModelType {
	newModelType := *model // dereferenced copy
	newModelType.modelName = newModelName
	// update the global registry for this ModelType
	return mongoidModelRegistry.updateModelTypeRegistration(&newModelType)
}

// GetModelName returns the current friendly name for this model type
func (model *ModelType) GetModelName() string {
	return model.modelName
}

// SetCollectionName changes the collection name used by the ModelType
func (model *ModelType) SetCollectionName(newCollectionName string) *ModelType {
	newModelType := *model // dereferenced copy
	newModelType.collectionName = newCollectionName
	// update the global registry for this ModelType
	return mongoidModelRegistry.updateModelTypeRegistration(&newModelType)
}

// GetCollectionName returns the current default collection name for this model type
func (model *ModelType) GetCollectionName() string {
	return model.collectionName
}

// SetDatabaseName changes the database name used by the ModelType
func (model *ModelType) SetDatabaseName(newDatabaseName string) *ModelType {
	newModelType := *model // dereferenced copy
	newModelType.databaseName = newDatabaseName
	// update the global registry for this ModelType
	return mongoidModelRegistry.updateModelTypeRegistration(&newModelType)
}

// GetDatabaseName returns the current default database for this model type
func (model *ModelType) GetDatabaseName() string {
	if model.databaseName == "" {
		return model.GetClient().Database
	}
	return model.databaseName
}

// SetClientName changes the client name used by the ModelType
func (model *ModelType) SetClientName(newClientName string) *ModelType {
	newModelType := *model // dereferenced copy
	newModelType.clientName = newClientName
	// update the global registry for this ModelType
	return mongoidModelRegistry.updateModelTypeRegistration(&newModelType)
}

// GetClientName returns the current default client name for this model type
func (model *ModelType) GetClientName() string {
	return model.clientName
}

// GetClient returns the current default client for this model type
func (model *ModelType) GetClient() *Client {
	if clientName := model.GetClientName(); clientName != "" {
		return ClientByName(clientName)
	}
	return DefaultClient()
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// New intantiates a new document model object of the registered type and returns a pointer to the new object.
// The returned object will be preset with the defaults specified during initial document/model registration.
// Note: Due to the strongly typed nature of Go, you'll need to perform a type assertion (as the value is returned as an interface{})
func (model *ModelType) New() IDocumentBase {
	log.Debug("ModelType.New()")
	typeRef := reflect.Indirect(reflect.ValueOf(model.rootTypeRef))                 // model.rootTypeRef is always a ptr to an example object, so we need to use Indirect()
	ret := reflect.New(typeRef.Type())                                              // finally have a solid object type, so make one
	retAsIDocumentBase := ret.Interface().(IDocumentBase)                           // convert into a IDocumentBase interface
	retAsIDocumentBase.initDocumentBase(retAsIDocumentBase, model.GetDefaultBSON()) // call the self init

	// apply the default values and reset "changed" state

	return retAsIDocumentBase // return the new object as an IDocumentBase interface
}

// GetDefaultBSON provides the default values for a ModelType returned as a BsonDocument.
// The returned value is deep-cloned to protect the original data, so you can begin using it directly without a second deep copy
func (model *ModelType) GetDefaultBSON() BsonDocument {
	log.Trace("GetDefaultBSON()")
	return BsonDocumentDeepCopy(model.defaultValue)
}

// NYI - ref: https://github.com/eshork/go-mongoid/issues/17
// func (model *ModelType) SetDefaultScope() {
// 	// log.Panic("NYI")
// 	log.Error("NYI - ModelType.SetDefaultScope")
// }

// NYI - ref: https://github.com/eshork/go-mongoid/issues/18
// func (model *ModelType) AddIndex() {
// 	// log.Panic("NYI")
// 	log.Error("NYI - ModelType.AddIndex")
// }

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

type modelTypeTagOpts struct {
	modelName      string
	collectionName string
	databaseName   string
	clientName     string
}

// parses a string of comma-separated modelTypeTagOpts key:value pairs and returns a new modelTypeTagOpts
func modelTypeTagOptsFromString(tagString string) modelTypeTagOpts {
	ret := modelTypeTagOpts{}
	// separate by commas
	segments := strings.Split(tagString, ",")
	for _, v := range segments {
		segment := strings.TrimSpace(v)
		// separate segments by colons to find key/value pairs
		pair := strings.Split(segment, ":")
		if len(pair) == 2 && len(pair[0]) > 0 && len(pair[1]) > 0 {
			switch strings.TrimSpace(pair[0]) {
			case "model":
				ret.modelName = strings.TrimSpace(pair[1])
			case "collection":
				ret.collectionName = strings.TrimSpace(pair[1])
			case "database":
				ret.databaseName = strings.TrimSpace(pair[1])
			case "client":
				ret.clientName = strings.TrimSpace(pair[1])
			}
		}
	}
	return ret
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// func AddCallback(when string, what func(*Base)) {}
// func AddValidation(when string, what func(*Base)(error)) {}
// func AddIndex()

// type Fields map[string]interface{}

// CreateCollection creates a collection on the Client's connected topology with the given databaseName and collectionName pair
// func (c *Client) CreateCollection(databaseName string, collectionName string) error {
// 	return nil
// }
