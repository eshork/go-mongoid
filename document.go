package mongoid

import (
	"mongoid/log"
)

// IDocumentBase ...
type IDocumentBase interface {
	initDocumentBase(selfRef IDocumentBase, defaultsBSON BsonDocument)
	DocumentBase() IDocumentBase
	Model() *ModelType

	ToBson() BsonDocument

	GetID() interface{}

	IsPersisted() bool
	IsChanged() bool
	Changes() BsonDocument

	Save() error

	// SetCollection(*mgo.Collection)
	// SetDocument(document IDocumentBase)
	// SetConnection(*Connection)

	// Create()
	// Create_()
	// Validate() error

	// Changes()
	// Changed_(fieldName)
	// Reset_(fieldName)
	// Was_(fieldName)

	SetField(fieldNamePath string, newValue interface{}) error
	GetField(fieldNamePath string) (interface{}, error)
}

// Base ...
type Base struct {
	rootTypeRef IDocumentBase // self-reference for future type recognition via interface{}
	// mongoidCollectionName string // change to collection pointer
	// mongoidDatabaseName   string // change to database pointer
	// mongoidClientName     string // change to client pointer

	persisted bool // persistence tracking (reflects the anticipated existence of a record within the datastore, based on the lifecycle of the instance)
	// privateID     string       // internal object ID tracker (string form in case a custom ID field is provided of a non-ObjectID type)
	previousValue BsonDocument // stores a BSON representation of the last values, used for change tracking
}

// initDocumentBase configures this IDocumentBase with a self-reference and a default state
// Self-reference is used to :
//   - store the original object-type
//   - store a copy of the initial object values for future change tracking
func (d *Base) initDocumentBase(selfRef IDocumentBase, defaultsBSON BsonDocument) {
	d.rootTypeRef = selfRef
	if d.rootTypeRef == nil {
		panic("cannot InitDocumentBase without a valid selfRef handle")
	}
	// log.Warn("YOLO")
	if defaultsBSON != nil {
		// log.Warn("YOLO2")
		// log.Warn(defaultsBSON)
		d.setPreviousValueBSON(defaultsBSON)
		structValuesFromBsonM(selfRef, defaultsBSON) // TODO: consider replacement with Reset() once that's written (assuming that won't have callback/event consequences)
	}
}

// force sets the last object values via BsonDocument
func (d *Base) setPreviousValueBSON(lastValue BsonDocument) {
	d.previousValue = lastValue
}

// updates the stored previousValue BSON with the current object values (resets value change tracking)
func (d *Base) refreshPreviousValueBSON() {
	d.setPreviousValueBSON(d.ToBson())
}

// Model returns the mongoid.ModelType of the document object, or nil if unknown
func (d *Base) Model() *ModelType {
	log.Trace("Base.Model()")
	if d.rootTypeRef == nil {
		log.Panic("tried to get Model() but selfRef not set - did you forget to InitDocumentBase()?")
	}
	return Model(d.rootTypeRef)
}

// DocumentBase returns the self-reference handle, which can be used to un-cast the object from *Base into an IDocumentBase (interface{}) of the original type
func (d *Base) DocumentBase() IDocumentBase {
	if d.rootTypeRef == nil {
		log.Panic("tried to get DocumentBase() but selfRef not set - did you forget to InitDocumentBase()?")
	}
	return d.rootTypeRef
}

// GetID returns an interface to the current document ID. Type assertion is left to the caller.
// TODO: make this work whether a custom ID field was explicitly declared in the document model (bson:"_id") or not
func (d *Base) GetID() interface{} {
	res, err := d.GetField("_id")
	if err != nil {
		return nil
	}
	return res
}
