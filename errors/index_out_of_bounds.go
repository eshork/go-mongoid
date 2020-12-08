package errors

// IndexOutOfBounds can occur when an indexed retrieval function or method is given a value that is beyond the available data set.
// Ie, you probably asked for data that doesn't exist, because the index value you gave is probably too big.
type IndexOutOfBounds struct {
	MongoidError
}

func (err *IndexOutOfBounds) Error() string {
	return "IndexOutOfBounds"
}
