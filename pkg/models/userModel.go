package models

import (
	"time"
	"tiny-url/pkg/utils"

	uuid "github.com/satori/go.uuid"
)

type User struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primary_key;"`
	Email      string    `json:"email" gorm:"unique;not null;" validate:"required,email"`
	Created_At time.Time `json:"created_at"`
	URLs       []URL     `json:"urls" gorm:"foreignKey:UserId"`
}

func (user *User) AddNewUserToDB() (*User, *utils.Error) {
	user.ID = uuid.NewV4()
	user.Created_At = time.Now()
	tx := db.Create(&user)
	if tx.Error != nil {
		return nil, &utils.Error{
			Code:   1,
			Status: 500,
			Msg:    tx.Error.Error(),
		}
	}
	return user, nil
}

func GetAllUsers() ([]User, *utils.Error) {
	var users []User
	tx := db.Preload("URLs").Find(&users)

	if tx.Error != nil {
		return nil, &utils.Error{
			Code:   1,
			Status: 500,
			Msg:    tx.Error.Error(),
		}
	}
	if tx.RowsAffected == 0 {
		return []User{}, nil
	}
	return users, nil

}
func GetUserById(userId string) (*User, *utils.Error) {
	var user User
	tx := db.Where("id=?", userId).Preload("URLs").Find(&user)
	if tx.Error != nil {
		return nil, &utils.Error{
			Code:   1,
			Status: 500,
			Msg:    tx.Error.Error(),
		}
	}
	return &user, nil
}

func CreateORFindUserByEmail(email string) (*User, *utils.Error) {
	var user User
	tx := db.Raw("select * from users where email=?", email).Scan(&user)
	if tx.Error != nil {
		return nil, &utils.Error{
			Code:   1,
			Status: 500,
			Msg:    tx.Error.Error(),
		}
	}
	if tx.RowsAffected == 0 {
		user.Email = email
		newUser, _ := user.AddNewUserToDB()
		return newUser, nil
	}
	return &user, nil
}
func FindUserByEmail(email string) (*User, *utils.Error) {
	var user User
	tx := db.Raw("select * from users where email=?", email).Scan(&user)
	if tx.Error != nil {
		return nil, &utils.Error{
			Code:   1,
			Status: 500,
			Msg:    tx.Error.Error(),
		}
	}
	if tx.RowsAffected == 0 {
		return nil, &utils.Error{
			Code:   1,
			Status: 404,
			Msg:    "no user exists with this email",
		}
	}
	return &user, nil
}
