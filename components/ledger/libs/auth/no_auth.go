package auth

import "net/http"

type noAuth struct{}

func (a noAuth) Authenticate(w http.ResponseWriter, r *http.Request) (bool, error) {
	return true, nil
}

func newNoAuth() *noAuth {
	return &noAuth{}
}
