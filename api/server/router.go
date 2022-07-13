package server

import "github.com/gorilla/mux"

func newRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc(healthCheckPath, healthCheckHandler)
	return router
}
