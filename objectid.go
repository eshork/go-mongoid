package mongoid

/*
These are convenience functions for working with ObjectID values.
Most of these are a nearly 1:1 pass-thru to the mongo-go-driver's implementation, found here:
https://github.com/mongodb/mongo-go-driver/blob/master/bson/primitive/objectid.go
*/

import (
	"time"

	bsonPrimitive "go.mongodb.org/mongo-driver/bson/primitive"
)

// ObjectID ...
// TODO: consider moving this to `type ObjectID = bsonPrimitive.ObjectID` instead
type ObjectID = bsonPrimitive.ObjectID

// ZeroObjectID is a zero-value ObjectID
func ZeroObjectID() ObjectID {
	return ObjectID{}
}

// NewObjectID generates a new ObjectID
func NewObjectID() ObjectID {
	return ObjectID(bsonPrimitive.NewObjectID())
}

//NewObjectIDFromTimestamp generates a new ObjectID based on the given time
func NewObjectIDFromTimestamp(timestamp time.Time) ObjectID {
	// 	return ObjectID(bsonPrimitive.NewObjectIDFromTimestamp(timestamp))
	return bsonPrimitive.NewObjectIDFromTimestamp(timestamp)
}

// NewObjectIDFromHex creates a new ObjectID from a hex string. It returns an error if the hex string is not a
// valid ObjectID.
func NewObjectIDFromHex(s string) (ObjectID, error) {
	// 	id, err := bsonPrimitive.ObjectIDFromHex(s)
	// 	return ObjectID(id), err
	return bsonPrimitive.ObjectIDFromHex(s)
}

// IsZero returns true if id is the empty ObjectID.
// func (id ObjectID) IsZero() bool {
// 	return id.primitive().IsZero()
// }

// Timestamp extracts the time part of the ObjectId.
// func (id ObjectID) Timestamp() time.Time {
// 	return id.primitive().Timestamp()
// }

// Hex returns the hex encoding of the ObjectID as a string
// func (id ObjectID) Hex() string {
// 	return id.primitive().Hex()
// }

// func (id ObjectID) String() string {
// 	return id.primitive().String()
// }

// MarshalJSON returns the ObjectID as a string
// func (id ObjectID) MarshalJSON() ([]byte, error) {
// 	return id.primitive().MarshalJSON()
// }

// func (id ObjectID) primitive() bsonPrimitive.ObjectID {
// 	return bsonPrimitive.ObjectID(id)
// }

// UnmarshalJSON populates the byte slice with the ObjectID. If the byte slice is 64 bytes long, it
// will be populated with the hex representation of the ObjectID. If the byte slice is twelve bytes
// long, it will be populated with the BSON representation of the ObjectID. Otherwise, it will
// return an error.
// func (id *ObjectID) UnmarshalJSON(b []byte) error {
// 	return (*bsonPrimitive.ObjectID)(id).UnmarshalJSON(b)
// }
