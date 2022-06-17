package routes

import (
	"net/http"
	"tiny-url/pkg/controllers"
	"tiny-url/pkg/middlewares"

	"github.com/gorilla/mux"
)

var userRoutes = func(router *mux.Router) {

	router.HandleFunc("/google/login", controllers.UserGoogleLogin).Methods(http.MethodGet)
	router.HandleFunc("/google/callback", controllers.GoogleCallBack).Methods(http.MethodGet)
	router.HandleFunc("/create", middlewares.IsAdmin(controllers.CreateNewUser)).Methods(http.MethodPost)
	router.HandleFunc("/showall", middlewares.IsAdmin(controllers.GetAllUsers)).Methods(http.MethodGet)
	router.HandleFunc("/show", middlewares.IfJWTLoggedIn(controllers.GetUserById)).Methods(http.MethodGet)
	router.HandleFunc("/expire/{urlId}", middlewares.IfJWTLoggedIn(controllers.SetURLExpired)).Methods(http.MethodPatch)
	router.HandleFunc("/checkloggedin", middlewares.IfJWTLoggedIn(controllers.CheckLoggedIn)).Methods(http.MethodGet)
}
