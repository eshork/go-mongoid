package errors

// ResultNotFound can occur when a query did not produce a result when it was was expected to produce one
type ResultNotFound struct{}

var _ error = new(ResultNotFound)
var _ error = ResultNotFound{}
var _ MongoidError = new(ResultNotFound)
var _ MongoidError = ResultNotFound{}

// ErrResultNotFound is the value used when ResultNotFound is given
var ErrResultNotFound error = ResultNotFound{}

//IsResultNotFound returns true if the given err is ErrResultNotFound
func IsResultNotFound(err error) bool {
	if err, ok := err.(ResultNotFound); ok {
		return err == ErrResultNotFound
	}
	if err, ok := err.(*ResultNotFound); ok {
		return *err == ErrResultNotFound
	}
	return false
}

// Error implements error interface
func (err ResultNotFound) Error() string {
	return "ResultNotFound"
}

// mongoidError implements MongoidError interface
func (err ResultNotFound) mongoidError() {}

// Unwrap implements MongoidError interface
func (err ResultNotFound) Unwrap() error { return nil }
