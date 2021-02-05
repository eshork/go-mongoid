package mongoid

import (
	"mongoid/log"
)

// IDocumentBase ...
type IDocumentBase interface {
	DocumentBase() IDocumentBase
	ModelType() ModelType

	ToBson() BsonDocument
	ToUpdateBson() BsonDocument

	GetID() interface{}

	IsPersisted() bool
	setPersisted(bool)
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

	implementIDocumentBase()
	initDocumentBase(modelType *ModelType, selfRef IDocumentBase, initialBSON BsonDocument)
}

// Document ...
type Document struct {
	rootTypeRef   IDocumentBase // self-reference for future type recognition via interface{}
	persisted     bool          // persistence tracking (reflects the anticipated existence of a record within the datastore, based on the lifecycle of the instance)
	previousValue BsonDocument  // stores a BSON representation of the last values, used for change tracking
	modelType     *ModelType    // the ModelType that was used to create this object
	// privateID     string       // internal object ID tracker (string form in case a custom ID field is provided of a non-ObjectID type)
}

// implementIDocumentBase implements IDocumentBase
func (d *Document) implementIDocumentBase() {}

// force sets previousValue (change tracking) to the given BsonDocument
func (d *Document) setPreviousValueBSON(lastValue BsonDocument) {
	d.previousValue = lastValue
}

// updates the stored previousValue BSON (change tracking) with the current object values (resets value change tracking)
func (d *Document) refreshPreviousValueBSON() {
	d.setPreviousValueBSON(d.ToBson())
}

// ModelType returns the ModelType of the document object
func (d *Document) ModelType() ModelType {
	log.Trace("Document.Model()")
	if d.modelType == nil {
		log.Trace("Document.Model() d.modelType is nil; creating ModelType on demand")
		mt := Model(d)
		d.modelType = &mt
	}
	return *d.modelType
}

// DocumentBase returns the self-reference handle, which can be used to un-cast the object from *Document into an IDocumentBase (interface{}) of the original type
func (d *Document) DocumentBase() IDocumentBase {
	if d.rootTypeRef == nil {
		log.Panic("DocumentBase() requires valid rootTypeRef")
	}
	return d.rootTypeRef
}

// GetID returns an interface to the current document ID. Type assertion is left to the caller.
// TODO: make this work whether a custom ID field was explicitly declared in the document model (bson:"_id") or not
func (d *Document) GetID() interface{} {
	res, err := d.GetField("_id")
	if err != nil {
		return nil
	}
	return res
}
