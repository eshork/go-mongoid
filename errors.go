package mongoid

// ErrorConfigNotFound means the specified configuration was not found
type ErrorConfigNotFound struct{}

func (err *ErrorConfigNotFound) Error() string {
	return "Specified Config was not found"
}

// ErrorConfigRefIsNil means the given Config reference was nil, and apparently that was bad
type ErrorConfigRefIsNil struct{}

func (err *ErrorConfigRefIsNil) Error() string {
	return "Given Config reference is Nil"
}

// ErrorAlreadyInitialized means mongoid was already initialized, so this is not a valid request
type ErrorAlreadyInitialized struct{}

func (err *ErrorAlreadyInitialized) Error() string {
	return "Already initialized"
}

// Database errors...

// ErrorNotConnected means mongoid client is not connected but use was attempted
type ErrorNotConnected struct{}

func (err *ErrorNotConnected) Error() string {
	return "Not Connected"
}

// ErrorConnectionTimedOut means mongoid client timed out while performing the action
type ErrorConnectionTimedOut struct{}

func (err *ErrorConnectionTimedOut) Error() string {
	return "Connection Timed Out"
}

// Validation errors...
