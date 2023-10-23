package client

import (
	"encoding/json"
	"fmt"
	"io"
)

type modulrError struct {
	StatusCode    int    `json:"-"`
	Field         string `json:"field"`
	Code          string `json:"code"`
	Message       string `json:"message"`
	ErrorCode     string `json:"errorCode"`
	SourceService string `json:"sourceService"`
}

func (me *modulrError) Error() error {
	if me.Message == "" {
		return fmt.Errorf("unexpected status code: %d", me.StatusCode)
	}

	return fmt.Errorf("%d: %s", me.StatusCode, me.Message)
}

func unmarshalError(statusCode int, body io.ReadCloser) *modulrError {
	var ces []modulrError
	_ = json.NewDecoder(body).Decode(&ces)

	if len(ces) == 0 {
		return &modulrError{
			StatusCode: statusCode,
		}
	}

	return &modulrError{
		StatusCode:    statusCode,
		Field:         ces[0].Field,
		Code:          ces[0].Code,
		Message:       ces[0].Message,
		ErrorCode:     ces[0].ErrorCode,
		SourceService: ces[0].SourceService,
	}
}
