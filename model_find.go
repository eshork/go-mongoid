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
func (model *ModelType) Find(ids ...ObjectID) *Result {
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
	return makeResult(ctx, cur, model)
}
