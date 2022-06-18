package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db            *gorm.DB
	configuration Configuration
	ctx           context.Context
	rdb           *redis.Client
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
	JWT struct {
		Secret string `json:"secret"`
		Exp    int    `json:"exp"`
	} `json:"jwt"`
	Token string `json:"token"`
	Redis struct {
		Addr     string `json:"addr"`
		Password string `json:"password"`
		DB       int    `json:"db"`
	} `json:"redis"`
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
	ctx = context.Background()
	//TODO: handle error
	setRedisClient()
}

func setRedisClient() {
	//TODO: handle error when redis is not up...
	rdb = redis.NewClient(&redis.Options{
		Addr:     configuration.Redis.Addr,
		Password: configuration.Redis.Password,
		DB:       configuration.Redis.DB,
	})
}
func GetRdb() *redis.Client {
	return rdb
}

func GetCtx() context.Context {
	return ctx
}

func GetDB() *gorm.DB {
	return db
}

func GetConfigurationFile() Configuration {
	return configuration
}
