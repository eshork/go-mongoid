package mongoid_test

import (
	"fmt"
	"math"
	"math/cmplx"
	"mongoid"
	"mongoid/util"

	"github.com/brianvoe/gofakeit"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/bson"
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

	IntSliceField []int
	// IntPtrSliceField     []*int
	// IntPtrSliceFieldNils []*int

	// // all 3 embed variations
	// SimpleEmbed       ExampleSimpleEmbeddableDocument
	// SimpleEmbedPtr    *ExampleSimpleEmbeddableDocument
	// SimpleEmbedPtrSet *ExampleSimpleEmbeddableDocument

	// // all 3 inline variations
	// SimpleInline       ExampleSimpeInlinableDocument1  `bson:",inline"`
	// SimpleInlinePtrSet *ExampleSimpeInlinableDocument2 `bson:",inline"`
	// SimpleInlinePtr    *ExampleSimpeInlinableDocument3 `bson:",inline"`

	// SimpleEmbedSliceField    []ExampleSimpleEmbeddableDocument
	// SimpleEmbedPtrSliceField []*ExampleSimpleEmbeddableDocument
	// SimpleEmbedSlicePtr      *[]ExampleSimpleEmbeddableDocument

	// StringPtrField1 *string // this little pointer gets one
	// StringPtrField2 *string // this little pointer gets none
	// IntPtrField1    *int    // this little pointer gets one
	// IntPtrField2    *int    // this little pointer gets none
	// OmittedBoolField   bool `bson:"-"`
	// privateStringField string
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
	StringField:    "tacocat is tacocat backwards",
	IntField:       42,
	BoolField:      true,
	StringPtrField: func() *string { retval := "test"; return &retval }(),
	IntSliceField:  []int{1, 2, 4, 8, 16},
	// IntPtrSliceField:     []*int{&TmpIntFieldValue1, &TmpIntFieldValue2},
	// IntPtrSliceFieldNils: []*int{nil, nil, nil},
	// SimpleEmbedPtrSet:    &TmpSimpleEmbedValue,
	// SimpleInlinePtrSet:   &TmpSimpleInlineValue,

	// SimpleEmbedSliceField:    []ExampleSimpleEmbeddableDocument{TmpSimpleEmbedValue, TmpSimpleEmbedValue},
	// SimpleEmbedPtrSliceField: []*ExampleSimpleEmbeddableDocument{&TmpSimpleEmbedValue, &TmpSimpleEmbedValue},
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
			// newObj.IntSliceField = []int{}
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

		It("identifies a string ptr (*string) field change via IsChanged()", func() {
			By("starting with a new ExampleDocument")
			newObj := ExampleDocuments.New().(*ExampleDocument)
			Expect(newObj.IsChanged()).To(BeFalse(), "expect no changes after creation")
			OnlineDatabaseOnly(func() {
				By("saving the ExampleDocument")
				newObj.Save()
				Expect(newObj.IsChanged()).To(BeFalse(), "expect no changes after save")
				By("unsetting StringPtrField (setting it to nil)")
				newObj.StringPtrField = nil
				Expect(newObj.IsChanged()).To(BeTrue(), "expect a change")
				Expect(newObj.Changes()).ToNot(Equal(bson.M(nil)), "expect a change")
				By("saving the ExampleDocument")
				newObj.Save()
				Expect(newObj.IsChanged()).To(BeFalse(), "expect no changes after save")
				// Expect(newObj.Changes()).To(Equal(bson.M(nil)), "expect no changes after save")
				By("changing StringPtrField to something else")
				newValue := "something else"
				newObj.StringPtrField = &newValue
				Expect(newObj.Changes()).ToNot(Equal(bson.M(nil)), "expect a change")
				Expect(newObj.IsChanged()).ToNot(BeFalse(), "expect a change")
				By("saving the ExampleDocument")
				newObj.Save()
				Expect(newObj.Changes()).To(Equal(bson.M(nil)), "expect no changes after save")
				Expect(newObj.IsChanged()).To(BeFalse(), "expect no changes after save")
			})
		})

		It("recalls a previous bool field value via Was(fieldName) [t6c81550b]", func() {
			newObj := ExampleDocuments.New().(*ExampleDocument)
			Expect(newObj.IsChanged()).To(BeFalse(), "expect no change")
			oldValue := newObj.BoolField
			newValue := false
			newObj.BoolField = newValue
			Expect(newObj.IsChanged()).To(BeTrue(), "expect a change")
			wasValue, changed := newObj.Was("bool_field")
			Expect(changed).To(BeTrue(), "expect field to indicate it was changed")
			Expect(wasValue).To(Equal(oldValue), "expect old value to be preserved")
		})
		It("recalls a previous string field value via Was(fieldName) [t8508c303]", func() {
			newObj := ExampleDocuments.New().(*ExampleDocument)
			Expect(newObj.IsChanged()).To(BeFalse(), "expect no change")
			oldValue := newObj.StringField
			newValue := gofakeit.HipsterWord()
			newObj.StringField = newValue
			Expect(newObj.IsChanged()).To(BeTrue(), "expect a change")
			wasValue, changed := newObj.Was("string_field")
			Expect(changed).To(BeTrue(), "expect field to indicate it was changed")
			Expect(wasValue).To(Equal(oldValue), "expect old value to be preserved")
		})

		It("generates concrete maps and arrays via ToBson()", func() {
			newObj := ExampleDocuments.New().(*ExampleDocument)
			bsonM := newObj.ToBson()
			invalidList := util.ValidateBson(bsonM)
			Expect(invalidList).To(BeEmpty(), "found non-marshallable complex types within BSON key path(s) from: "+fmt.Sprintf("%+v", bsonM))
		})

		It("can be Save()'ed and Find()'ed [t62adaaca]", func() {
			OnlineDatabaseOnly(func() {

				By("object creation")
				newObj := mongoid.M("ExampleDocument").New().(*ExampleDocument)
				Expect(newObj).ToNot(BeNil(), "expects a real object to be created")

				By("object Persisted()==false check")
				Expect(newObj.IsPersisted()).To(BeFalse(), "expects to not yet be persisted")

				By("object IsChanged()==false check")
				Expect(newObj.IsChanged()).To(BeFalse(), "expects to be unchanged")

				initialObjectID := newObj.ID
				Expect(initialObjectID).To(Equal(mongoid.ObjectID{}), "expects initialObjectID to be zero-value")

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
				res := mongoid.M("ExampleDocument").Find(objectID)
				foundObj := res.One().(*ExampleDocument)
				Expect(foundObj.ID).To(Equal(newObj.ID), "expects foundObj.ID == newObj.ID")
			})
		})

	})

})

