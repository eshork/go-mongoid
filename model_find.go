package mongoid

import (

	// "strings"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	// "go.mongodb.org/mongo-driver/mongo"
	"mongoid/log"
)

// Find a document or multiple documents by their ids
func (model *ModelType) Find(ids ...ObjectID) (*Result, error) {
	log.Debug("ModelType.Find()")
	q := bson.D{}
	if len(ids) <= 0 {
		// q = bson.D{} // default (empty) is already correct - ie, find all records
	} else if len(ids) == 1 {
		q = bson.D{primitive.E{Key: "_id", Value: ids[0]}}
	} else if len(ids) > 1 {
		q = bson.D{
			primitive.E{
				Key: "_id",
				Value: bson.D{
					primitive.E{
						Key:   "$in",
						Value: ids,
					},
				}}}
	}

	collection := model.getMongoCollectionHandle()
	// ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second) // todo context with arbitrary 5sec timeout
	// defer cancel()
	ctx := context.TODO() // todo context with unlimited timeout

	cur, err := collection.Find(ctx, q)
	if err != nil {
		// this is a panic at the moment, because no one has yet looked to see what these errors might be, so we can't assume any of them are recoverable
		log.Panic(err) // unknown bad stuff happened within the driver
	}

	return makeResult(ctx, cur, model), nil

	// // old stuff below
	// defer cur.Close(ctx)
	// for cur.Next(ctx) {
	// 	var result bson.M
	// 	err := cur.Decode(&result)
	// 	if err != nil {
	// 		log.Panic(err)
	// 	}
	// 	// log.Panic(result)
	// 	// do something with result...
	// 	//   Create new object instance
	// 	//   load values from bson.D into object instance
	// 	//   reset change tracking
	// 	// yield to caller
	// }
	// if err := cur.Err(); err != nil {
	// 	log.Panic(err)
	// }
	// return nil, nil
}

// func (model *ModelType) New() IDocumentBase {

// // Start Example 9
// cursor, err := coll.Find(
// 	context.Background(),
// 	bson.D{{"status", "D"}},
// )
// // End Example 9

// // Start Example 10
// cursor, err := coll.Find(
// 	context.Background(),
// 	bson.D{{"status", bson.D{{"$in", bson.A{"A", "D"}}}}})
// // End Example 10
