package examples

import (
	"log"
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

func findPetByID(id mongoid.ObjectID) mongoid.ObjectID {
	// res, _ := Pets.Find(id)              // use Model.Find() to retrieve one or more records by _id
	// foundPet := res.OneAndClose().(*Pet) // a Result can contain several records, but we're only expecting one, so just grab the First() and cast it back into a *Pet
	foundPet := Pets.Find(id).OneAndClose().(*Pet)
	return foundPet.ID // return back some (weak) proof that we found the record
}

func findTwoPetsByID(id1 mongoid.ObjectID, id2 mongoid.ObjectID) (mongoid.ObjectID, mongoid.ObjectID) {
	res := Pets.Find(id1, id2) // Model.Find() can retrieve multiple records by _id
	defer res.Close()          // be sure to Close the Result when you're done with it
	pet1 := res.First().(*Pet) // since we expect only 2 results, we can use First() and Last()
	pet2 := res.Last().(*Pet)  // and it doesn't matter which order we use them in
	return pet1.ID, pet2.ID
}

func panicOnFindError(res *mongoid.Result, err error) *mongoid.Result {
	if err != nil {
		log.Panic(err)
	}
	return res
}
