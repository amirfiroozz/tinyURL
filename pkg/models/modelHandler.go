package models

import (
	"tiny-url/pkg/config"

	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	db = config.GetDB()
	db.AutoMigrate(&User{}, &URL{})
}
