package examples

import (
	"mongoid"
)

func findPetByID(id mongoid.ObjectID) *Pet {
	foundPet := Pets.Find(id).One().(*Pet) // use Model.Find() to retrieve one record by _id, then One() to retrieve a single expected record
	return foundPet                        // return the record
	// if the record was not found, then One() will panic with a NotFound{} object
}

func findTwoPetsByID(id1 mongoid.ObjectID, id2 mongoid.ObjectID) (*Pet, *Pet) {
	res := Pets.Find(id1, id2) // Model.Find() can retrieve multiple records by _id
	pet1 := res.First().(*Pet) // since we expect only 2 results, we can use First() and Last()
	pet2 := res.Last().(*Pet)
	return pet1, pet2
}
