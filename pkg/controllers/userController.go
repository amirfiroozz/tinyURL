package controllers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"tiny-url/pkg/config"
	"tiny-url/pkg/models"
	"tiny-url/pkg/utils"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

func CreateNewUser(w http.ResponseWriter, r *http.Request) {
	var userBody models.User
	err := utils.ParseBody(r, &userBody)
	if err != nil {
		utils.SendError(w, r, *err)
		return
	}

	user, err := userBody.AddNewUserToDB()
	if err != nil {
		utils.SendError(w, r, *err)
		return
	}
	utils.SendResponse(w, r, user)

}
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := models.GetAllUsers()
	if err != nil {
		utils.SendError(w, r, *err)
	}
	utils.SendResponse(w, r, users)
}
func GetUserById(w http.ResponseWriter, r *http.Request) {
	//TODO: uncomment this
	// userId :=mux.Vars(r)["sessionUserId"]
	userId := mux.Vars(r)["userId"]
	user, err := models.GetUserById(userId)
	if err != nil {
		utils.SendError(w, r, *err)
		return
	}
	utils.SendResponse(w, r, user)
}

func UserGoogleLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	googleConfig := config.GetGoogleConfig()
	url := googleConfig.AuthCodeURL(config.GetConfigurationFile().State)
	http.Redirect(w, r, url, http.StatusSeeOther)

}
func GoogleCallBack(w http.ResponseWriter, r *http.Request) {
	type userEmailInfo struct {
		Email string `json:"email"`
	}
	state := r.URL.Query()["state"][0]
	if state != config.GetConfigurationFile().State {
		utils.SendError(w, r, generateUtilError(1, 500, "states are different"))
		return
	}
	code := r.URL.Query()["code"][0]
	googleConfig := config.GetGoogleConfig()
	token, err := googleConfig.Exchange(context.Background(), code)
	if err != nil {
		utilError := generateUtilError(1, 500, err.Error())
		utils.SendError(w, r, utilError)
		return
	}
	URL := config.GetConfigurationFile().Google.UserInfoAccessTokenURL
	resp, err := http.Get(URL + token.AccessToken)
	if err != nil {
		utilError := generateUtilError(2, 500, err.Error())
		utils.SendError(w, r, utilError)
		return
	}
	userData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		utilError := generateUtilError(3, 500, err.Error())
		utils.SendError(w, r, utilError)
		return
	}
	var userEmailInfoData userEmailInfo
	unmarshalingError := json.Unmarshal([]byte(userData), &userEmailInfoData)
	if unmarshalingError != nil {
		utilError := generateUtilError(4, 500, unmarshalingError.Error())
		utils.SendError(w, r, utilError)
		return
	}
	user, userError := models.CreateORFindUserByEmail(userEmailInfoData.Email)
	if userError != nil {
		utilError := generateUtilError(5, 500, userError.Msg)
		utils.SendError(w, r, utilError)
		return
	}
	sessionError := setSession(w, r, user.Email)
	if sessionError != nil {
		utilError := generateUtilError(6, 500, sessionError.Error())
		utils.SendError(w, r, utilError)
		return
	}
	utils.SendResponse(w, r, user)

}

func setSession(w http.ResponseWriter, r *http.Request, email string) error {
	sessionConf := config.GetConfigurationFile().Session
	var store = sessions.NewCookieStore([]byte(sessionConf.Secret))
	store.Options = &sessions.Options{
		Domain: sessionConf.Domain,
		Path:   sessionConf.Path,
		MaxAge: sessionConf.MaxAge,
	}
	session, _ := store.Get(r, "session")
	session.Values["email"] = email
	return session.Save(r, w)
}

func generateUtilError(code int, status int, msg string) utils.Error {
	return utils.Error{
		Code:   code,
		Status: status,
		Msg:    msg,
	}
}
