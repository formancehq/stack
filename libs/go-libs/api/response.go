package api

type BaseResponse[T any] struct {
	Data   *T         `json:"data,omitempty"`
	Cursor *Cursor[T] `json:"cursor,omitempty"`
}

type Cursor[T any] struct {
	PageSize int    `json:"pageSize,omitempty"`
	HasMore  bool   `json:"hasMore"`
	Previous string `json:"previous,omitempty"`
	Next     string `json:"next,omitempty"`
	Data     []T    `json:"data"`
}

type ErrorResponse struct {
	ErrorCode    string `json:"errorCode,omitempty"`
	ErrorMessage string `json:"errorMessage,omitempty"`
	Details      string `json:"details,omitempty"`
}
