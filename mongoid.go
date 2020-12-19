package mongoid

import (
	"mongoid/log"
)

// Model returns a ModelType interface, used as the root interface the query and creation methods that may be performed related to a registered Document/Model
// The single input argument is either the string name that the model was registered under, or a pointer to an object of the desired type.
// If the input does not resolve to a registered type, an error will be logged and Nil will be returned.
// (Unfound Document is not Fatal by design, the log should indicate the issue, and your code will most certainly panic soon enough if it doesn't explicitly validate return values)
func Model(modelRef interface{}) *ModelType {
	switch modelRef.(type) {
	case string:
		return getRegisteredModelTypeByName(modelRef.(string))
	case IDocumentBase:
		return getRegisteredModelTypeByDocRef(modelRef.(IDocumentBase))
	default:
		log.Error("given modelRef could not type assert to string or IDocumentBase")
		return nil
	}
}

// M is a shorthand function for Model()
func M(modelRef interface{}) *ModelType {
	return Model(modelRef)
}

// mongoid.M(&MyThing{}).New()
// mongoid.M(&MyThing{}).Create(options)
// mongoid.M(&MyThing{}).Find(options)

// mongoid.New("ModelName", contentBSON{}).(ModelObject)
// mongoid.New(ModelObject{}).(ModelObject)

// mongoid.Create("ModelName", contentBSON{}).(ModelObject)
// mongoid.Create(ModelObject{}).(ModelObject)
// mongoid.CreateV("ModelName", contentBSON{}).(ModelObject)
// mongoid.CreateV(ModelObject{}).(ModelObject)

// model-identifier + source-selector + query-builder + finality

//Source-Selector
// WithDatabase()
// WithCollection()
// WithClient()

//Finality
// First
// Last
// Pluck
// All
// Each? Iter? Next?

// mongoid.M("Thing").Where(bson.Reasons{}).First
