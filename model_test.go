package mongoid_test

import (
	"mongoid"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Model()", func() {
	type ModelTest struct {
		mongoid.Document
		ID           mongoid.ObjectID `bson:"_id"`
		ExampleField string
	}
	Context("document struct by value", func() {
		It("accepts inscope variable", func() {
			var doc ModelTest = ModelTest{ExampleField: "by value variable!"}
			m := mongoid.Model(doc)
			By("checking default value")
			Expect(m.New().(*ModelTest).ExampleField).To(Equal(doc.ExampleField), "expects default value to be preset on new object")
		})
		It("accepts inline", func() {
			const exampleStr = "by value inline!"
			m := mongoid.Model(ModelTest{ExampleField: exampleStr})
			By("checking default value")
			Expect(m.New().(*ModelTest).ExampleField).To(Equal(exampleStr), "expects default value to be preset on new object")
		})
	})
	Context("document struct by reference", func() {
		It("accepts inscope variable", func() {
			var docPtr *ModelTest = &ModelTest{ExampleField: "by reference variable!"}
			m := mongoid.Model(docPtr)
			By("checking default value")
			Expect(m.New().(*ModelTest).ExampleField).To(Equal(docPtr.ExampleField), "expects default value to be preset on new object")
		})
		It("accepts 'new' allocated variable", func() {
			var docPtr *ModelTest = new(ModelTest)
			docPtr.ExampleField = "by new'ed variable!"
			m := mongoid.Model(docPtr)
			By("checking default value")
			Expect(m.New().(*ModelTest).ExampleField).To(Equal(docPtr.ExampleField), "expects default value to be preset on new object")
		})
		It("accepts inline", func() {
			const exampleStr = "by reference inline!"
			m := mongoid.Model(&ModelTest{ExampleField: exampleStr})
			By("checking default value")
			Expect(m.New().(*ModelTest).ExampleField).To(Equal(exampleStr), "expects default value to be preset on new object")
		})
	})
})
