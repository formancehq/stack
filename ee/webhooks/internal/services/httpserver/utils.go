package httpserver

type MethodHTTP string


const (
	POST MethodHTTP = "POST"
	GET MethodHTTP = "GET"
	PUT MethodHTTP = "PUT"
	DELETE MethodHTTP = "DELETE"
) 

type Route struct {
	Method MethodHTTP
	Url string
}

func NewRoute(m MethodHTTP, u string ) Route {
	return Route{
		Method: m,
		Url:u,
	}
}