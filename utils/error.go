package utils

// Error …
type Error struct {
	message string
}

func (e *Error) Error() string {
	return e.message
}

// NewError …
func NewError(message string) error {
	return &Error{message}
}
