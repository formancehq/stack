package core

const (
	CodePending       = "pending"
	CodeDeleted       = "object deleted"
	CodeStackDisabled = "stack disabled"
)

type ApplicationError struct {
	code string
}

func (e *ApplicationError) Error() string {
	return e.code
}

func NewApplicationError(code string) *ApplicationError {
	return &ApplicationError{
		code: code,
	}
}

func NewDeletedError() *ApplicationError {
	return NewApplicationError(CodeDeleted)
}

func NewStackDisabledError() *ApplicationError {
	return NewApplicationError(CodeStackDisabled)
}

func NewPendingError() *ApplicationError {
	return NewApplicationError(CodePending)
}

func IsApplicationError(err error) bool {
	_, ok := err.(*ApplicationError)
	return ok
}
