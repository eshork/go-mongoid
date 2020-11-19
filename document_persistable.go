package mongoid

import (
	"context"
	"mongoid/log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

// IsPersisted returns true if the document has been saved to the database.
// Returns false if the document is new or has been destroyed.
// This is not a change tracker -- see: Changed() for that
func (d *Base) IsPersisted() bool {
	log.Debug("Base.IsPersisted()")
	return d.persisted
}

// IsChanged returns true if the document instance has changed from the last retrieved value from the datastore, false otherwise.
// Newly created but documents begin in a Changed()=false state, as the document begins with all default values.
func (d *Base) IsChanged() bool {
	log.Debug("Base.IsChanged()")
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
	log.Debug("Base.Changes()")
	currentBson := d.ToBson()
	previousBson := d.previousValue
	diffBson := makeBsonDocumentDiff(previousBson, currentBson)
	return diffBson
}

// Was provides the previous field value and indicates if a change has occurred
func (d *Base) Was(fieldPath string) (interface{}, bool) {
	log.Panicf("NYI -Base.Was_(%s)", fieldPath)
	return nil, false
}

// Save will store the changed attributes to the database atomically, or insert the document if flagged as a new record via Model#new_record?
// Can bypass validations if wanted.
func (d *Base) Save() error {
	log.Debug("Base.Save()")

	// if already persisted, this is an update, otherwise it's a new insert
	if d.IsPersisted() {
		// update goes here
		return d.saveByUpdate()
	}
	return d.saveByInsert()
}

func (d *Base) saveByUpdate() error {
	log.Debug("saveByUpdate()")
	// insert a new object
	log.Fatal("NYI Save() - PERSISTED")
	return nil
}

func (d *Base) saveByInsert() error {
	log.Debug("saveByInsert()")
	// insert a new object

	collection := d.getMongoCollectionHandle()
	ctx, _ := context.WithTimeout(context.TODO(), 5*time.Second) // todo context with arbitrary 5sec timeout

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

	res, err := collection.InsertOne(ctx, insertBson)
	if err != nil {
		log.Fatal(err)
	}

	id := res.InsertedID
	// log.Error("id: ", id)
	if err := d.SetField("_id", id); err != nil {
		log.Error(err)
		log.Panic(err)
	}

	d.persisted = true           // this is now persisted
	d.refreshPreviousValueBSON() // update change tracking with current values
	return nil
}

// returns a handle to the mongo driver collection for this document instance
func (d *Base) getMongoCollectionHandle() *mongo.Collection {
	dModel := d.Model()
	return dModel.getMongoCollectionHandle()
}
