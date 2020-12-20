package errors

// DocumentFieldNotFound -
type DocumentFieldNotFound struct {
	MongoidError
	FieldName string
}

func (err *DocumentFieldNotFound) Error() string {
	if err.FieldName != "" {
		return "DocumentFieldNotFound: " + err.FieldName
	}
	return "DocumentFieldNotFound"
}
