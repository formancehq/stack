package sharedapi

import "encoding/json"

type BaseResponse[T any] struct {
	Data   *T         `json:"data,omitempty"`
	Cursor *Cursor[T] `json:"cursor,omitempty"`
}

type Cursor[T any] struct {
	PageSize int    `json:"page_size,omitempty"`
	HasMore  bool   `json:"has_more"`
	Previous string `json:"previous,omitempty"`
	Next     string `json:"next,omitempty"`
	Data     []T    `json:"data"`
}

type cursor[T any] Cursor[T]

func (c Cursor[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		cursor[T]
		// Keep those fields to ensure backward compatibility, even if it will be
		Total     int64 `json:"total,omitempty"`
		Remaining int   `json:"remaining_results,omitempty"`
	}{
		cursor: cursor[T](c),
	})
}

type ErrorResponse struct {
	ErrorCode    string `json:"error_code,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
}
