package main

import (
	"encoding/json"
	"log"
	"os"
	"strconv"

	"github.com/ryanhatfield/connectlabs-login/ap"
	"github.com/ryanhatfield/connectlabs-login/app"
	"github.com/ryanhatfield/connectlabs-login/sso"
)

func main() {

	getEnv := func(key, def string) string {
		//helper method for getting environment variables with defaults
		if e := os.Getenv(key); e != "" {
			return e
		}
		return def
	}

	getEnvInt := func(key string, def int) int {
		e := getEnv(key, "0")
		i, err := strconv.ParseInt(e, 0, 32)
		if err != nil || i == 0 {
			return def
		}
		return int(i)
	}

	getEnvBool := func(key string, def bool) bool {
		e := getEnv(key, "false")
		i, err := strconv.ParseBool(e)
		if err != nil {
			return def
		}
		return i
	}

	debug := getEnvBool("CONNECTLABS_DEBUG", false)

	//sso.UserStorage and ap.SessionStorage can be their own interfaces, but all
	//are implemented in app.Data, so just cast it to the specific interface type
	database := &app.Data{
		DatabaseURL:      os.Getenv("DATABASE_URL"),
		MaxDBConnections: getEnvInt("MAX_DATABASE_CONNECTIONS", 20),
		Debug:            debug,
	}

	application := app.App{
		AccessPointHandler: &ap.AP{
			Secret:   getEnv("SECRET", "default"),
			Sessions: database,
		},
		Database: database,
		Debug:    debug,
		Port:     getEnv("PORT", "8080"),
		SingleSignOnHandler: &sso.SSO{
			Users:          database,
			KeyFacebook:    os.Getenv("FACEBOOK_KEY"),
			SecretFacebook: os.Getenv("FACEBOOK_SECRET"),
			KeyTwitter:     os.Getenv("TWITTER_KEY"),
			SecretTwitter:  os.Getenv("TWITTER_SECRET"),
			KeyGPlus:       os.Getenv("GPLUS_KEY"),
			SecretGPlus:    os.Getenv("GPLUS_SECRET"),
			CallbackURL:    os.Getenv("CALLBACK_URL"),
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
