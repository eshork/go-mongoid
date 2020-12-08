package errors

// NotFound can occur when a query did not produce a result when it was was expected to do so
type NotFound struct {
	MongoidError
}

func (err *NotFound) Error() string {
	return "NotFound"
}
