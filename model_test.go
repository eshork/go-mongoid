package mongoid_test

import (
	"mongoid"
	// "github.com/brianvoe/gofakeit"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type TestModel1 struct {
	mongoid.Base
	ID mongoid.ObjectID `bson:"_id"`
}

var TestModel1s = mongoid.Register(&TestModel1{})

type TestModel2 struct {
	mongoid.Base `mongoid:"collection:otherCollection"`
	ID           mongoid.ObjectID `bson:"_id"`
}

var TestModel2s = mongoid.Register(&TestModel2{})

type TestModel3 struct {
	mongoid.Base `mongoid:"database:otherDatabase"`
	ID           mongoid.ObjectID `bson:"_id"`
}

var TestModel3s = mongoid.Register(&TestModel3{})

var _ = Describe("Model", func() {
	Context("TestModel1", func() {
		It("is verifyably registered", func() {
			By("struct name")
			Expect(mongoid.M("TestModel1")).ToNot(BeNil())
			By("example ref object")
			Expect(mongoid.M(&TestModel1{})).ToNot(BeNil())
			By("convenience handle")
			Expect(TestModel1s).ToNot(BeNil())
			By("convenience handle type verification")
			Expect(TestModel1s).To(BeAssignableToTypeOf(mongoid.Model(&TestModel1{})))
		})
		It("is accessible via either M and Model methods", func() {
			Expect(mongoid.M(&TestModel1{})).To(BeAssignableToTypeOf(mongoid.Model(&TestModel1{})))
		})
		It("reports default Collection name", func() {
			Expect(mongoid.M("TestModel1").GetCollectionName()).To(Equal("test_model_1"))
		})
		It("reports default Database name", func() {
			OnlineDatabaseOnly(func() {
				Expect(mongoid.M("TestModel1").GetDatabaseName()).To(Equal(mongoid.DefaultClient().Database))
			})
		})
	})
	Context("TestModel2", func() {
		It("reports struct-tag declared Collection name", func() {
			Expect(mongoid.M("TestModel2").GetCollectionName()).To(Equal("otherCollection"))
		})
	})
	Context("TestModel3", func() {
		It("reports struct-tag declared Database name", func() {
			Expect(mongoid.M("TestModel3").GetDatabaseName()).To(Equal("otherDatabase"))
		})
	})
})
