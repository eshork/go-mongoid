package errors

// ResultNotFound can occur when a query did not produce a result when it was was expected to produce one
type ResultNotFound struct {
	MongoidError
}

func (err *ResultNotFound) Error() string {
	return "ResultNotFound"
}
