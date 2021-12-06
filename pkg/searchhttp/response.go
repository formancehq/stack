package searchhttp

type Response struct {
	Cursor *Page       `json:"cursor,omitempty"`
	Data   interface{} `json:"data,omitempty"`
}
