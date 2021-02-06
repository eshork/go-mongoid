package mongoid

import (
	"mongoid/log"
)

// IDocument is the interface that structs must implement to be used with go-mongoid.
// The Document struct implements this interface, and it can be easily applied to your
// own structs by anonymous field.
//
// For example:
//     type ExampleMinimalDocumentStruct struct {
//         mongoid.Document // implements IDocument
//     }
// Refer to examples/ for additional usage examples.
type IDocument interface {
	DocumentBase() IDocument
	ModelType() ModelType

	ToBson() BsonDocument
	ToUpdateBson() BsonDocument

	GetID() interface{}

	IsPersisted() bool
	IsChanged() bool
	Changes() BsonDocument
	Was(fieldPath string) (interface{}, bool)
	// Reset(fieldPath string)
	// ResetAll()

	Save() error

	// Changes()
	// Changed_(fieldName)
	// Reset_(fieldName)
	// Was_(fieldName)

	SetField(fieldNamePath string, newValue interface{}) error
	GetField(fieldNamePath string) (interface{}, error)

	implementIDocumentBase()
	initDocumentBase(modelType *ModelType, selfRef IDocument, initialBSON BsonDocument)
}

// Document ...
type Document struct {
	rootTypeRef   IDocument    // self-reference for future type recognition via interface{}
	persisted     bool         // persistence tracking (reflects the anticipated existence of a record within the datastore, based on the lifecycle of the instance)
	previousValue BsonDocument // stores a BSON representation of the last values, used for change tracking
	modelType     *ModelType   // the ModelType that was used to create this object
	// privateID     string       // internal object ID tracker (string form in case a custom ID field is provided of a non-ObjectID type)
}

// implementIDocumentBase implements IDocument
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

// DocumentBase returns the self-reference handle, which can be used to un-cast the object from *Document into an IDocument (interface{}) of the original type
func (d *Document) DocumentBase() IDocument {
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
