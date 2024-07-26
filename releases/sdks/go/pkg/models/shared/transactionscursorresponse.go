// Code generated by Speakeasy (https://speakeasyapi.com). DO NOT EDIT.

package shared

type TransactionsCursorResponseCursor struct {
	Data     []Transaction `json:"data"`
	HasMore  bool          `json:"hasMore"`
	Next     *string       `json:"next,omitempty"`
	PageSize int64         `json:"pageSize"`
	Previous *string       `json:"previous,omitempty"`
}

func (o *TransactionsCursorResponseCursor) GetData() []Transaction {
	if o == nil {
		return []Transaction{}
	}
	return o.Data
}

func (o *TransactionsCursorResponseCursor) GetHasMore() bool {
	if o == nil {
		return false
	}
	return o.HasMore
}

func (o *TransactionsCursorResponseCursor) GetNext() *string {
	if o == nil {
		return nil
	}
	return o.Next
}

func (o *TransactionsCursorResponseCursor) GetPageSize() int64 {
	if o == nil {
		return 0
	}
	return o.PageSize
}

func (o *TransactionsCursorResponseCursor) GetPrevious() *string {
	if o == nil {
		return nil
	}
	return o.Previous
}

type TransactionsCursorResponse struct {
	Cursor TransactionsCursorResponseCursor `json:"cursor"`
}

func (o *TransactionsCursorResponse) GetCursor() TransactionsCursorResponseCursor {
	if o == nil {
		return TransactionsCursorResponseCursor{}
	}
	return o.Cursor
}
