package mongoid

import (
	"context"
	mongoidError "mongoid/errors"
	"mongoid/log"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// IsPersisted returns true if the document has been saved to the database.
// Returns false if the document is new or has been destroyed.
// This is not a change tracker -- see: Changed() for that
func (d *Base) IsPersisted() bool {
	log.Trace("Base.IsPersisted()")
	return d.persisted
}

// sets the peristed state
func (d *Base) setPersisted(val bool) {
	d.persisted = val
}

// IsChanged returns true if the document instance has changed from the last retrieved value from the datastore, false otherwise.
// Newly created but documents begin in a Changed()=false state, as the document begins with all default values.
func (d *Base) IsChanged() bool {
	log.Trace("Base.IsChanged()")
	if len(d.Changes()) > 0 {
		return true
	}
	return false
}

// Changes gives a BsonDocument representation of all changes detected on the document/model instance.
// Entries will only exist when a change for that specific key/value is detected.
// Entries that are unchanged are excluded from the output BsonDocument.
// New or changed values will have a key/value pair that reflects the newly set entry value.
// Unset or missing values will have an key/value pair with 'nil' as the value side, to reflect the unset status.
func (d *Base) Changes() BsonDocument {
	log.Trace("Base.Changes()")
	currentBson := d.ToBson()
	previousBson := d.previousValue
	diffBson := makeBsonDocumentDiff(previousBson, currentBson)
	return diffBson
}

// Was provides the previous field value and indicates if a change has occurred
func (d *Base) Was(fieldPath string) (interface{}, bool) {
	value, err := d.GetField(fieldPath)
	if err != nil {
		log.Panic(err)
	}
	prevValue, prevFound := d.GetFieldPrevious(fieldPath)
	if !prevFound {
		log.Panic(mongoidError.DocumentFieldNotFound{FieldName: fieldPath})
	}
	if !verifyBothAreSameSame(value, prevValue) {
		// different types is definitely a change
		return prevValue, true
	}
	return prevValue, !reflect.DeepEqual(value, prevValue)
}

// Save will store the changed attributes to the database atomically, or insert the document if flagged as a new record via Model#new_record?
// Can bypass validations if wanted.
func (d *Base) Save() error {
	log.Debugf("%v.Save()", d.Model().modelName)

	// if already persisted, this is an update, otherwise it's a new insert
	if d.IsPersisted() {
		// update goes here
		return d.saveByUpdate()
	}
	return d.saveByInsert()
}

func (d *Base) saveByUpdate() error {
	log.Trace("saveByUpdate()")

	collection := d.getMongoCollectionHandle()
	ctx, ctxCancel := context.WithTimeout(context.TODO(), 5*time.Second) // todo context with arbitrary 5sec timeout
	defer ctxCancel()

	selectFilter := bson.M{"_id": d.GetID()}
	updateBson := d.ToUpdateBson()
	log.Debugf("collection[%s].UpdateOne %v %v", collection.Name(), selectFilter, updateBson)
	_, err := collection.UpdateOne(ctx, selectFilter, updateBson)
	if err != nil {
		log.Fatal(err)
	}
	d.setPersisted(true)         // this is now persisted
	d.refreshPreviousValueBSON() // update change tracking with current values
	return nil
}

func (d *Base) saveByInsert() error {
	log.Trace("saveByInsert()")
	// insert a new object

	collection := d.getMongoCollectionHandle()
	ctx, ctxCancel := context.WithTimeout(context.TODO(), 5*time.Second) // todo context with arbitrary 5sec timeout
	defer ctxCancel()

	insertBson := d.ToBson()
	// log.Error("insertBson: ", insertBson)

	// if there's a root _id field with a zero-value ObjectID, just drop it
	idObjInterface, found := insertBson["_id"]
	if found {
		objectID, ok := idObjInterface.(ObjectID)
		if ok && objectID == ZeroObjectID() {
			// REF ISSUE #19 - is it better to make a new ObjectID here, or let the MongoDB driver do it for us?

			// METHOD 1 - this way simply makes a new ObjectID here (ie within go-mongoid)
			// insertBson["_id"] = NewObjectID()

			// METHOD 2 - this way deletes _id field so that the Mongo driver will be forced to figure it out for us
			delete(insertBson, "_id")
		}
	}

	log.Debugf("collection[%s].InsertOne %v", collection.Name(), insertBson)
	res, err := collection.InsertOne(ctx, insertBson)
	if err != nil {
		log.Fatal(err)
	}

	id := res.InsertedID
	// log.Error("id: ", id)
	if err := d.SetField("_id", id); err != nil {
		log.Panic(err)
	}

	d.setPersisted(true)         // this is now persisted
	d.refreshPreviousValueBSON() // update change tracking with current values
	return nil
}

// returns a handle to the mongo driver collection for this document instance
func (d *Base) getMongoCollectionHandle() *mongo.Collection {
	dModel := d.Model()
	return dModel.getMongoCollectionHandle()
}
