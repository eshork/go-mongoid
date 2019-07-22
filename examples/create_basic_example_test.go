package examples

import (
	// "bce/incidental/dao/mongodb/mongoid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	// "time"
)

var _ = Describe("Create Basic Example", func() {
	Describe("createNewPet", func() {
		It("creates a new pet record", func() {
			By("Running the example to create a new record")
			newPetID := createNewPet()
			Expect(newPetID.IsZero()).To(BeFalse())
			By("Retrieving the new record")
			// findRes := Pets.Find(newPetID)
			// Expect(len(findRes) > 0 ).To(BeFalse())

		})
	})

	Describe("createNewPetByLateBinding", func() {
	})

})
