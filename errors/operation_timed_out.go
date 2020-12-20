package errors

// OperationTimedOut can occur when a function or method is called in an unexpected manner
// (methods called out of expected/intended order, using uninitialized subsystems, unusable inputs etc)
type OperationTimedOut struct {
	MongoidError
	MethodName string
	Reason     string
}

func (err *OperationTimedOut) Error() string {
	// example: "InvalidOperation [struct.MethodName] - Reason goes here"
	msg := "OperationTimedOut"
	if err.MethodName != "" {
		msg = msg + " [" + err.MethodName + "]"
	}
	if err.Reason != "" {
		msg = msg + " - " + err.Reason
	}
	return msg
}
