package controllers

import (
	"fmt"
	"net/http"
	"time"
	redis "tiny-url/pkg/middlewares"
	"tiny-url/pkg/models"
	"tiny-url/pkg/utils"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

func CreateNewURL(w http.ResponseWriter, r *http.Request) {
	urlBody := &models.URL{}
	setURLDefaults(r, urlBody)
	err := utils.ParseBody(r, urlBody)
	if err != nil {
		utils.SendError(w, r, *err)
		return
	}
	url, err := models.CreateNewURL(*urlBody)
	if err != nil {
		utils.SendError(w, r, *err)
		return
	}
	if url.ExpiresIn > 0 {
		fmt.Printf("%v:	URL:%v sent to scheduler to get expired....\n", time.Now(), url.ID)
		setURLExpired(*url)
	}
	err = redis.Set(url.ExpiresIn, url.ShortURL, url.OriginalURl)
	if err != nil {
		fmt.Println("redis memory set error")
	}

	utils.SendResponse(w, r, url)
}

func GetOriginalURLAndUpdateClickedCount(w http.ResponseWriter, r *http.Request) {
	type originalURL struct {
		URL string `json:"url"`
	}
	var originalUrl originalURL

	shortURL := mux.Vars(r)["shortURL"]
	if redis.IfExists(shortURL) {
		val, err := redis.Get(shortURL)
		if err != nil {
			utils.SendError(w, r, *err)
			return
		}
		originalUrl.URL = *val
		updateClickedCount(w, r, shortURL)
		utils.SendResponse(w, r, originalUrl)
		return
	}
	url, err := models.GetOriginalURLAndUpdateClickedCount(shortURL)
	if err != nil {
		utils.SendError(w, r, *err)
		return
	}
	originalUrl.URL = *url
	utils.SendResponse(w, r, originalUrl)

}

func SetURLExpired(w http.ResponseWriter, r *http.Request) {

	urlId := mux.Vars(r)["urlId"]
	userId := mux.Vars(r)["sessionUserId"]

	_, err := models.SetURLExpired(uuid.FromStringOrNil(urlId), userId)
	if err != nil {
		utils.SendError(w, r, *err)
		return
	}

	utils.SendResponse(w, r, &utils.Succes{
		Msg: "successfuly set expired",
	})

}

func CheckLoggedIn(w http.ResponseWriter, r *http.Request) {
	utils.SendResponse(w, r, &utils.Succes{
		Msg: "user is logged in",
	})
}

func setURLExpired(url models.URL) {
	go func() {
		time.Sleep(time.Duration(url.ExpiresIn) * time.Minute)
		_, err := models.SetURLExpired(url.ID, "")
		if err != nil {
			fmt.Printf("error ecurred: %v\n", err.Msg)
			return
		}
		fmt.Printf("%v:	URL:%v expired succesfully....\n", time.Now(), url.ID)
	}()
}

func setURLDefaults(r *http.Request, urlBody *models.URL) {
	urlBody.ID = uuid.NewV4()
	urlBody.Created_At = time.Now()
	urlBody.ClickedCount = 0
	urlBody.Expired = false
	urlBody.ShortURL = generateShortHashString()
	urlBody.UserId = uuid.FromStringOrNil(mux.Vars(r)["sessionUserId"])

}

func generateShortHashString() string {
	return fmt.Sprintf("%v", uuid.NewV4())[:7]
}

func updateClickedCount(w http.ResponseWriter, r *http.Request, shortURL string) {
	go func() {
		fmt.Printf("Increasing clicked count of shortURL:%v\n", shortURL)
		err := models.UpdateClickedCount(shortURL)
		if err != nil {
			utils.SendError(w, r, *err)
			fmt.Printf("error while increasing clicked count of shortURL:%v\n", shortURL)
			return
		}
	}()
}
