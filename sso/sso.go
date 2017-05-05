package sso

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/gplus"
	"github.com/markbates/goth/providers/twitter"
	"github.com/ryanhatfield/connectlabs-login/sso/oauth"
)

type SSO struct {
	Users          UserStorage
	Sites          []SiteConfig
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

	user, err := oauth.CompleteUserAuth(res, req, "")
	if err != nil {
		log.Printf("error getting user from oauth provider: %s", err.Error())
		return
	}

	dbUser := &User{}

	sessionNode, _ := oauth.GetFromSession("node", req)
	sessionMac, _ := oauth.GetFromSession("mac", req)

	if err = sso.Users.FindUserByDevice(sessionMac, sessionNode, dbUser); err != nil {
		log.Printf("error getting user by device mac / node: %s\n", err.Error())
	}

	if err = sso.Users.AddLoginToUser(dbUser.ID, UserLogin{
		AccessToken:       user.AccessToken,
		AccessTokenSecret: user.AccessTokenSecret,
		AvatarURL:         user.AvatarURL,
		Description:       user.Description,
		Email:             user.Email,
		ExpiresAt:         user.ExpiresAt,
		FirstName:         user.FirstName,
		LastName:          user.LastName,
		Location:          user.Location,
		Name:              user.Name,
		NickName:          user.NickName,
		Provider:          user.Provider,
		RefreshToken:      user.RefreshToken,
	}); err != nil {
		log.Printf("error getting user by device mac / node: %s\n", err.Error())
	}

	userurl, _ := oauth.GetFromSession("userurl", req)
	res.Header().Set("Location", userurl)
	res.WriteHeader(http.StatusTemporaryRedirect)

	log.Println("finished auth callback")
}

func (sso *SSO) HandleAuthLogout(res http.ResponseWriter, req *http.Request) {
	oauth.Logout(res, req)
	userurl, _ := oauth.GetFromSession("userurl", req)
	res.Header().Set("Location", userurl)
	res.WriteHeader(http.StatusTemporaryRedirect)
}

func (sso *SSO) HandleAuthLogin(res http.ResponseWriter, req *http.Request) {
	// try to get the user without re-authenticating
	if gothUser, err := oauth.CompleteUserAuth(res, req, ""); err == nil {
		j, _ := json.MarshalIndent(gothUser, "", "  ")
		res.Write(j)
	} else {
		oauth.BeginAuthHandler(res, req)
	}
}

func (sso *SSO) HandleLoginPage(res http.ResponseWriter, req *http.Request) {
	//get variables from URL and save in session
	req.ParseForm()
	node := strings.Replace(req.Form.Get("called"), "-", ":", -1)
	mac := strings.Replace(req.Form.Get("mac"), "-", ":", -1)
	ssid := req.Form.Get("ssid")
	userurl := req.Form.Get("userurl")
	uamip := req.Form.Get("uamip")
	uamport := "8085"
	sessionToken := "default"

	oauth.StoreInSession("node", node, req, res)
	oauth.StoreInSession("mac", mac, req, res)
	oauth.StoreInSession("ssid", ssid, req, res)
	oauth.StoreInSession("userurl", userurl, req, res)
	oauth.StoreInSession("uamip", userurl, req, res)
	oauth.StoreInSession("uamport", userurl, req, res)

	log.Println("handling login page")

	t, err := template.ParseFiles("www/login.gohtml")
	if err != nil {
		log.Println(err)
	}

	getConf := func(ssid string) *SiteConfig {
		for _, c := range sso.Sites {
			if c.Name == ssid {
				return &c
			}
		}
		return nil
	}

	config := getConf(ssid)
	if config == nil {
		config = getConf("Default")
	}
	if config == nil {
		config = &sso.Sites[0]
	}

	dbUser := &User{}

	if err = sso.Users.FindUserByDevice(mac, node, dbUser); err != nil {
		log.Printf("error getting user by device mac / node: %s\n", err.Error())
	}

	view := &LoginView{
		Title:            config.Title,
		SubTitle:         config.SubTitle,
		SSID:             config.Name,
		TwitterProvider:  true,
		FacebookProvider: true,
		GPlusProvider:    true,
		Email:            dbUser.Email,
		UserAccountManagementURL:     fmt.Sprintf("http://%s:%s/cgi-bin/submit.cgi", uamip, uamport),
		UserAccountManagementIP:      uamip,
		UserAccountManagementSession: sessionToken,
		MacAddress:                   mac,
	}

	if user, err := oauth.CompleteUserAuth(res, req, "facebook"); err == nil {
		view.FacebookUser = user
		if view.Email == "" {
			view.Email = user.Email
		}
	}
	if user, err := oauth.CompleteUserAuth(res, req, "gplus"); err == nil {
		view.GPlusUser = user
		if view.Email == "" {
			view.Email = user.Email
		}
	}
	if user, err := oauth.CompleteUserAuth(res, req, "twitter"); err == nil {
		view.TwitterUser = user
		if view.Email == "" {
			view.Email = user.Email
		}
	}

	//nice for debugging
	// m, _ := json.MarshalIndent(*view, "", "  ")
	// view.Message = string(m)

	t.Execute(res, *view)
}
