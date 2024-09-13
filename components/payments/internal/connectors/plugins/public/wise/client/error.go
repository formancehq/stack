package client

import (
	"encoding/json"
	"fmt"
	"io"
)

type wiseErrors struct {
	Errors []*wiseError `json:"errors"`
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

func unmarshalError(statusCode int, body io.ReadCloser) *wiseError {
	var ces wiseErrors
	_ = json.NewDecoder(body).Decode(&ces)

	if len(ces.Errors) == 0 {
		return &wiseError{
			StatusCode: statusCode,
		}
	}

	return &wiseError{
		StatusCode: statusCode,
		Code:       ces.Errors[0].Code,
		Message:    ces.Errors[0].Message,
	}
}
