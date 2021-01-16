# go-mongoid

**a work in progress to make go webscale**


![alt text](etc/assets/go-mongoid-100.png "Mongoid for Go")

This is a (sort of) reimplementation of [Mongoid](https://github.com/mongodb/mongoid) for Go, using [mongo-go-driver](https://github.com/mongodb/mongo-go-driver) as the connection interface. The primary focus is on ease of use, convenience.

> Mongoid is an ODM (Object-Document Mapper) framework for MongoDB in ~Ruby~ Go

Many things (most things?) don't directly translate from Ruby to Go, but the major Document-related interfaces from the Mongoid API are replicated as closely as possible, with adjustments as needed to facilitate language differences.

# Target features (for v1)

- Uses Go structs as the primary document interface - ie, build your own custom document definitions using native syntax
  - Supports builtin Go data-types as document field-types
  - Supports custom structs as document field-types (embedded documents)
  - Supports maps and slices/arrays as dynamic field-types
  - Supports custom field data-types (custom structs with their own bson marshal/unmarshal methods)
- Default values for new document objects
- Change tracking - identify which fields have been altered since new object creation or since loading from the database, as well as the previous values
- Atomic updates - only changed fields are written to the datastore during save operations, same as Ruby Mongoid
- Query builder interface - concatenate method calls to dynamically build queries
  - Save and recall query Scopes (as well as default scopes per ModelType)

- Model relationships: one-to-one, one-to-many, many-to-many (and the inverses)
  - Lazy loading for cross-document associations by default
  - Easy basis to spawn new custom Query builders

- Custom Callbacks based on document lifecycle events (onCreate, onUpdate, onDelete)

- Custom Validations for document lifecycle events (onCreate, onUpdate, onDelete)

- Plugin architecture allows for adhoc add-on functionality (think Mongoid::Paranoia, Mongoid::Versioning, etc)

- MongoDB connection configuration via JSON, YAML, or ENV vars



# Installation & Usage

Within your project directory, add/update the library to the latest stable release version

```bash
go get -u github.com/eshork/go-mongoid
```

Within any relevant project files, import the package as usual
```go
// Add this basic import to your project file
import mongoid "github.com/eshork/go-mongoid"
```

Refer to the `examples` directory for use case examples to get you started.

Run `grift docs` to start a godoc server to view embedded documentation.

Be sure to check [the wiki](https://github.com/eshork/go-mongoid/wiki) for additional information and examples.