var _ = Describe("Document", func() {

	Context("verifying storage of Document field types", func() {

		It("bool field [td5459dce]", func() {
			type BoolTestStruct struct {
				mongoid.Base `mongoid:"collection:bool_test"`
				ID           mongoid.ObjectID `bson:"_id"`
				Field        bool
			}
			BoolTestStructs := mongoid.Register(&BoolTestStruct{Field: false})
			newObj := BoolTestStructs.New().(*BoolTestStruct)
			newObj.Field = true
			OnlineDatabaseOnly(func() {
				newObj.Save()
				sameObj := BoolTestStructs.Find(newObj.ID).One().(*BoolTestStruct)
				Expect(sameObj.Field).To(Equal(newObj.Field), "retrieved document should have same value as original")
				newObj.Field = false
				newObj.Save()
				Expect(sameObj.Field).ToNot(Equal(newObj.Field), "retrieved document should have different value as original before refetch")
				sameObj = BoolTestStructs.Find(newObj.ID).One().(*BoolTestStruct)
				Expect(sameObj.Field).To(Equal(newObj.Field), "retrieved document should have same value as original after refetch")
			})
		})
		It("string field [tdb0f1026]", func() {
			type StringTestStruct struct {
				mongoid.Base `mongoid:"collection:bool_test"`
				ID           mongoid.ObjectID `bson:"_id"`
				Field        string
			}
			StringTestStructs := mongoid.Register(&StringTestStruct{Field: "original value"})
			newObj := StringTestStructs.New().(*StringTestStruct)
			newObj.Field = "something else"
			OnlineDatabaseOnly(func() {
				newObj.Save()
				sameObj := StringTestStructs.Find(newObj.ID).One().(*StringTestStruct)
				Expect(sameObj.Field).To(Equal(newObj.Field), "retrieved document should have same value as original")
				newObj.Field = "crystal pepsi"
				newObj.Save()
				Expect(sameObj.Field).ToNot(Equal(newObj.Field), "retrieved document should have different value as original before refetch")
				sameObj = StringTestStructs.Find(newObj.ID).One().(*StringTestStruct)
				Expect(sameObj.Field).To(Equal(newObj.Field), "retrieved document should have same value as original after refetch")
			})
		})
		It("int field [t1c6a3f54]", func() {
			type IntTestStruct struct {
				mongoid.Base `mongoid:"collection:int_test"`
				ID           mongoid.ObjectID `bson:"_id"`
				Field        int
			}
			IntTestStructs := mongoid.Register(&IntTestStruct{Field: 0})
			newObj := IntTestStructs.New().(*IntTestStruct)
			newObj.Field = 42
			OnlineDatabaseOnly(func() {
				newObj.Save()
				sameObj := IntTestStructs.Find(newObj.ID).One().(*IntTestStruct)
				Expect(sameObj.Field).To(Equal(newObj.Field), "retrieved document should have same value as original")
				newObj.Field = 7
				newObj.Save()
				Expect(sameObj.Field).ToNot(Equal(newObj.Field), "retrieved document should have different value as original before refetch")
				sameObj = IntTestStructs.Find(newObj.ID).One().(*IntTestStruct)
				Expect(sameObj.Field).To(Equal(newObj.Field), "retrieved document should have same value as original after refetch")
			})
		})
		It("int8 field [tb575e350]", func() {
			type Int8TestStruct struct {
				mongoid.Base `mongoid:"collection:int8_test"`
				ID           mongoid.ObjectID `bson:"_id"`
				Field        int8
			}
			Int8TestStructs := mongoid.Register(&Int8TestStruct{Field: 0})
			newObj := Int8TestStructs.New().(*Int8TestStruct)
			newObj.Field = int8(math.MaxInt8)
			OnlineDatabaseOnly(func() {
				newObj.Save()
				sameObj := Int8TestStructs.Find(newObj.ID).One().(*Int8TestStruct)
				Expect(sameObj.Field).To(Equal(newObj.Field), "retrieved document should have same value as original")
				newObj.Field = 7
				newObj.Save()
				Expect(sameObj.Field).ToNot(Equal(newObj.Field), "retrieved document should have different value as original before refetch")
				sameObj = Int8TestStructs.Find(newObj.ID).One().(*Int8TestStruct)
				Expect(sameObj.Field).To(Equal(newObj.Field), "retrieved document should have same value as original after refetch")
			})
		})
		It("int16 field [ta1075c83]", func() {
			type Int16TestStruct struct {
				mongoid.Base `mongoid:"collection:int16_test"`
				ID           mongoid.ObjectID `bson:"_id"`
				Field        int16
			}
			Int16TestStructs := mongoid.Register(&Int16TestStruct{Field: 0})
			newObj := Int16TestStructs.New().(*Int16TestStruct)
			newObj.Field = int16(math.MaxInt16)
			OnlineDatabaseOnly(func() {
				newObj.Save()
				sameObj := Int16TestStructs.Find(newObj.ID).One().(*Int16TestStruct)
				Expect(sameObj.Field).To(Equal(newObj.Field), "retrieved document should have same value as original")
				newObj.Field = 7
				newObj.Save()
				Expect(sameObj.Field).ToNot(Equal(newObj.Field), "retrieved document should have different value as original before refetch")
				sameObj = Int16TestStructs.Find(newObj.ID).One().(*Int16TestStruct)
				Expect(sameObj.Field).To(Equal(newObj.Field), "retrieved document should have same value as original after refetch")
			})
		})
		It("int32 field [t61c5309f]", func() {
			type Int32TestStruct struct {
				mongoid.Base `mongoid:"collection:int32_test"`
				ID           mongoid.ObjectID `bson:"_id"`
				Field        int32
			}
			Int32TestStructs := mongoid.Register(&Int32TestStruct{Field: 0})
			newObj := Int32TestStructs.New().(*Int32TestStruct)
			newObj.Field = int32(math.MaxInt32)
			OnlineDatabaseOnly(func() {
				newObj.Save()
				sameObj := Int32TestStructs.Find(newObj.ID).One().(*Int32TestStruct)
				Expect(sameObj.Field).To(Equal(newObj.Field), "retrieved document should have same value as original")
				newObj.Field = 7
				newObj.Save()
				Expect(sameObj.Field).ToNot(Equal(newObj.Field), "retrieved document should have different value as original before refetch")
				sameObj = Int32TestStructs.Find(newObj.ID).One().(*Int32TestStruct)
				Expect(sameObj.Field).To(Equal(newObj.Field), "retrieved document should have same value as original after refetch")
			})
		})
		It("int64 field [t486a2496]", func() {
			type Int64TestStruct struct {
				mongoid.Base `mongoid:"collection:int64_test"`
				ID           mongoid.ObjectID `bson:"_id"`
				Field        int64
			}
			Int64TestStructs := mongoid.Register(&Int64TestStruct{Field: 0})
			newObj := Int64TestStructs.New().(*Int64TestStruct)
			newObj.Field = int64(math.MaxInt64)
			OnlineDatabaseOnly(func() {
				newObj.Save()
				sameObj := Int64TestStructs.Find(newObj.ID).One().(*Int64TestStruct)
				Expect(sameObj.Field).To(Equal(newObj.Field), "retrieved document should have same value as original")
				newObj.Field = 7
				newObj.Save()
				Expect(sameObj.Field).ToNot(Equal(newObj.Field), "retrieved document should have different value as original before refetch")
				sameObj = Int64TestStructs.Find(newObj.ID).One().(*Int64TestStruct)
				Expect(sameObj.Field).To(Equal(newObj.Field), "retrieved document should have same value as original after refetch")
			})
		})
		It("no field [tbf3ad97f]", func() {
			type TestStruct struct {
				mongoid.Base `mongoid:"collection:nofield_test"`
				ID           mongoid.ObjectID `bson:"_id"`
			}
			TestStructs := mongoid.Register(&TestStruct{})
			newObj := TestStructs.New().(*TestStruct)
			OnlineDatabaseOnly(func() {
				newObj.Save()
				sameObj := TestStructs.Find(newObj.ID).One().(*TestStruct)
				Expect(newObj.ID).To(Equal(sameObj.ID), "retrieved document should have same ID")
			})
		})
		It("uint field [te355fc16]", func() {
			type UintTestStruct struct {
				mongoid.Base `mongoid:"collection:uint_test"`
				ID           mongoid.ObjectID `bson:"_id"`
				Field        uint
			}
			UintTestStructs := mongoid.Register(&UintTestStruct{Field: 0})
			newObj := UintTestStructs.New().(*UintTestStruct)
			newObj.Field = uint(42)
			OnlineDatabaseOnly(func() {
				newObj.Save()
				sameObj := UintTestStructs.Find(newObj.ID).One().(*UintTestStruct)
				Expect(sameObj.Field).To(Equal(newObj.Field), "retrieved document should have same value as original")
				newObj.Field = 7
				newObj.Save()
				Expect(sameObj.Field).ToNot(Equal(newObj.Field), "retrieved document should have different value as original before refetch")
				sameObj = UintTestStructs.Find(newObj.ID).One().(*UintTestStruct)
				Expect(sameObj.Field).To(Equal(newObj.Field), "retrieved document should have same value as original after refetch")
			})
		})
		It("uint8 field [t1e63da25]", func() {
			type Uint8TestStruct struct {
				mongoid.Base `mongoid:"collection:uint8_test"`
				ID           mongoid.ObjectID `bson:"_id"`
				Field        uint8
			}
			Uint8TestStructs := mongoid.Register(&Uint8TestStruct{Field: 0})
			newObj := Uint8TestStructs.New().(*Uint8TestStruct)
			newObj.Field = uint8(math.MaxUint8)
			OnlineDatabaseOnly(func() {
				newObj.Save()
				sameObj := Uint8TestStructs.Find(newObj.ID).One().(*Uint8TestStruct)
				Expect(sameObj.Field).To(Equal(newObj.Field), "retrieved document should have same value as original")
				newObj.Field = 7
				newObj.Save()
				Expect(sameObj.Field).ToNot(Equal(newObj.Field), "retrieved document should have different value as original before refetch")
				sameObj = Uint8TestStructs.Find(newObj.ID).One().(*Uint8TestStruct)
				Expect(sameObj.Field).To(Equal(newObj.Field), "retrieved document should have same value as original after refetch")
			})
		})
		It("uint16 field [tf2960ac5]", func() {
			type Uint16TestStruct struct {
				mongoid.Base `mongoid:"collection:uint16_test"`
				ID           mongoid.ObjectID `bson:"_id"`
				Field        uint16
			}
			Uint16TestStructs := mongoid.Register(&Uint16TestStruct{Field: 0})
			newObj := Uint16TestStructs.New().(*Uint16TestStruct)
			newObj.Field = uint16(math.MaxUint16)
			OnlineDatabaseOnly(func() {
				newObj.Save()
				sameObj := Uint16TestStructs.Find(newObj.ID).One().(*Uint16TestStruct)
				Expect(sameObj.Field).To(Equal(newObj.Field), "retrieved document should have same value as original")
				newObj.Field = 7
				newObj.Save()
				Expect(sameObj.Field).ToNot(Equal(newObj.Field), "retrieved document should have different value as original before refetch")
				sameObj = Uint16TestStructs.Find(newObj.ID).One().(*Uint16TestStruct)
				Expect(sameObj.Field).To(Equal(newObj.Field), "retrieved document should have same value as original after refetch")
			})
		})
		It("uint32 field [t38e40270]", func() {
			type Uint32TestStruct struct {
				mongoid.Base `mongoid:"collection:uint32_test"`
				ID           mongoid.ObjectID `bson:"_id"`
				Field        uint32
			}
			Uint32TestStructs := mongoid.Register(&Uint32TestStruct{Field: 0})
			newObj := Uint32TestStructs.New().(*Uint32TestStruct)
			newObj.Field = uint32(math.MaxUint32)
			OnlineDatabaseOnly(func() {
				newObj.Save()
				sameObj := Uint32TestStructs.Find(newObj.ID).One().(*Uint32TestStruct)
				Expect(sameObj.Field).To(Equal(newObj.Field), "retrieved document should have same value as original")
				newObj.Field = 7
				newObj.Save()
				Expect(sameObj.Field).ToNot(Equal(newObj.Field), "retrieved document should have different value as original before refetch")
				sameObj = Uint32TestStructs.Find(newObj.ID).One().(*Uint32TestStruct)
				Expect(sameObj.Field).To(Equal(newObj.Field), "retrieved document should have same value as original after refetch")
			})
		})
		It("uint64 field [tc7be9cc8]", func() {
			type Uint64TestStruct struct {
				mongoid.Base `mongoid:"collection:uint64_test"`
				ID           mongoid.ObjectID `bson:"_id"`
				Field        uint64
			}
			Uint64TestStructs := mongoid.Register(&Uint64TestStruct{Field: 0})
			newObj := Uint64TestStructs.New().(*Uint64TestStruct)
			newObj.Field = uint64(math.MaxUint64)
			OnlineDatabaseOnly(func() {
				newObj.Save()
				sameObj := Uint64TestStructs.Find(newObj.ID).One().(*Uint64TestStruct)
				Expect(sameObj.Field).To(Equal(newObj.Field), "retrieved document should have same value as original")
				newObj.Field = 7
				newObj.Save()
				Expect(sameObj.Field).ToNot(Equal(newObj.Field), "retrieved document should have different value as original before refetch")
				sameObj = Uint64TestStructs.Find(newObj.ID).One().(*Uint64TestStruct)
				Expect(sameObj.Field).To(Equal(newObj.Field), "retrieved document should have same value as original after refetch")
			})
		})
		It("complex64 field [t9c9d2b4e]", func() {
			type Complex64TestStruct struct {
				mongoid.Base `mongoid:"collection:int16_test"`
				ID           mongoid.ObjectID `bson:"_id"`
				Field        complex64
			}
			Complex64TestStructs := mongoid.Register(&Complex64TestStruct{Field: 0})
			newObj := Complex64TestStructs.New().(*Complex64TestStruct)
			newObj.Field = complex64(cmplx.Inf())
			OnlineDatabaseOnly(func() {
				newObj.Save()
				sameObj := Complex64TestStructs.Find(newObj.ID).One().(*Complex64TestStruct)
				Expect(sameObj.Field).To(Equal(newObj.Field), "retrieved document should have same value as original")
				newObj.Field = 7
				newObj.Save()
				Expect(sameObj.Field).ToNot(Equal(newObj.Field), "retrieved document should have different value as original before refetch")
				sameObj = Complex64TestStructs.Find(newObj.ID).One().(*Complex64TestStruct)
				Expect(sameObj.Field).To(Equal(newObj.Field), "retrieved document should have same value as original after refetch")
			})
		})
		It("complex128 field [tb975eac8]", func() {
			type Complex128TestStruct struct {
				mongoid.Base `mongoid:"collection:int16_test"`
				ID           mongoid.ObjectID `bson:"_id"`
				Field        complex128
			}
			Complex128TestStructs := mongoid.Register(&Complex128TestStruct{Field: 0})
			newObj := Complex128TestStructs.New().(*Complex128TestStruct)
			newObj.Field = complex128(cmplx.Inf())
			OnlineDatabaseOnly(func() {
				newObj.Save()
				sameObj := Complex128TestStructs.Find(newObj.ID).One().(*Complex128TestStruct)
				Expect(sameObj.Field).To(Equal(newObj.Field), "retrieved document should have same value as original")
				newObj.Field = 7
				newObj.Save()
				Expect(sameObj.Field).ToNot(Equal(newObj.Field), "retrieved document should have different value as original before refetch")
				sameObj = Complex128TestStructs.Find(newObj.ID).One().(*Complex128TestStruct)
				Expect(sameObj.Field).To(Equal(newObj.Field), "retrieved document should have same value as original after refetch")
			})
		})
	})
})
