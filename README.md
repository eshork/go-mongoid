# go-mongoid

**a work in progress to make go webscale**

![alt text](etc/assets/go-mongoid-100.png "Mongoid for Go")

This is a (sort of) reimplementation of [Mongoid](https://github.com/mongodb/mongoid) for Go, using [mongo-go-driver](https://github.com/mongodb/mongo-go-driver) as the connection interface. The primary focus is on ease of use, convenience.

> Mongoid is an ODM (Object-Document Mapper) framework for MongoDB in ~Ruby~ Go

Many things (most things?) don't directly translate from Ruby to Go, but major themes of the Document-related interfaces from the Mongoid API are replicated as closely as possible, with adjustments as needed to facilitate language differences.

Also, I don't represent or work for MongoDB, Inc.

# Current features
- Map database documents onto native Go structs (struct fields map to BSON document fields)
- Supports all builtin Go data-types as document field-types
- Document lifecycle (persistence and change tracking)
- Find by ID or query builder interface

# Future features
- MongoDB connection configuration via JSON, YAML, or ENV vars
- Model relationships: one-to-one, one-to-many, many-to-many
- Custom callbacks based on document lifecycle events (onCreate, onUpdate, onDelete)
- Field validation callbacks
- Plugin architecture for adhoc add-on functionality (think Mongoid::Paranoia, Mongoid::Versioning, etc)


# Installation & Usage

Add the library to your project

```bash
cd ~/yourGoProjectDir
go get -u github.com/eshork/go-mongoid
```

Configure a MongoDB server connection

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

var MyDocuments = mongoid.Collection(&MyDocument{})
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

----

[MIT License](LICENSE)
