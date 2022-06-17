package middlewares

import (
	"fmt"
	"net/http"
	"tiny-url/pkg/config"
	"tiny-url/pkg/models"
	"tiny-url/pkg/utils"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

type Error struct {
	Code int
	Msg  string
}

func IfLoggedIn(handlerFunc http.HandlerFunc) http.HandlerFunc {
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
		//TODO: why not accepting email??
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
