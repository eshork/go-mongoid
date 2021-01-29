package errors

// ResultNotExpected can occur when a query produces an unexpected result.
// For example, if a query was expected to return a single result, but more than one result was present, you may encounter ResultNotExpected
type ResultNotExpected struct{}

var _ error = new(ResultNotExpected)
var _ error = ResultNotExpected{}
var _ MongoidError = new(ResultNotExpected)
var _ MongoidError = ResultNotExpected{}

// ErrResultNotExpected is the value used when ResultNotExpected is given
var ErrResultNotExpected error = ResultNotExpected{}

//IsResultNotExpected returns true if the given err is ErrResultNotExpected
func IsResultNotExpected(err error) bool {
	if err, ok := err.(ResultNotExpected); ok {
		return err == ErrResultNotExpected
	}
	if err, ok := err.(*ResultNotExpected); ok {
		return *err == ErrResultNotExpected
	}
	return false
}

// Error implements error interface
func (err ResultNotExpected) Error() string {
	return "ResultNotExpected"
}

// mongoidError implements MongoidError interface
func (err ResultNotExpected) mongoidError() {}

// Unwrap implements MongoidError interface
func (err ResultNotExpected) Unwrap() error { return nil }
