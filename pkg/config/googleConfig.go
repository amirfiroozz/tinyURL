package config

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func GetGoogleConfig() *oauth2.Config {
	googleConfigFile := GetConfigurationFile().Google
	googleConfig := &oauth2.Config{
		ClientID:     googleConfigFile.ClientID,
		ClientSecret: googleConfigFile.ClientSecret,
		RedirectURL:  googleConfigFile.RedirectURL,
		Scopes:       googleConfigFile.Scopes,
		Endpoint:     google.Endpoint,
	}
	return googleConfig
}
