package seng

// Error Built-in error
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Error implements error interface
func (e *Error) Error() string {
	return e.Message
}

func NewError(code int, message ...string) *Error {
	e := &Error{
		Code: code,
	}
	if len(message) > 0 {
		e.Message = message[0]
	} else {
		e.Message = StatusMessage(code)
	}
	return e
}
