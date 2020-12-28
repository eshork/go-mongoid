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

// ObjectID is the BSON ObjectID type, as implemented by mongo-go-driver.
// More details about ObjectID can be found within the MongoDB documentation:
// https://docs.mongodb.com/manual/reference/method/ObjectId/
type ObjectID = bsonPrimitive.ObjectID

// ZeroObjectID returns a zero-value ObjectID, which can be used as a sentinel value.
// Within the context of go-mongoid, a zero-value ObjectID is used as the initial value of the document `_id` field,
// and it will be replaced by a newly generated unique ObjectID value during the first Save() operation.
func ZeroObjectID() ObjectID {
	return ObjectID{}
}

// NewObjectID generates a new (populated) ObjectID that is ready to be used as a unique identifier for documents.
// This will never produce a zero-value ObjectID.
func NewObjectID() ObjectID {
	return bsonPrimitive.NewObjectID()
}

// NewObjectIDFromTimestamp generates a new ObjectID based on the given time
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
