package sso

import (
	"encoding/json"
	"fmt"
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
	KeyFacebook    string
	SecretFacebook string
	KeyGPlus       string
	SecretGPlus    string
	KeyTwitter     string
	SecretTwitter  string
	CallbackURL    string
	initialized    bool
}

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

	user, err := oauth.CompleteUserAuth(res, req)
	if err != nil {
		fmt.Fprintln(res, err)
		return
	}
	dbUser := &User{}
	j, _ := json.MarshalIndent(user, "", "  ")
	res.Write(j)
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
