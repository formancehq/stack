// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package shared

type PaymentsCursorCursor struct {
	Data     []Payment `json:"data"`
	HasMore  bool      `json:"hasMore"`
	Next     *string   `json:"next,omitempty"`
	PageSize int64     `json:"pageSize"`
	Previous *string   `json:"previous,omitempty"`
}

func (o *PaymentsCursorCursor) GetData() []Payment {
	if o == nil {
		return []Payment{}
	}
	return o.Data
}

func (o *PaymentsCursorCursor) GetHasMore() bool {
	if o == nil {
		return false
	}
	return o.HasMore
}

func (o *PaymentsCursorCursor) GetNext() *string {
	if o == nil {
		return nil
	}
	return o.Next
}

func (o *PaymentsCursorCursor) GetPageSize() int64 {
	if o == nil {
		return 0
	}
	return o.PageSize
}

func (o *PaymentsCursorCursor) GetPrevious() *string {
	if o == nil {
		return nil
	}
	return o.Previous
}

type PaymentsCursor struct {
	Cursor PaymentsCursorCursor `json:"cursor"`
}

func (o *PaymentsCursor) GetCursor() PaymentsCursorCursor {
	if o == nil {
		return PaymentsCursorCursor{}
	}
	return o.Cursor
}
