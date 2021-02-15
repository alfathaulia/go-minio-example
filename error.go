package upload

import (
	"fmt"
)

//ErrorCode errors
type ErrorCode int

//Role enum of user permission roles
const (
	BaseError ErrorCode = iota
	RecordNotFound
	ParserError
	NoSuchFile
	DBError
)

func (e ErrorCode) String() string {
	return [...]string{
		"Error",
		"RecordNotFound",
		"ParserError",
		"NoSuchFile",
		"DBError",
	}[e]
}

//Error uploader errors
type Error struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

func (err *Error) Error() string {
	return fmt.Sprintf("%s: %s, with code: %v", ErrorCode(err.Code).String(), err.Message, err.Code)
}
