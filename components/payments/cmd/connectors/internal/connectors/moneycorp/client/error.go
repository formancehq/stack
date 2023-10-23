package client

import (
	"encoding/json"
	"fmt"
	"io"
)

type moneycorpErrors struct {
	Errors []*moneycorpError `json:"errors"`
}

type moneycorpError struct {
	StatusCode int    `json:"-"`
	Code       string `json:"code"`
	Title      string `json:"title"`
	Detail     string `json:"detail"`
}

func (me *moneycorpError) Error() error {
	if me.Detail == "" {
		return fmt.Errorf("unexpected status code: %d", me.StatusCode)
	}

	return fmt.Errorf("%d: %s", me.StatusCode, me.Detail)
}

func unmarshalError(statusCode int, body io.ReadCloser) *moneycorpError {
	var ces moneycorpErrors
	_ = json.NewDecoder(body).Decode(&ces)

	if len(ces.Errors) == 0 {
		return &moneycorpError{
			StatusCode: statusCode,
		}
	}

	return &moneycorpError{
		StatusCode: statusCode,
		Code:       ces.Errors[0].Code,
		Title:      ces.Errors[0].Title,
		Detail:     ces.Errors[0].Detail,
	}
}
