package main

import (
	"encoding/json"
	"log"
	"os"
	"strconv"

	"github.com/ryanhatfield/connectlabs-login/ap"
	"github.com/ryanhatfield/connectlabs-login/api"
	"github.com/ryanhatfield/connectlabs-login/app"
	"github.com/ryanhatfield/connectlabs-login/server"
	"github.com/ryanhatfield/connectlabs-login/sso"
)

func main() {

	getEnv := func(key string) string {
		return os.Getenv(key)
	}

	getEnvD := func(key, def string) string {
		//helper method for getting environment variables with defaults
		if e := getEnv(key); e != "" {
			return e
		}
		return def
	}

	getEnvInt := func(key string, def int) int {
		e := getEnvD(key, "0")
		i, err := strconv.ParseInt(e, 0, 32)
		if err != nil || i == 0 {
			return def
		}
		return int(i)
	}

	getEnvBool := func(key string, def bool) bool {
		e := getEnvD(key, "false")
		i, err := strconv.ParseBool(e)
		if err != nil {
			return def
		}
		return i
	}

	debug := getEnvBool("CONNECTLABS_DEBUG", false)

	//sso.UserStorage and ap.SessionStorage can be their own interfaces, but all
	//are implemented in app.Data, so just cast it to the specific interface type
	database := &server.Data{
		DatabaseURL:      getEnv("DATABASE_URL"),
		Debug:            debug,
		MaxDBConnections: getEnvInt("MAX_DATABASE_CONNECTIONS", 20),
	}

	application := app.App{
		AccessPointHandler: &ap.AP{
			Secret:   getEnvD("SECRET", "default"),
			Sessions: database,
		},
		API: &api.API{
			AuthorizeToken: getEnv("API_AUTHORIZATION_TOKEN"),
			Users:          database,
		},
		Database: database,
		Debug:    debug,
		Port:     getEnvD("PORT", "8080"),
		SingleSignOnHandler: &sso.SSO{
			Users: database,
			Sites: []sso.SiteConfig{
				sso.SiteConfig{
					Name:      "Default",
					Title:     "Connect Labs WiFi",
					SubTitle:  "Log In for Specials",
					Providers: []string{"email"},
				},
				sso.SiteConfig{
					Name:      "BeardFromFargo Guest",
					Title:     "Beard Wifi",
					SubTitle:  "Log In for Specials",
					Providers: []string{"facebook", "gplus", "twitter", "email"},
				},
			},
			KeyFacebook:    getEnv("FACEBOOK_KEY"),
			SecretFacebook: getEnv("FACEBOOK_SECRET"),
			KeyTwitter:     getEnv("TWITTER_KEY"),
			SecretTwitter:  getEnv("TWITTER_SECRET"),
			KeyGPlus:       getEnv("GPLUS_KEY"),
			SecretGPlus:    getEnv("GPLUS_SECRET"),
			CallbackURL:    getEnv("CALLBACK_URL"),
		},
	}

	if application.Debug {
		log.Println("Debug logging is enabled")
		j, _ := json.MarshalIndent(os.Environ(), "", "  ")
		log.Printf("Environment:\n%s", j)

		j, _ = json.MarshalIndent(application, "", "  ")
		log.Printf("Application:\n%s", j)
	}

	err := application.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}
