package mongoid

import (
	"mongoid/log"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
)

// makeDocument creates a new object of type docType, populated with the given srcDoc
func makeDocument(docType *ModelType, srcDoc bson.M) IDocumentBase {
	log.Debug("ModelType.New()")
	typeRef := reflect.Indirect(reflect.ValueOf(docType.rootTypeRef)) // model.rootTypeRef is always a ptr to an example object, so we need to use Indirect()
	ret := reflect.New(typeRef.Type())                                // finally have a solid object type, so make one
	retAsIDocumentBase := ret.Interface().(IDocumentBase)             // convert into a IDocumentBase interface
	retAsIDocumentBase.initDocumentBase(retAsIDocumentBase, srcDoc)   // call the self init
	return retAsIDocumentBase
}

// initDocumentBase configures this IDocumentBase with a self-reference and an initial state
// Self-reference is used to :
//   - store the original object-type
//   - store a copy of the initial object values for future change tracking
func (d *Base) initDocumentBase(selfRef IDocumentBase, initialBSON BsonDocument) {
	d.rootTypeRef = selfRef
	if d.rootTypeRef == nil {
		panic("cannot initDocumentBase without a valid selfRef handle")
	}
	if initialBSON != nil {
		d.setPreviousValueBSON(initialBSON)
		structValuesFromBsonM(selfRef, initialBSON)
	}
}
