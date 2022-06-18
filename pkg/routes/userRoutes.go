package routes

import (
	"net/http"
	"tiny-url/pkg/controllers"
	auth "tiny-url/pkg/middlewares"

	"github.com/gorilla/mux"
)

var userRoutes = func(router *mux.Router) {

	router.HandleFunc("/google/login", controllers.UserGoogleLogin).Methods(http.MethodGet)
	router.HandleFunc("/google/callback", controllers.GoogleCallBack).Methods(http.MethodGet)
	router.HandleFunc("/create", auth.IsAdmin(controllers.CreateNewUser)).Methods(http.MethodPost)
	router.HandleFunc("/showall", auth.IsAdmin(controllers.GetAllUsers)).Methods(http.MethodGet)
	router.HandleFunc("/show", auth.IfJWTLoggedIn(controllers.GetUserById)).Methods(http.MethodGet)
	router.HandleFunc("/expire/{urlId}", auth.IfJWTLoggedIn(controllers.SetURLExpired)).Methods(http.MethodPatch)
	router.HandleFunc("/checkloggedin", auth.IfJWTLoggedIn(controllers.CheckLoggedIn)).Methods(http.MethodGet)
}
