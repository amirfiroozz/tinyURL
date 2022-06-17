package routes

import (
	"net/http"
	"tiny-url/pkg/controllers"
	"tiny-url/pkg/middlewares"

	"github.com/gorilla/mux"
)

var userRoutes = func(router *mux.Router) {

	//TODO: check auth
	router.HandleFunc("/google/login", controllers.UserGoogleLogin).Methods(http.MethodGet)
	//TODO: make sure no body calls this endpoint
	router.HandleFunc("/google/callback", controllers.GoogleCallBack).Methods(http.MethodGet)
	router.HandleFunc("/create", controllers.CreateNewUser).Methods(http.MethodPost)
	router.HandleFunc("/show", middlewares.IfJWTLoggedIn(controllers.GetUserById)).Methods(http.MethodGet)
	router.HandleFunc("/showall", middlewares.IfJWTLoggedIn(controllers.GetAllUsers)).Methods(http.MethodGet)
	router.HandleFunc("/expire/{urlId}", middlewares.IfJWTLoggedIn(controllers.SetURLExpired)).Methods(http.MethodPatch)
	router.HandleFunc("/checkloggedin", middlewares.IfJWTLoggedIn(controllers.CheckLoggedIn)).Methods(http.MethodGet)
}
