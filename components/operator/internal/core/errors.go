package core

const (
	CodePending       = "pending"
	CodeDeleted       = "object deleted"
	CodeStackDisabled = "stack disabled"
	CodeStackNotFound = "stack not found"
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

func NewStackNotFoundError() *ApplicationError {
	return NewApplicationError(CodeStackNotFound)
}

func NewPendingError() *ApplicationError {
	return NewApplicationError(CodePending)
}

func IsApplicationError(err error) bool {
	_, ok := err.(*ApplicationError)
	return ok
}
