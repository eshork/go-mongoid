package examples

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Create Basic Example", func() {
	Describe("createNewPet", func() {
		PIt("creates a new pet record", func() {
			OnlineDatabaseOnly(func() {
				By("Running the example to create a new record")
				newPetID := createNewPet()
				Expect(newPetID.IsZero()).To(BeFalse())
				By("Retrieving the new record")
				foundPetID := findPetByID(newPetID)
				Expect(foundPetID).To(Equal(newPetID))
			})
		})
	})
})
