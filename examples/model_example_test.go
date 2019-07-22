package examples

import (
	"mongoid"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	// "time"
)

// Verify model is registered
// via convenience reference
// via example object

// Verify model collection name

// Verify model registration name

type unknownModel struct {
	mongoid.Base
}

// Verify unknown model is not registered
// via example object
// via model registration name

// "fmt"
// "go.mongodb.org/mongo-driver/bson"
// "strconv"

var _ = Describe("Model Example", func() {

	Describe("Pet model registration", func() {
		It("has a non-nil convenience handle", func() {
			Expect(Pets).ToNot(BeNil())
		})
		It("is findable by example object", func() {
			By("mongoid.Model")
			petModelObj := mongoid.Model(&Pet{})
			Expect(petModelObj).ToNot(BeNil())
			By("mongoid.M")
			petMObj := mongoid.Model(&Pet{})
			Expect(petMObj).ToNot(BeNil())
			By("equality check")
			Expect(petModelObj).To(Equal(petMObj))
		})

		It("is findable by registration name", func() {
			By("mongoid.Model(&Pet{})")
			petModelObj := mongoid.Model(&Pet{})
			By("mongoid.Model(\"OurPets\")")
			petNamedObj := mongoid.Model("OurPets")
			By("equality check")
			Expect(petNamedObj).To(Equal(petModelObj))
		})

		It("has an updated model registration name", func() {
			Expect(mongoid.Model("OurPets").GetModelName()).To(Equal("OurPets"))
		})

		It("has an updated default collection name", func() {
			Expect(mongoid.Model("OurPets").GetCollectionName()).To(Equal("our_pets"))
		})

	})

})
