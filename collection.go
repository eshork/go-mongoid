package mongoid

import (
	"context"
	"time"
)

// ICollection provides methods to interact with Document records.
// Record access is scoped to a specific collection name within a specific database via a specific client.
// Use one of the With... methods to create a new collectionHandle with an updated scope.
//
type ICollection interface {
	String() string
	TypeName() string
	WithCollectionName(newCollectionName string) ICollection
	CollectionName() string
	WithDatabaseName(newDatabaseName string) ICollection
	DatabaseName() string
	WithClientName(newClientName string) ICollection
	ClientName() string
	Client() *Client

	New() IDocument

	DefaultBSON() BsonDocument
	Find(ids ...ObjectID) *Result
	FindCtx(ctx context.Context, ids ...ObjectID) *Result
	FindByDeadline(d time.Time, ids ...ObjectID) *Result
	FindByTimeout(t time.Duration, ids ...ObjectID) *Result
}
