package errors

// InvalidOperation can occur when a function or method is called in an unexpected manner
// (methods called out of expected/intended order, using uninitialized subsystems, unusable inputs etc)
type InvalidOperation struct {
	Wrapped    error
	MethodName string
	Reason     string
}

var _ error = new(InvalidOperation)
var _ error = InvalidOperation{}
var _ MongoidError = new(InvalidOperation)
var _ MongoidError = InvalidOperation{}

//IsInvalidOperation returns true if the given err is a InvalidOperation
func IsInvalidOperation(err error) bool {
	if _, ok := err.(InvalidOperation); ok {
		return true
	}
	if _, ok := err.(*InvalidOperation); ok {
		return true
	}
	return false
}

// Error implements error interface
func (err InvalidOperation) Error() string {
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

// mongoidError implements MongoidError interface
func (err InvalidOperation) mongoidError() {}

// Unwrap implements MongoidError interface
func (err InvalidOperation) Unwrap() error { return err.Wrapped }
