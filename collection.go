package mongoid

import (
	"context"
	"time"
)

// ICollection ...
type ICollection interface {
	String() string
	Name() string
	WithCollectionName(newCollectionName string) ICollection
	GetCollectionName() string
	WithDatabaseName(newDatabaseName string) ICollection
	GetDatabaseName() string
	WithClientName(newClientName string) ICollection
	GetClientName() string
	GetClient() *Client

	New() IDocument
	GetDefaultBSON() BsonDocument
	Find(ids ...ObjectID) *Result
	FindCtx(ctx context.Context, ids ...ObjectID) *Result
	FindByDeadline(d time.Time, ids ...ObjectID) *Result
	FindByTimeout(t time.Duration, ids ...ObjectID) *Result
}
