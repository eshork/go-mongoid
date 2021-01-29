package errors

// DocumentFieldNotFound -
type DocumentFieldNotFound struct {
	FieldName string
}

var _ error = new(DocumentFieldNotFound)
var _ error = DocumentFieldNotFound{}
var _ MongoidError = new(DocumentFieldNotFound)
var _ MongoidError = DocumentFieldNotFound{}

//IsDocumentFieldNotFound returns true if the given err is a DocumentFieldNotFound
func IsDocumentFieldNotFound(err error) bool {
	if _, ok := err.(DocumentFieldNotFound); ok {
		return true
	}
	if _, ok := err.(*DocumentFieldNotFound); ok {
		return true
	}
	return false
}

// Error implements error interface
func (err DocumentFieldNotFound) Error() string {
	if err.FieldName != "" {
		return "DocumentFieldNotFound: " + err.FieldName
	}
	return "DocumentFieldNotFound"
}

// mongoidError implements MongoidError interface
func (err DocumentFieldNotFound) mongoidError() {}

// Unwrap implements MongoidError interface
func (err DocumentFieldNotFound) Unwrap() error { return nil }
