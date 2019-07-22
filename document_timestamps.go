package mongoid

import (
	// "mongoid/log"
	"time"
)

// ITimestampCreated ...
type ITimestampCreated interface{}

// TimestampCreated ...
type TimestampCreated struct {
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

// ITimestampUpdated ...
type ITimestampUpdated interface{}

// TimestampUpdated ...
type TimestampUpdated struct {
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

// Timestamps ...
type Timestamps struct {
	TimestampCreated `bson:",inline"`
	TimestampUpdated `bson:",inline"`
}
