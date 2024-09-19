package client

import (
	"fmt"
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
	var err error
	if me.Message == "" {
		err = fmt.Errorf("unexpected status code: %d", me.StatusCode)
	} else {
		err = fmt.Errorf("%d: %s", me.StatusCode, me.Message)
	}

	return err
}
