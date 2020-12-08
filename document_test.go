package mongoid_test

import (
	"mongoid"
	"mongoid/util"

	"fmt"

	"github.com/brianvoe/gofakeit"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// left unregistered for test purposes
type UnknownExampleDocument struct {
	mongoid.Base
	StringField string
}

// left unregistered for test purposes
// var UnknownExampleDocuments = mongoid.Register(&UnknownExampleDocument{})

type ExampleSimpeInlinableDocument1 struct {
	Inlined1StringField    string
	Inlined1IntField       int
	Inlined1BoolField      bool
	Inlined1StringPtrField *string
	Inlined1IntPtrField    *int
	Inlined1BoolPtrField   *bool
}

type ExampleSimpeInlinableDocument2 struct {
	Inlined2StringField    string
	Inlined2IntField       int
	Inlined2BoolField      bool
	Inlined2StringPtrField *string
	Inlined2IntPtrField    *int
	Inlined2BoolPtrField   *bool
}

type ExampleSimpeInlinableDocument3 struct {
	Inlined3StringField    string
	Inlined3IntField       int
	Inlined3BoolField      bool
	Inlined3StringPtrField *string
	Inlined3IntPtrField    *int
	Inlined3BoolPtrField   *bool
}

type ExampleSimpleEmbeddableDocument struct {
	StringField    string
	IntField       int
	BoolField      bool
	StringPtrField *string
	IntPtrField    *int
	BoolPtrField   *bool
}

// type ExampleEmbedableStruct2 struct {
// 	DifferentStringField string `bson:"renamed_string_field"`
// }

// core example doc type - used as the basis of all document related tests
type ExampleDocument struct {
	mongoid.Base `mongoid:"collection:tacocat"`
	// mongoid.Timestamps `bson:",inline"`
	ID mongoid.ObjectID `bson:"_id"`

	StringField    string
	IntField       int
	BoolField      bool
	StringPtrField *string
	IntPtrField    *int
	BoolPtrField   *bool

	IntSliceField        []int
	IntPtrSliceField     []*int
	IntPtrSliceFieldNils []*int

	// all 3 embed variations
	SimpleEmbed       ExampleSimpleEmbeddableDocument
	SimpleEmbedPtr    *ExampleSimpleEmbeddableDocument
	SimpleEmbedPtrSet *ExampleSimpleEmbeddableDocument

	// all 3 inline variations
	SimpleInline       ExampleSimpeInlinableDocument1  `bson:",inline"`
	SimpleInlinePtrSet *ExampleSimpeInlinableDocument2 `bson:",inline"`
	SimpleInlinePtr    *ExampleSimpeInlinableDocument3 `bson:",inline"`

	SimpleEmbedSliceField    []ExampleSimpleEmbeddableDocument
	SimpleEmbedPtrSliceField []*ExampleSimpleEmbeddableDocument
	// SimpleEmbedSlicePtr      *[]ExampleSimpleEmbeddableDocument

	// StringPtrField1 *string // this little pointer gets one
	// StringPtrField2 *string // this little pointer gets none
	// IntPtrField1    *int    // this little pointer gets one
	// IntPtrField2    *int    // this little pointer gets none
	// OmittedBoolField   bool `bson:"-"`
	privateStringField string
	// IntArrayField      []int // TODO: array known does not work yet

	// DefaultEmbeddedStructPtr *ExampleEmbedableStruct1
	// EmbeddedStructPtr1       *ExampleEmbedableStruct1
	// EmbeddedStructPtr1Filled *ExampleEmbedableStruct1

	// EmbeddedStruct2 ExampleEmbedableStruct2 `bson:"renamed_embedded_struct"`

	// RenamedEmbeddedStruct1 ExampleEmbedableStruct1 `bson:"some_embedded_struct"`

	// RelationshipsAreHard interface{} // need to make sure relationships stay hard

	// ExampleEmbedableStruct1
	// InlineStructField    ExampleEmbedableStruct1  `bson:", inline"`
	// InlineStructPtrField *ExampleEmbedableStruct1 //`bson:", inline"`
	// StructPtrField *ExampleEmbedableStruct1
	// InlineStructPtrField *ExampleEmbedableStruct1 `bson:", inline"`
}

var TmpStringFieldValue1 = "!racecar!"
var TmpIntFieldValue1 = 41
var TmpIntFieldValue2 = 42
var TmpSimpleEmbedValue = ExampleSimpleEmbeddableDocument{
	StringField: "TmpSimpleEmbedValue = ExampleSimpleEmbeddableDocument{StringField}",
}
var TmpSimpleInlineValue = ExampleSimpeInlinableDocument2{Inlined2StringField: "I am inlined"}
var TmpSimpleEmbedSliceValue = []ExampleSimpleEmbeddableDocument{TmpSimpleEmbedValue, TmpSimpleEmbedValue}

// register the model with some default values
var ExampleDocuments = mongoid.Register(&ExampleDocument{
	StringField:          "tacocat is tacocat backwards",
	IntField:             42,
	IntSliceField:        []int{1, 2, 4, 8, 16},
	IntPtrSliceField:     []*int{&TmpIntFieldValue1, &TmpIntFieldValue2},
	IntPtrSliceFieldNils: []*int{nil, nil, nil},
	SimpleEmbedPtrSet:    &TmpSimpleEmbedValue,
	SimpleInlinePtrSet:   &TmpSimpleInlineValue,

	SimpleEmbedSliceField:    []ExampleSimpleEmbeddableDocument{TmpSimpleEmbedValue, TmpSimpleEmbedValue},
	SimpleEmbedPtrSliceField: []*ExampleSimpleEmbeddableDocument{&TmpSimpleEmbedValue, &TmpSimpleEmbedValue},
	// SimpleEmbedSlicePtr:      &TmpSimpleEmbedSliceValue,

})

var _ = Describe("Document", func() {

	Context("an unknown UnknownExampleDocument document model", func() {
		It("is not verifyably registered", func() {
			By("struct name")
			Expect(mongoid.M("UnknownExampleDocument")).To(BeNil())
			By("example ref object")
			Expect(mongoid.M(&UnknownExampleDocument{})).To(BeNil())
		})
	})

	Context("ExampleDocument document model", func() {

		It("can be New()'ed", func() {
			newObj := ExampleDocuments.New()
			By("being an actual object")
			Expect(newObj).ToNot(BeNil())
			By("being the correct object type")
			Expect(newObj).To(BeAssignableToTypeOf(&ExampleDocument{}))
		})

		It("begins unpersisted", func() {
			newObj := ExampleDocuments.New()
			Expect(newObj.IsPersisted()).To(BeFalse(), "expects to not yet be persisted")
		})

		It("begins unchanged", func() {
			newObj := ExampleDocuments.New()
			Expect(newObj.IsChanged()).To(BeFalse(), "expects to not yet be changed")
		})

		It("identifies a simple change via IsChanged()", func() {
			newObj := ExampleDocuments.New().(*ExampleDocument)
			newObj.IntSliceField = []int{}
			newObj.StringField = gofakeit.HipsterWord()
			Expect(newObj.IsChanged()).To(BeTrue(), "expect a change")
		})

		It("identifies a simple slice field change via IsChanged()", func() {
			newObj := ExampleDocuments.New().(*ExampleDocument)
			newObj.IntSliceField = []int{gofakeit.Number(1, 99)}
			Expect(newObj.IsChanged()).To(BeTrue(), "expect a change")
		})

		It("identifies a slice field clearing via IsChanged()", func() {
			newObj := ExampleDocuments.New().(*ExampleDocument)
			newObj.IntSliceField = []int{}
			Expect(newObj.IsChanged()).To(BeTrue(), "expect a change")
		})

		PIt("recalls a previous field value via Was(fieldName)", func() {
			// marked pending because .Was is currently NYI - ref gihub issue #5
			newObj := ExampleDocuments.New().(*ExampleDocument)
			oldValue := newObj.StringField
			newValue := gofakeit.HipsterWord()
			newObj.StringField = newValue
			Expect(newObj.IsChanged()).To(BeTrue(), "expect a change")
			wasValue, _ := newObj.Was("string_field")
			Expect(wasValue).To(Equal(oldValue), "expect old value to be preserved")
		})

		It("generates concrete maps and arrays via ToBson()", func() {
			newObj := ExampleDocuments.New().(*ExampleDocument)
			bsonM := newObj.ToBson()
			invalidList := util.ValidateBson(bsonM)
			Expect(invalidList).To(BeEmpty(), "found non-marshallable complex types within BSON key path(s) from: "+fmt.Sprintf("%+v", bsonM))
		})

		PIt("can be Save()'ed and Find()'ed", func() {
			OnlineDatabaseOnly(func() {

				By("object creation")
				newObj := mongoid.M("ExampleDocument").New().(*ExampleDocument)
				Expect(newObj).ToNot(BeNil(), "expects a real object to be created")

				By("object Persisted()==false check")
				Expect(newObj.IsPersisted()).To(BeFalse(), "expects to not yet be persisted")

				By("object IsChanged()==false check")
				Expect(newObj.IsChanged()).To(BeFalse(), "expects to be unchanged")

				initialObjectID := newObj.ID

				By("Save()'ing")
				Expect(newObj.Save()).To(BeNil(), "expects no errors")
				Expect(newObj.IsPersisted()).To(BeTrue(), "expects to now be persisted")

				actualObjectID := newObj.ID
				Expect(initialObjectID).ToNot(Equal(actualObjectID), "expects objectID to be updated")

				By("another object IsChanged()==false check")
				Expect(newObj.IsChanged()).To(BeFalse(), fmt.Sprintf("expects to be unchanged but found: %+v", newObj.Changes()))

				By("GetID()")
				objectIDInt := newObj.GetID()
				Expect(objectIDInt).ToNot(BeNil(), "expects an ID value")
				objectID, ok := objectIDInt.(mongoid.ObjectID)
				Expect(ok).To(BeTrue(), "expects ID to type-assert into ObjectID")
				Expect(newObj.GetID().(mongoid.ObjectID)).To(Equal(newObj.ID), "expects newObj.GetID().(ObjectID) == newObj.ID")

				By("Find()'ing")
				res, err := mongoid.M("ExampleDocument").Find(objectID)
				Expect(err).ToNot(HaveOccurred())
				foundObj := res.OneAndClose().(*ExampleDocument)
				Expect(foundObj.ID).To(Equal(newObj.ID), "expects foundObj.ID == newObj.ID")
			})
		})

	})

})
