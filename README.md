# go-mongoid

**a work in progress to make go webscale**


![alt text](etc/assets/go-mongoid-100.png "Mongoid for Go")

This is a (sort of) reimplementation of [Mongoid](https://github.com/mongodb/mongoid) for Go, using [mongo-go-driver](https://github.com/mongodb/mongo-go-driver) as the connection interface. The primary focus is on ease of use, convenience.

> Mongoid is an ODM (Object-Document Mapper) framework for MongoDB in ~Ruby~ Go

Many things (most things?) don't directly translate from Ruby to Go, but major themes of the Document-related interfaces from the Mongoid API are replicated as closely as possible, with adjustments as needed to facilitate language differences.

Also, I don't represent or work for MongoDB, Inc.

# Target features for v1.0.0

- Uses Go structs as the primary document interface - ie, build your own custom document definitions using native syntax
  - Supports all builtin Go data-types as document field-types
  - Supports custom structs as document field-types (embedded documents)
  - Supports maps and slices/arrays as dynamic/flexible field-types
  - Supports custom field data-types (custom structs with their own bson marshaling methods)
- Default values for new document objects
- Change tracking - identify which fields have been altered since new object creation or since loading from the database, as well as the previous values
- Atomic updates - only changed fields are written to the datastore during save operations, same as Ruby Mongoid
- Query builder interface - concatenating method calls to build complex queries

---
# Future features
- Save and recall query Scopes (as well as default scopes per ModelType)

- Model relationships: one-to-one, one-to-many, many-to-many (and the inverses)
  - Lazy loading for cross-document associations by default
  - Easy basis to spawn new custom Query builders

- Custom Callbacks based on document lifecycle events (onCreate, onUpdate, onDelete)

- Custom Validations for document lifecycle events (onCreate, onUpdate, onDelete)

- Plugin architecture allows for adhoc add-on functionality (think Mongoid::Paranoia, Mongoid::Versioning, etc)

- MongoDB connection configuration via JSON, YAML, or ENV vars


# Installation & Usage

Add the library to your project

```bash
cd ~/yourGoProjectDir
go get -u github.com/eshork/go-mongoid
```

Configure a MongoDB server

```
import mongoid "https://github.com/eshork/go-mongoid"

gMongoidConfig := mongoid.Config{
  Clients: []mongoid.Client{
    {
      Name:     "default",
      Hosts:    []string{"localhost:27017"},
      Database: "yourDbNameHere",
    },
  },
}
mongoid.Configure(&gMongoidConfig)
```

Define a bespoke struct and use it to access database records

```
type MyDocument struct {
	mongoid.Document
	MyValue string
}

var MyDocuments = mongoid.Model(&MyDocument{})
```

Make a new item and save it

```go
newDoc := MyDocuments.New().(*MyDocument)
newDoc.MyValue = "something worth keeping"
newDoc.Save()

var mongoid.ObjectID recordId = newDoc.ID
```

Retrieve stored records by ID

```go
foundDoc := MyDocuments.Find(recordId).One().(*MyDocument)
```

Check [the wiki](https://github.com/eshork/go-mongoid/wiki) for additional setup information and examples.

Refer to the [examples/](https://github.com/eshork/go-mongoid/tree/master/examples) directory for some use case examples to get you started.

Run `grift docs` to start a local godoc server to view the embedded source documentation.
