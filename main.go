package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/pat"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/gplus"
	"github.com/markbates/goth/providers/twitter"
	"github.com/ryanhatfield/connectlabs-login/ap"
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

	maxDBCon, _ := strconv.ParseInt(getEnv("MAX_DATABASE_CONNECTIONS", "20"), 0, 32)

	dat := &data{
		DatabaseURL:      getEnv("DATABASE_URL", ""),
		MaxDBConnections: int(maxDBCon),
	}

	if err := dat.InitializeDB(); err != nil {
		log.Printf("error initializing db:\n%+v", err)
	}

	router := pat.New()
	sso := sso.SSO{
		Users: dat,
	}

	accessPoints := ap.AP{
		Secret:   getEnv("SECRET", "default"),
		Sessions: ap.SessionStorage(dat),
	}

	callbackURL := getEnv("CALLBACK_URL", "http://localhost:8080/auth/")

	goth.UseProviders(
		facebook.New(os.Getenv("FACEBOOK_KEY"), os.Getenv("FACEBOOK_SECRET"), callbackURL+"facebook/callback"),
		twitter.NewAuthenticate(os.Getenv("TWITTER_KEY"), os.Getenv("TWITTER_SECRET"), callbackURL+"twitter/callback"),
		gplus.New(os.Getenv("GPLUS_KEY"), os.Getenv("GPLUS_SECRET"), callbackURL+"gplus/callback"),
	)

	accessPoints.AddRoutes(router)
	sso.AddRoutes(router)

	log.Println(http.ListenAndServe(":"+getEnv("PORT", "8080"), router))
}
