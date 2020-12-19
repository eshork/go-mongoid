package errors

// ResultNotExpected can occur when a query produces an unexpected result.
// For example, if a query was expected to return a single result, but more than one result was present, you may encounter ResultNotExpected
type ResultNotExpected struct {
	MongoidError
}

func (err *ResultNotExpected) Error() string {
	return "ResultNotExpected"
}
