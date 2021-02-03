package mongoid

import (
	"mongoid/log"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
)

// makeDocument creates a new object of type docType, populated with the given srcDoc
func makeDocument(docType *ModelType, srcDoc bson.M) IDocumentBase {
	log.Trace("makeDocument()")
	typeRef := reflect.Indirect(reflect.ValueOf(docType.rootTypeRef))        // model.rootTypeRef is always a ptr to an example object, so we need to use Indirect()
	ret := reflect.New(typeRef.Type())                                       // finally have a solid object type, so make one
	retAsIDocumentBase := ret.Interface().(IDocumentBase)                    // convert into a IDocumentBase interface
	retAsIDocumentBase.initDocumentBase(docType, retAsIDocumentBase, srcDoc) // call the self init
	return retAsIDocumentBase
}

// initDocumentBase configures this IDocumentBase with a self-reference and an initial state
// Self-reference is used to :
//   - store the original object-type
//   - store a copy of the initial object values for future change tracking
func (d *Base) initDocumentBase(modelType *ModelType, selfRef IDocumentBase, initialBSON BsonDocument) {
	d.rootTypeRef = selfRef
	d.modelType = modelType
	if d.rootTypeRef == nil {
		panic("cannot initDocumentBase without a valid selfRef handle")
	}
	if d.modelType == nil {
		panic("cannot initDocumentBase without a valid modelType handle")
	}
	if initialBSON != nil {
		structValuesFromBsonM(selfRef, initialBSON)
		// benefits to using d.setPreviousValueBSON instead of d.refreshPreviousValueBSON here:
		//  - skip a call to ToBson(), since we already have a BSON formatted representation of the desired state (ie, faster)
		//  - tests have more opportunity to uncover issues with to/from bson converters and the value initialization code
		d.setPreviousValueBSON(initialBSON)
	}
}
