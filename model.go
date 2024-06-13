package mongoid

import (
	"fmt"
	"mongoid/log"
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
)

type collectionHandle struct {
	rootTypeRef    IDocument // a reference to an object instance of the document/model type given during registration (for future type sanity)
	modelName      string
	modelFullName  string
	collectionName string
	databaseName   string
	clientName     string
	client         *Client      // populated on creation based on clientName
	defaultValue   BsonDocument // bson representation of default values to be applied during creation of brand new document/model instances
}

var _ fmt.Stringer = collectionHandle{} // assert implements Stringer interface
var _ ICollection = collectionHandle{}  // assert implements ICollection interface

// String implements fmt.Stringer interface
func (col collectionHandle) String() string {
	// collectionHandle(col.modelFullName)[client:cName,database:dbName,collection:colName]
	extras := ""
	appendExtras := func(name, value string) {
		if extras != "" {
			extras = fmt.Sprintf("%s,%s:%s", extras, name, value)
		}
		extras = fmt.Sprintf("%s:%s", name, value)
	}
	appendExtras("client", col.clientName)
	appendExtras("database", col.databaseName)
	appendExtras("collection", col.collectionName)
	return fmt.Sprintf("Collection(%s/%s)[%s]", col.modelName, col.modelFullName, extras)
}

// Name implements ICollection.TypeName()
func (col collectionHandle) TypeName() string {
	return col.modelName
}

// reflect-ive rootTypeRef type equality check
func (col *collectionHandle) equalsRootType(comparisonModel *collectionHandle) bool {
	if col != nil && comparisonModel != nil && col.rootTypeRef != nil && comparisonModel.rootTypeRef != nil {
		if reflect.TypeOf(col.rootTypeRef) == reflect.TypeOf(comparisonModel.rootTypeRef) {
			if reflect.TypeOf(col.rootTypeRef).Kind() == reflect.Ptr {
				if reflect.TypeOf(col.rootTypeRef).Elem() == reflect.TypeOf(comparisonModel.rootTypeRef).Elem() {
					return true
				}
				return false
			}
			return true
		}
	}
	return false
}

// GetModelName returns the current friendly name for this model type
func (col collectionHandle) GetModelName() string {
	return col.modelName
}

// WithCollectionName returns a collectionHandle with the collectionName altered as directed
func (col collectionHandle) WithCollectionName(newCollectionName string) ICollection {
	var newCol collectionHandle = col
	newCol.collectionName = newCollectionName
	return newCol
}

// GetCollectionName returns the current default collection name for this model type
func (col collectionHandle) CollectionName() string {
	return col.collectionName
}

// WithDatabaseName returns a collectionHandle with the databaseName altered as directed
func (col collectionHandle) WithDatabaseName(newDatabaseName string) ICollection {
	var newModel collectionHandle = col
	newModel.databaseName = newDatabaseName
	return newModel
}

// GetDatabaseName returns the database name for this collectionHandle
func (col collectionHandle) DatabaseName() string {
	if col.databaseName == "" {
		return col.Client().Database
	}
	return col.databaseName
}

// WithClientName returns a collectionHandle with the clientName altered as directed
func (col collectionHandle) WithClientName(newClientName string) ICollection {
	newModel := collectionHandle{}
	newModel = col
	newModel.clientName = newClientName
	return newModel
}

// GetClientName returns the custom client name for this collectionHandle, or "" if using the default
func (col collectionHandle) ClientName() string {
	return col.clientName
}

// GetClient returns the Client used by this collectionHandle
func (col collectionHandle) Client() *Client {
	if col.client != nil {
		return col.client
	}
	if col.clientName != "" {
		col.client = ClientByName(col.clientName)
	}
	if col.client == nil {
		col.client = DefaultClient()
	}
	return col.client
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// New intantiates a new document model object of the registered type and returns a pointer to the new object.
// The returned object will be preset with the defaults specified during initial document/model registration.
// Note: Due to the strongly typed nature of Go, you'll need to perform a type assertion (as the value is returned as an interface{})
func (col collectionHandle) New() IDocument {
	log.Debugf("%v.New()", col.GetModelName())
	retAsIDocumentBase := makeDocument(&col, col.DefaultBSON())
	return retAsIDocumentBase // return the new object as an IDocument interface
}

// DefaultBSON provides the default values for a collectionHandle returned as a BsonDocument.
// The returned value is deep-cloned to protect the original data, so you can begin using it directly without a second deep copy
func (col collectionHandle) DefaultBSON() BsonDocument {
	log.Trace("DefaultBSON()")
	if col.defaultValue == nil {
		return BsonDocument{} // empty document if no defaults available
	}
	return BsonDocumentDeepCopy(col.defaultValue)
}

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

// returns a handle to the mongo driver collection for this collectionHandle
func (col collectionHandle) getMongoCollectionHandle() *mongo.Collection {
	client := col.Client()
	dbName := col.DatabaseName()
	collectionName := col.CollectionName()
	collectionRef := client.getMongoCollectionHandle(dbName, collectionName)
	return collectionRef
}
