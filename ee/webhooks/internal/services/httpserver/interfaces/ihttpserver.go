package interfaces 

import (
	"context"
	"net/http"
)

type IHTTPServer interface {
	
	Register(method string, url string, handler func(http.ResponseWriter, *http.Request))
	Run(context.Context) error
	Stop(context.Context) error
}