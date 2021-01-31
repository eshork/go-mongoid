package errors

// IndexOutOfBounds can occur when an indexed retrieval function or method is given a value that is beyond the available data set.
// Ie, you probably asked for data that doesn't exist, because the index value you gave is probably too big.
type IndexOutOfBounds struct{}

var _ error = new(IndexOutOfBounds)
var _ error = IndexOutOfBounds{}
var _ MongoidError = new(IndexOutOfBounds)
var _ MongoidError = IndexOutOfBounds{}

// ErrIndexOutOfBounds is the value used when IndexOutOfBounds is given
var ErrIndexOutOfBounds error = IndexOutOfBounds{}

//IsIndexOutOfBounds returns true if the given err is ErrIndexOutOfBounds
func IsIndexOutOfBounds(err error) bool {
	if err, ok := err.(IndexOutOfBounds); ok {
		return err == ErrIndexOutOfBounds
	}
	if err, ok := err.(*IndexOutOfBounds); ok {
		return *err == ErrIndexOutOfBounds
	}
	return false
}

// Error implements error interface
func (err IndexOutOfBounds) Error() string {
	return "IndexOutOfBounds"
}

// mongoidError implements MongoidError interface
func (err IndexOutOfBounds) mongoidError() {}

// Unwrap implements MongoidError interface
func (err IndexOutOfBounds) Unwrap() error { return nil }
