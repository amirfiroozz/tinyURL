package config

import (
	"encoding/json"
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db            *gorm.DB
	configuration Configuration
)

type Configuration struct {
	DB struct {
		Host     string `json:"host"`
		User     string `json:"user"`
		Password string `json:"password"`
		DBName   string `json:"dbname"`
		Port     string `json:"port"`
	} `json:"db"`
	Session struct {
		Secret string `json:"secret"`
		Domain string `json:"domain"`
		Path   string `json:"path"`
		MaxAge int    `json:"maxAge"`
	} `json:"session"`
	State  string `json:"state"`
	Google struct {
		ClientID               string   `json:"clientID"`
		ClientSecret           string   `json:"clientSecret"`
		RedirectURL            string   `json:"redirectURL"`
		Scopes                 []string `json:"scopes"`
		UserInfoAccessTokenURL string   `json:"userInfoAccessTokenURL"`
	} `json:"google"`
}

func genarateConfigFile() {
	file, _ := os.Open("./pkg/config/conf.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&configuration)
	if err != nil {
		errorMsg := fmt.Sprintf("config file error: %v", err)
		panic(errorMsg)
	}
}

func init() {
	genarateConfigFile()
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable", configuration.DB.Host, configuration.DB.User, configuration.DB.Password, configuration.DB.DBName, configuration.DB.Port)
	d, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db = d
}
func GetDB() *gorm.DB {
	return db
}

func GetConfigurationFile() Configuration {
	return configuration
}
