package response

import "errors"

// ErrorSomethingWentWrong http api errors
var (
	ErrorSomethingWentWrong error = errors.New("something went wrong")
	ErrUnauthorized               = errors.New("user is unauthorized")
	TokenExpired                  = errors.New("token is expired")
)

type Error struct {
	error
	StatusCode int
	Message    string
}

func (err *Error) Error() string {
	return err.Message
}

func NewError(statusCode int, message string) *Error {
	return &Error{
		StatusCode: statusCode,
		Message:    message,
	}
}
