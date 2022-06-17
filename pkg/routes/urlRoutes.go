package routes

import (
	"net/http"
	"tiny-url/pkg/controllers"
	"tiny-url/pkg/middlewares"

	"github.com/gorilla/mux"
)

var urlRoutes = func(router *mux.Router) {
	//TODO: remove {userId} parameter. it will be added by session controlling middleware
	router.HandleFunc("/create", middlewares.IfJWTLoggedIn(controllers.CreateNewURL)).Methods(http.MethodPost)
	router.HandleFunc("/show/{shortURL}", controllers.GetOriginalURLAndUpdateClickedCount).Methods(http.MethodGet)
}
