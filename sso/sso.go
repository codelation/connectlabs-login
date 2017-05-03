package sso

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/gplus"
	"github.com/markbates/goth/providers/twitter"
	"github.com/ryanhatfield/connectlabs-login/sso/oauth"
)

type SSO struct {
	Users          UserStorage
	SessionStore   string
	KeyFacebook    string
	SecretFacebook string
	KeyGPlus       string
	SecretGPlus    string
	KeyTwitter     string
	SecretTwitter  string
	CallbackURL    string
	initialized    bool
}

const SessionKey = "connectlabs-login"

type SiteConfigStorage interface {
	GetSiteConfig(site string) interface{}
}

func (sso *SSO) Initialize() {
	if sso.initialized {
		return
	}
	goth.UseProviders(
		facebook.New(sso.KeyFacebook, sso.SecretFacebook, sso.CallbackURL+"facebook/callback"),
		twitter.NewAuthenticate(sso.KeyTwitter, sso.SecretTwitter, sso.CallbackURL+"twitter/callback"),
		gplus.New(sso.KeyGPlus, sso.SecretGPlus, sso.CallbackURL+"gplus/callback"),
	)
	sso.initialized = true
}

func (sso *SSO) HandleAuthCallback(res http.ResponseWriter, req *http.Request) {
	log.Println("starting auth callback")
	store := sso.Users.SessionStore()

	session, sessionerr := store.Get(req, SessionKey)
	if sessionerr != nil {
		log.Printf("error getting session store: %s", sessionerr.Error())
	}

	user, err := oauth.CompleteUserAuth(res, req)
	if err != nil {
		log.Printf("error getting user from oauth provider: %s", err.Error())
		return
	}

	dbUser := &User{}

	sessionNode, _ := session.Values["node"].(string)
	sessionMac, _ := session.Values["mac"].(string)

	if err = sso.Users.FindUserByDevice(dbUser, sessionMac, sessionNode); err != nil {
		log.Printf("error getting user by device mac / node: %s\n", err.Error())
	} else {
		userjs, _ := json.MarshalIndent(dbUser, "", "  ")
		res.Write(userjs)
	}

	j, _ := json.MarshalIndent(user, "", "  ")
	res.Write(j)

	log.Println("finished auth callback")
}

func (sso *SSO) HandleAuthLogout(res http.ResponseWriter, req *http.Request) {
	oauth.Logout(res, req)
	res.Header().Set("Location", "/")
	res.WriteHeader(http.StatusTemporaryRedirect)
}

func (sso *SSO) HandleAuthLogin(res http.ResponseWriter, req *http.Request) {
	// try to get the user without re-authenticating
	if gothUser, err := oauth.CompleteUserAuth(res, req); err == nil {
		j, _ := json.MarshalIndent(gothUser, "", "  ")
		res.Write(j)
	} else {
		oauth.BeginAuthHandler(res, req)
	}
}

func (sso *SSO) HandleLoginPage(res http.ResponseWriter, req *http.Request) {

	log.Println("handling login page")

	t, err := template.ParseFiles("www/login.gohtml")
	if err != nil {
		log.Println(err)
	}
	t.Execute(res,
		SiteConfig{
			Name:      "BeardFromFargo",
			Title:     "Beard Wifi",
			SubTitle:  "Log In for Specials",
			Providers: []string{"facebook", "gplus", "twitter", "email"},
		})
}
