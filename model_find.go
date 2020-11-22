package mongoid

import (
	// "fmt"
	// "mongoid/log"
	// "reflect"
	// "strings"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/mongo"
	"mongoid/log"
	"time"
)

func (model *ModelType) Find(ids ...ObjectID) (*Result, error) {
	log.Debug("ModelType.Find()")

	collection := model.getMongoCollectionHandle()
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second) // todo context with arbitrary 5sec timeout
	defer cancel()
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		log.Panic(err)
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var result bson.M
		err := cur.Decode(&result)
		if err != nil {
			log.Panic(err)
		}
		log.Panic(result)
		// do something with result....
	}
	if err := cur.Err(); err != nil {
		log.Panic(err)
	}
	return nil, nil
}

// // IsPersisted returns true if the document has been saved to the database.
// // Returns false if the document is new or has been destroyed.
// // This is not a change tracker -- see: Changed() for that
// func (d *Base) IsPersisted() bool {
// 	log.Debug("Base.IsPersisted()")
// 	return d.persisted
// }

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
