package routes

import (
	"net/http"
	"tiny-url/pkg/controllers"
	"tiny-url/pkg/middlewares"

	"github.com/gorilla/mux"
)

var userRoutes = func(router *mux.Router) {

	router.HandleFunc("/google/login", controllers.UserGoogleLogin).Methods(http.MethodGet)
	//TODO: make sure no body calls this endpoint
	router.HandleFunc("/google/callback", controllers.GoogleCallBack).Methods(http.MethodGet)
	router.HandleFunc("/create", controllers.CreateNewUser).Methods(http.MethodPost)
	//TODO: check auth
	router.HandleFunc("/show/{userId}", controllers.GetUserById).Methods(http.MethodGet)
	router.HandleFunc("/showall", controllers.GetAllUsers).Methods(http.MethodGet)
	//TODO: it should be patch method
	router.HandleFunc("/expire/{urlId}", middlewares.IfLoggedIn(controllers.SetURLExpired)).Methods(http.MethodGet)
}
