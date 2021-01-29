package errors

// MongoidError is the base type for all errors and panics generated within go-mongoid.
type MongoidError interface {
	error
	Unwrap() error
	mongoidError()
}

//IsMongoidError returns true if the given err is a MongoidError
func IsMongoidError(err error) bool {
	if _, ok := err.(MongoidError); ok {
		return true
	}
	return false
}
