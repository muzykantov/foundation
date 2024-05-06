package reflect

import "errors"

// Error types.
var (
	ErrMethodNotFound          = errors.New("method not found")
	ErrIncorrectArgumentCount  = errors.New("incorrect number of arguments")
	ErrUnsupportedArgumentType = errors.New("unsupported argument type")
	ErrInvalidArgumentValue    = errors.New("invalid argument value")
)
