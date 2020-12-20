package errors

// MongoidError is the base type for all errors and panics generated within go-mongoid.
type MongoidError interface {
	error
}
