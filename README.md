# go-mongoid

**work in progress ([all work is currently happening here](https://github.com/eshork/go-mongoid/tree/v1dev))**


![alt text](etc/assets/go-mongoid-100.png "Mongoid for Go")


This is a reimplementation of [https://github.com/mongodb/mongoid](Mongoid) for Go, using the (still in progress) official [mongo-go-driver](https://github.com/mongodb/mongo-go-driver) as the connection interface.

Not everything perfectly translates from Ruby to Go, but the Document-related interfaces of the Mongoid API are replicated as closely as possible with adjustments as necessarily to facilitate language differences.

> **Note:** Preserving syntactic sugar and developer experience/productivity is held as a primary concern. This can result in liberal usage of Go's _reflect_ package, but it is arguable that servers are fast/cheap and the costs of performing reflection is more than paid by the dividends gained by developer productivity. Also, let's not forget that Go is typically quite fast in comparison to almost every runtime interpretted scripting language, so I'd be quite interested to see some real-world benchmarks against the Ruby Mongoid implementation as a comparison to identify "slowness" within this library.
>
> If your specific use-case is highly sensitive to latency, you might elect to start here but then later optimize the time sensitive bits of your application to a more direct driver interface (for which this library is intended to also happily provide direct access into).

# Target features (for v1)

- MongoDB connection configuration via JSON, YAML, or ENV vars

- Uses Go structs as the primary document interface - ie, build your own custom document definitions using native syntax
  - Supports all native Go data-types as document field-types
  - Supports custom field data-types (nyi - some initial work already done, but api needs firming up)
  - Supports generic Go maps as dynamic field-types

- Compatible with officially released mongo-go-driver versions
  - not yet true as of pre-v1; currently uses a commit-ref-id off master as certain bleeding-edge functionality is required/better; this will be resolved before a 1.0 release

- Change tracking - identify which fields have been altered since new object creation or since loading from the database, as well as the previous values

- Default values for new document objects

- Atomic updates - only changed fields are written to the datastore during save operations, same as Ruby Mongoid

- Query builder interface - concatenate method calls to dynamically build queries
  - Save and recall query Scopes (as well as default scopes per ModelType)

- Model relationships: one-to-one, one-to-many, many-to-many (and the inverses)
  - Lazy loading for cross-document associations by default
  - Easy basis to spawn new custom Query builders

- Custom Callbacks based on document lifecycle events (onCreate, onUpdate, onDelete)

- Custom Validations for document lifecycle events (onCreate, onUpdate, onDelete)

- Plugin architecture allows for adhoc add-on functionality (think: Mongoid::Paranoia, Mongoid::Versioning, etc)
