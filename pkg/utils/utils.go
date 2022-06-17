package utils

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/go-playground/validator"
)

type Error struct {
	Code   int    `json:"code"`
	Status int    `json:"status"`
	Msg    string `json:"msg"`
}
type Succes struct {
	Msg string `json:"msg"`
}

func ParseBody(r *http.Request, x interface{}) *Error {
	body, err := ioutil.ReadAll(r.Body)
	if err == nil {
		err := json.Unmarshal([]byte(body), x)
		if err != nil {
			return &Error{
				Code:   1,
				Status: 400,
				Msg:    err.Error(),
			}
		}
		invalidDataError := checkValidData(x)
		if invalidDataError != nil {
			return invalidDataError
		}
	}
	return nil
}

func checkValidData(x interface{}) *Error {
	validate := validator.New()
	err := validate.Struct(x)
	if err != nil {
		return &Error{
			Code:   1,
			Status: 400,
			Msg:    err.Error(),
		}
	}
	return nil
}

func SendError(w http.ResponseWriter, r *http.Request, err Error) {
	res, _ := json.Marshal(err)
	w.Header().Set("Content-Type", "Application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(err.Status)
	w.Write(res)
}

func SendResponse(w http.ResponseWriter, r *http.Request, x interface{}) {
	res, _ := json.Marshal(x)
	w.Header().Set("Content-Type", "Application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
