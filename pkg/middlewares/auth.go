package middlewares

import (
	"fmt"
	"net/http"
	"tiny-url/pkg/config"
	"tiny-url/pkg/models"
	"tiny-url/pkg/utils"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

type Error struct {
	Code int
	Msg  string
}

func IfSessionLoggedIn(handlerFunc http.HandlerFunc) http.HandlerFunc {
	SESSION_SECRET := config.GetConfigurationFile().Session.Secret
	var store = sessions.NewCookieStore([]byte(SESSION_SECRET))
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session")
		email, ok := session.Values["email"]
		if !ok {
			var resError utils.Error = utils.Error{
				Code:   1,
				Status: 403,
				Msg:    "Unauthorized!!!",
			}
			utils.SendError(w, r, resError)
			return
		}
		user, err := models.FindUserByEmail(fmt.Sprintf("%v", email))
		if err != nil {
			var resError utils.Error = utils.Error{
				Code:   err.Code,
				Status: err.Status,
				Msg:    err.Msg,
			}
			utils.SendError(w, r, resError)
			return
		}
		mux.Vars(r)["sessionUserId"] = user.ID.String()
		handlerFunc.ServeHTTP(w, r)
	}
}

func IfJWTLoggedIn(handlerFunc http.HandlerFunc) http.HandlerFunc {
	JWT := config.GetConfigurationFile().JWT
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] == nil {
			var resError utils.Error = utils.Error{
				Code:   1,
				Status: 403,
				Msg:    "Unauthorized!!!",
			}
			utils.SendError(w, r, resError)
			return
		}

		var mySigningKey = []byte(JWT.Secret)

		token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				str := "There was an error in parsing"
				return nil, fmt.Errorf("%v", str)
			}
			return mySigningKey, nil
		})

		if err != nil {
			var resError utils.Error = utils.Error{
				Code:   1,
				Status: 500,
				Msg:    err.Error(),
			}
			utils.SendError(w, r, resError)
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if ok && token.Valid {
			email := claims["email"]
			user, err := models.FindUserByEmail(fmt.Sprintf("%v", email))
			if err != nil {
				var resError utils.Error = utils.Error{
					Code:   err.Code,
					Status: err.Status,
					Msg:    err.Msg,
				}
				utils.SendError(w, r, resError)
				return
			}
			mux.Vars(r)["sessionUserId"] = user.ID.String()
			handlerFunc.ServeHTTP(w, r)
			return
		}
		var resError utils.Error = utils.Error{
			Code:   1,
			Status: 500,
			Msg:    "error",
		}
		utils.SendError(w, r, resError)

	}
}
