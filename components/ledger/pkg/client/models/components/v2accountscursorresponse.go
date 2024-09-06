// Code generated by Speakeasy (https://speakeasyapi.com). DO NOT EDIT.

package components

type V2AccountsCursorResponseCursor struct {
	PageSize int64       `json:"pageSize"`
	HasMore  bool        `json:"hasMore"`
	Previous *string     `json:"previous,omitempty"`
	Next     *string     `json:"next,omitempty"`
	Data     []V2Account `json:"data"`
}

func (o *V2AccountsCursorResponseCursor) GetPageSize() int64 {
	if o == nil {
		return 0
	}
	return o.PageSize
}

func (o *V2AccountsCursorResponseCursor) GetHasMore() bool {
	if o == nil {
		return false
	}
	return o.HasMore
}

func (o *V2AccountsCursorResponseCursor) GetPrevious() *string {
	if o == nil {
		return nil
	}
	return o.Previous
}

func (o *V2AccountsCursorResponseCursor) GetNext() *string {
	if o == nil {
		return nil
	}
	return o.Next
}

func (o *V2AccountsCursorResponseCursor) GetData() []V2Account {
	if o == nil {
		return []V2Account{}
	}
	return o.Data
}

type V2AccountsCursorResponse struct {
	Cursor V2AccountsCursorResponseCursor `json:"cursor"`
}

func (o *V2AccountsCursorResponse) GetCursor() V2AccountsCursorResponseCursor {
	if o == nil {
		return V2AccountsCursorResponseCursor{}
	}
	return o.Cursor
}
