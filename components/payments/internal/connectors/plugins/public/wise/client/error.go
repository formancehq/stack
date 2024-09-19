package client

import (
	"fmt"
)

type wiseErrors struct {
	Errors []*wiseError `json:"errors"`
}

func (we *wiseErrors) Error(statusCode int) *wiseError {
	if len(we.Errors) == 0 {
		return &wiseError{StatusCode: statusCode}
	}
	we.Errors[0].StatusCode = statusCode
	return we.Errors[0]
}

type wiseError struct {
	StatusCode int    `json:"-"`
	Code       string `json:"code"`
	Message    string `json:"message"`
}

func (me *wiseError) Error() error {
	if me.Message == "" {
		return fmt.Errorf("unexpected status code: %d", me.StatusCode)
	}

	return fmt.Errorf("%s: %s", me.Code, me.Message)
}
