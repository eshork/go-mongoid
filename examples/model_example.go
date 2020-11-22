package examples

import (
	"mongoid"
	"time"
)

// Pet is a record for a pet in our adoption app.
// Mongoid will automatically store these records in a collection named after the pluralized form of the struct name. ('pets' in this case)
type Pet struct {
	mongoid.Base       // add Mongoid functionality to this struct
	mongoid.Timestamps // add the optional automatic record timestamp fields (created_at, updated_at)

	// Every document has an ID field.
	// Explicitly declaring the ID in your struct allows you to access the value and control the data type (though ObjectID is a solid choice)
	ID mongoid.ObjectID `bson:"_id"`

	// Mongoid automatically determines the field storage type based on the struct field type.
	// You can optionally specify the field name to be used in the database via the bson struct tag, otherwise exported struct fields will be
	// stored as a snake_case variant of their struct field name.
	Name string `bson:"name"`

	// You can optionally exclude exported fields from the database by setting the bson name to "-"
	UnstoredField string `bson:"-"`

	// All unexported struct fields (non-capitalized) will not be stored to the database unless a name is explicitly given via the bson tag
	otherUnstoredField string `bson:"stored_anyhow"`

	// For many data types, default values can be specified via the "default" struct tag. These will be automatically applied to objects
	// created via the ModelType.New() or ModelType.Create() actions at instantiation. Objects intantiated natively will not have a default
	// value applied (there's no reliable way to intercept direct object allocation)
	Breed string `default:"mutt"`

	// If a database field is expected to be null-able, use a pointer type in your struct.
	// On record loading, if the database stored value is null, the field value in the struct object will be set to nil.
	// Otherwise, an instance of the field-type will be created and the field pointer will point to the object.
	// (This is the only way that null values are stored in the database. Concrete type zero-values in Go are too uncertain to be used as indicators across data types.)
	AdoptionDate *time.Time
}

// Pets model registration must occur before objects can be stored to or read from the database; all it takes is a pointer to a reference object of the desired type.
// You can register models before or after database connections have been established (via `mongoid.Configure()`), but the recommendation is to register models prior.
//
// The following example is a typical way to register a model type while storing a global convenience handle for future use.
// rather than name, so name overlap is technically permitted (convenience handles like this one make those situations easier to deal with)
var Pets = mongoid.Register(&Pet{})

// If you'd rather not create a convenience handle at the module-global level, you can simply discard the return value like so.
//     var _ = mongoid.Register(&Pet{})
// Without a convenience handle, you'll need to use the `mongoid.Model(ref_or_name)` method to retrieve a handle when you need one.
// Note: Model names must be globally unique to reliably use the string-based lookup function `mongoid.M("MyModelName")`, but models are ultimately registered by actual type,

func init() {
	// You can register models within an init function as well, if that's preferable to you.
	// mongoid.Register(&Pet{})

	// Once a model is registered, you can declare one or more indexes for it, which will be automatically created if missing once database connections are established (via background index creation so as not to block)
	// Pets.Index( /* indexDefinition */ ) // TODO

	// You can reconfigure various options after a model has been registered (once mongoid.Configured() is true, these may generate warning log messages)
	// For instace you can change the name of the collection that will be used
	Pets = Pets.SetCollectionName("our_pets")

	// You can even redefine the model name as used by `mongoid.M("MyModelName")` and  `mongoid.Model("MyModelName")`
	Pets = Pets.SetModelName("OurPets")

	// Take note that the above registration changes also updated the convenience handle value - these are stored by value, so you can create different convenience handles based on different application needs

	// We'll typically be using the convenience handle that we created above within the other example files, but model_example_test.go demonstrates these registration changes did occur.
}
