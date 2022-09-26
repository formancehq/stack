package api

import (
	"net"
	"net/http"
)

func StartServer(addr string, handler http.Handler) error {
	socket, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	go func() {
		err := http.Serve(socket, handler)
		if err != nil {
			panic(err)
		}
	}()
	return nil
}
