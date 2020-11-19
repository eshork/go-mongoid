package examples

import "mongoid"

func createNewPet() mongoid.ObjectID {
	newPet := Pets.New().(*Pet) // you _must_ use the New() method to create new document model objects (this maintains lifecycle for validation/callbacks/etc)

	// note: if you create a document/model object manually via new() or via stack object, it won't be linked into mongoid, and attempts to use it with mongoid will likely panic

	newPet.Name = "scruffy" // you can access struct fields as you normally would for both read and write
	newPet.Save()           // Save() will store the document to the database
	return newPet.ID        // If an ID was not explicitly provided, one will be automatically created
}

func findPetByID(id mongoid.ObjectID) mongoid.ObjectID {
	res, _ := Pets.Find(id)        // use Model.Find() to retrieve one or more records by _id
	foundPet := res.First().(*Pet) // a Result can contain several records, but we're only expecting one, so just grab the First() and cast it back into a *Pet
	return foundPet.ID             // return back some (weak) proof that we found the record
}
