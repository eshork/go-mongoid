package mongoid

import (
	"fmt"
	"mongoid/log"
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
)

// ModelType represents a mongoid model/document type and provides methods to interact with records.
// The ModelType scopes document access to a specific collection within a specific database via a specific client.
// Use one of the With... methods to create a new ModelType with an updated scope.
//
type ModelType struct {
	rootTypeRef    IDocumentBase // a reference to an object instance of the document/model type given during registration (for future type sanity)
	modelName      string
	modelFullName  string
	collectionName string
	// Client - determined by combination of databaseName and clientName
	databaseName string //
	clientName   string
	defaultValue BsonDocument // bson representation of default values to be applied during creation of brand new document/model instances
}

var _ fmt.Stringer = ModelType{} // assert implements Stringer interface

// String implements fmt.Stringer interface
func (model ModelType) String() string {
	// ModelType(model.modelFullName)[client:cName,database:dbName,collection:colName]
	extras := ""
	appendExtras := func(name, value string) {
		if extras != "" {
			extras = fmt.Sprintf("%s,%s:%s", extras, name, value)
		}
		extras = fmt.Sprintf("%s:%s", name, value)
	}
	appendExtras("client", model.clientName)
	appendExtras("database", model.databaseName)
	appendExtras("collection", model.collectionName)
	return fmt.Sprintf("ModelType(%s/%s)[%s]", model.modelName, model.modelFullName, extras)
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

// WithModelName returns a ModelType with the modelName altered as directed
func (model ModelType) WithModelName(newModelName string) ModelType {
	var newModel ModelType = model
	newModel.modelName = newModelName
	return newModel
}

// GetModelName returns the current friendly name for this model type
func (model ModelType) GetModelName() string {
	return model.modelName
}

// WithCollectionName returns a ModelType with the collectionName altered as directed
func (model ModelType) WithCollectionName(newCollectionName string) ModelType {
	var newModel ModelType = model
	newModel.collectionName = newCollectionName
	return newModel
}

// GetCollectionName returns the current default collection name for this model type
func (model ModelType) GetCollectionName() string {
	return model.collectionName
}

// WithDatabaseName returns a ModelType with the databaseName altered as directed
func (model ModelType) WithDatabaseName(newDatabaseName string) ModelType {
	var newModel ModelType = model
	newModel.databaseName = newDatabaseName
	return newModel
}

// GetDatabaseName returns the database name for this ModelType
func (model ModelType) GetDatabaseName() string {
	if model.databaseName == "" {
		return model.GetClient().Database
	}
	return model.databaseName
}

// WithClientName returns a ModelType with the clientName altered as directed
func (model ModelType) WithClientName(newClientName string) ModelType {
	newModel := ModelType{}
	newModel = model
	newModel.clientName = newClientName
	return newModel
}

// GetClientName returns the custom client name for this ModelType, or "" if using the default
func (model *ModelType) GetClientName() string {
	return model.clientName
}

// GetClient returns the Client use by this ModelType
func (model *ModelType) GetClient() *Client {
	if clientName := model.GetClientName(); clientName != "" {
		log.Fatal("no")
		return ClientByName(clientName)
	}
	return DefaultClient()
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// New intantiates a new document model object of the registered type and returns a pointer to the new object.
// The returned object will be preset with the defaults specified during initial document/model registration.
// Note: Due to the strongly typed nature of Go, you'll need to perform a type assertion (as the value is returned as an interface{})
func (model ModelType) New() IDocumentBase {
	log.Debugf("%v.New()", model.GetModelName())
	retAsIDocumentBase := makeDocument(&model, model.GetDefaultBSON())
	return retAsIDocumentBase // return the new object as an IDocumentBase interface
}

// GetDefaultBSON provides the default values for a ModelType returned as a BsonDocument.
// The returned value is deep-cloned to protect the original data, so you can begin using it directly without a second deep copy
func (model ModelType) GetDefaultBSON() BsonDocument {
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

// returns a handle to the mongo driver collection for this ModelType
func (model *ModelType) getMongoCollectionHandle() *mongo.Collection {
	client := model.GetClient()
	dbName := model.GetDatabaseName()
	collectionName := model.GetCollectionName()
	collectionRef := client.getMongoCollectionHandle(dbName, collectionName)
	return collectionRef
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
