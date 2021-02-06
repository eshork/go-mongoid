package mongoid

import (
	"context"
	"runtime"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"mongoid/log"
	"mongoid/util"
)

// Find a document or multiple documents by their ids
func (model collectionHandle) Find(ids ...ObjectID) *Result {
	log.Debugf("%v.Find(%v)", model.GetModelName(), ids)
	return model.find(nil, ids...)
}

// FindCtx finds a document or multiple documents by their ids, bound by a new context
func (model collectionHandle) FindCtx(ctx context.Context, ids ...ObjectID) *Result {
	log.Debugf("%v.FindCtx(%v)", model.GetModelName(), ids)
	return model.find(ctx, ids...)
}

// FindByDeadline locates a document or multiple documents by their ids, bound by a time deadline.
func (model collectionHandle) FindByDeadline(d time.Time, ids ...ObjectID) *Result {
	log.Debugf("%v.FindByDeadline(%v)", model.GetModelName(), ids)
	ctx := model.Client().Context()
	newCtx, cancel := context.WithDeadline(ctx, d)
	res := model.find(newCtx, ids...)
	runtime.SetFinalizer(res, func(r *Result) {
		cancel()
	})
	return res
}

// FindByTimeout locates a document or multiple documents by their ids, bound by a timeout duration.
func (model collectionHandle) FindByTimeout(t time.Duration, ids ...ObjectID) *Result {
	log.Debugf("%v.FindByTimeout(%v)", model.GetModelName(), ids)
	ctx := model.Client().Context()
	newCtx, cancel := context.WithTimeout(ctx, t)
	res := model.find(newCtx, ids...)
	runtime.SetFinalizer(res, func(r *Result) {
		cancel()
	})
	return res
}

func (model *collectionHandle) find(ctx context.Context, ids ...ObjectID) *Result {
	modelContext := model.Client().Context()
	if ctx == nil {
		ctx = modelContext
	} else {
		ctx = util.ContextWithContext(ctx, modelContext)
	}

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
	cur, err := collection.Find(ctx, q)
	if err != nil {
		// this is a panic at the moment, because no one has yet looked to see what these errors might be, so we can't assume any of them are recoverable
		log.Panic(err) // unknown bad stuff happened within the driver
	}
	return makeResult(modelContext, cur, model)
}
