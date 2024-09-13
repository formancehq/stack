package client

import (
	"fmt"
)

type mangopayError struct {
	StatusCode int               `json:"-"`
	Message    string            `json:"Message"`
	Type       string            `json:"Type"`
	Errors     map[string]string `json:"Errors"`
}

func (me *mangopayError) Error() error {
	var errorMessage string
	if len(me.Errors) > 0 {
		for _, message := range me.Errors {
			errorMessage = message
			break
		}
	}

	var err error
	if errorMessage == "" {
		err = fmt.Errorf("unexpected status code: %d", me.StatusCode)
	} else {
		err = fmt.Errorf("%d: %s", me.StatusCode, errorMessage)
	}

	return err
}
