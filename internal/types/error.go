package types

import "fmt"

type ErrorWithCode interface {
	error
	fmt.Stringer
	ErrorCode() int
	Error() string
	ErrorMessage() string
}