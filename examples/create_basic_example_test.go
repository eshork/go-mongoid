package examples

import (
	"mongoid"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Create Basic Example", func() {
	Describe("createNewPet", func() {
		It("creates a new pet record", func() {
			OnlineDatabaseOnly(func() {
				By("Running the example to create a new record")
				newPetID := createNewPet()
				Expect(newPetID.IsZero()).To(BeFalse())
				By("Retrieving the new record")
				foundPetID := findPetByID(newPetID)
				Expect(foundPetID).To(Equal(newPetID))
				By("Running the example again to create another new record")
				newPetID2 := createNewPet()
				By("Retrieving the new record along with the previous record")
				// findTwoPetsByID(newPetID2, newPetID) // test
				foundPetID1, foundPetID2 := findTwoPetsByID(newPetID2, newPetID)
				Expect([]mongoid.ObjectID{foundPetID1, foundPetID2}).To(ContainElement(newPetID))
				Expect([]mongoid.ObjectID{foundPetID1, foundPetID2}).To(ContainElement(newPetID2))
			})
		})

	})
})
