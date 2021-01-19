package examples

import (
	"mongoid"
)

func createNewPet() mongoid.ObjectID {
	newPet := Pets.New().(*Pet) // You must use the New() method to create new document model objects.
	//                             This allows mongoid to perform some necessary object initialization
	//                             and maintains lifecycle for validation/callbacks/etc.
	//
	//                             Use type assertion to cast returned Document interface objects into their original type

	// note: if you create a document/model object manually (via the golang new() operator or via stack object),
	//       then it won't be properly linked with mongoid, and attempts to use it with mongoid will likely panic
	// badPet := new(Pet) <- So don't do this
	// badPet := &Pet{}   <- And don't do this

	newPet.Name = "scruffy" // you can access struct fields as you normally would for both read and write
	newPet.Save()           // Save() will store this new document to the database
	return newPet.ID        // If an ID was not explicitly provided, one will be automatically created (ObjectID type) and can be read back after the document is Save()'ed
}
