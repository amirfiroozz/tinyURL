package models

import (
	"time"
	"tiny-url/pkg/utils"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type URL struct {
	ID           uuid.UUID `json:"id" gorm:"primary_key;type:uuid"`
	Created_At   time.Time `json:"created_at"`
	ClickedCount int       `json:"clickedCount"`
	ExpiresIn    int       `json:"expIn" validate:"numeric,gte=0"`
	Expired      bool      `json:"expired"`
	OriginalURl  string    `json:"originalURL" validate:"required,url"`
	ShortURL     string    `json:"shortURL" validate:"required,len=7" gorm:"type:varchar(7);unique"`
	UserId       uuid.UUID `json:"userId"`
}

func CreateNewURL(url URL) (*URL, *utils.Error) {
	tx := db.Create(&url)
	if tx.Error != nil {
		return nil, &utils.Error{
			Code:   1,
			Status: 500,
			Msg:    tx.Error.Error(),
		}
	}
	return &url, nil
}

func GetOriginalURLAndUpdateClickedCount(shortURL string) (*string, *utils.Error) {
	var url URL
	tx := db.Raw("select * from urls where short_url=?", shortURL).First(&url)
	if tx.Error != nil {
		return nil, &utils.Error{
			Code:   3,
			Status: 500,
			Msg:    tx.Error.Error(),
		}
	}
	if tx.RowsAffected == 0 {
		return nil, &utils.Error{
			Code:   1,
			Status: 404,
			Msg:    "nothing founded",
		}
	}
	EXPIRED := true
	if url.Expired == EXPIRED {
		return nil, &utils.Error{
			Code:   1,
			Status: 410,
			Msg:    "this url is expired!!!",
		}
	}
	err := updateClickedCount(&url)
	if err != nil {
		return nil, &utils.Error{
			Code:   3,
			Status: 500,
			Msg:    err.Error(),
		}
	}
	return &url.OriginalURl, nil

}

func SetURLExpired(urlId uuid.UUID, userId string) *utils.Error {
	var tx *gorm.DB
	if userId == "" {
		tx = db.Exec("UPDATE urls SET expired = true WHERE id = ?", urlId)
	} else {
		tx = db.Exec("UPDATE urls SET expired = true WHERE id = ? and user_id=?", urlId, userId)
	}
	if tx.Error != nil {
		return &utils.Error{
			Code:   1,
			Msg:    tx.Error.Error(),
			Status: 500,
		}
	}
	if tx.RowsAffected == 0 {
		return &utils.Error{
			Code:   2,
			Msg:    "404 not found!",
			Status: 404,
		}
	}
	return nil
}

func updateClickedCount(url *URL) error {
	url.ClickedCount++
	return db.Save(&url).Error
}
