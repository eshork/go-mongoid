package errors

// OperationTimedOut can occur when attempting to perform operations after the current context expires or is canceled
type OperationTimedOut struct {
	Wrapped    error
	MethodName string
	Reason     string
}

var _ error = new(OperationTimedOut)
var _ error = OperationTimedOut{}
var _ MongoidError = new(OperationTimedOut)
var _ MongoidError = OperationTimedOut{}

//IsOperationTimedOut returns true if the given err is a OperationTimedOut
func IsOperationTimedOut(err error) bool {
	if _, ok := err.(OperationTimedOut); ok {
		return true
	}
	if _, ok := err.(*OperationTimedOut); ok {
		return true
	}
	return false
}

// Error implements error interface
func (err OperationTimedOut) Error() string {
	// example: "OperationTimedOut [struct.MethodName] - Reason goes here"
	msg := "OperationTimedOut"
	if err.MethodName != "" {
		msg = msg + " [" + err.MethodName + "]"
	}
	if err.Reason != "" {
		msg = msg + " - " + err.Reason
	}
	return msg
}

// mongoidError implements MongoidError interface
func (err OperationTimedOut) mongoidError() {}

// Unwrap implements MongoidError interface
func (err OperationTimedOut) Unwrap() error { return err.Wrapped }
