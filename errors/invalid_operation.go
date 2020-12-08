package errors

// InvalidOperation can occur when a function or method is called in an unexpected manner
// (methods called out of expected/intended order, using uninitialized subsystems, unusable inputs etc)
type InvalidOperation struct {
	MongoidError
	MethodName string
	Reason     string
}

func (err *InvalidOperation) Error() string {
	// example: "InvalidOperation [struct.MethodName] - Reason goes here"
	msg := "InvalidOperation"
	if err.MethodName != "" {
		msg = msg + " [" + err.MethodName + "]"
	}
	if err.Reason != "" {
		msg = msg + " - " + err.Reason
	}
	return msg
}
