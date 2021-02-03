package examples

import (
	"mongoid"
	"time"
)

// Pet is a record for a pet in our adoption app.
// By default, Mongoid will store these records in a collection named after the pluralized form of the struct name. ('pets' in this case)
type Pet struct {
	mongoid.Document   // add Document functionality to this struct
	mongoid.Timestamps // add the optional automatic record timestamp fields (created_at, updated_at)

	// Every document has an ID field.
	// Explicitly declaring the ID in your struct allows you to access the value and control the data type (though ObjectID is a solid choice)
	ID mongoid.ObjectID `bson:"_id"`

	// Mongoid automatically determines the field storage type based on the struct field type.
	// You can optionally specify the field name to be used in the database via the bson struct tag, otherwise exported struct fields will be
	// stored as a snake_case variant of their struct field name.
	Name string `bson:"name"`

	// You can optionally exclude exported struct fields from the database by setting the bson name to "-"
	UnstoredField string `bson:"-"`

	// All unexported struct fields (non-capitalized) will not be stored to the database.
	otherUnstoredField string

	// If a database field is expected to be null-able, use a pointer type in your struct.
	// On record loading, if the database stored value is null, the field value in the struct object will be set to nil.
	// Otherwise, an instance of the field-type will be created and the field pointer will point to the object.
	AdoptionDate *time.Time
}

// Call the Model function with an example object to obtain a ModelType that can be used to create, save,
// and find records based on that struct.
var Pets = mongoid.Model(&Pet{
	Name: "spot", // the current field values will be used as default values for your Documents
})

// Some ModelType attributes can be changed via the With... methods, which return a new ModelType.
// For example, the WithClientName method can be used to work with records across multiple database server connections.
var LostPets = Pets.WithCollectionName("lost_pets").WithDatabaseName("LostPetsDB")
