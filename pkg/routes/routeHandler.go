package routes

import "github.com/gorilla/mux"

var CreateRoutes = func(router *mux.Router) {
	userRoutes(router.PathPrefix("/user").Subrouter())
	urlRoutes(router.PathPrefix("/url").Subrouter())
}
