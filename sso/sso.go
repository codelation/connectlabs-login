package sso

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/pat"
	"github.com/ryanhatfield/connectlabs-login/sso/util"
)

type SSO struct {
	Providers []Provider
	Users     UserStorage
}

type SiteConfigStorage interface {
	GetSiteConfig(site string) interface{}
}

func (sso *SSO) AddRoutes(p *pat.Router) {
	p.Get("/auth/{provider}/callback", func(res http.ResponseWriter, req *http.Request) {

		p := req.URL.Query().Get(":provider")
		log.Printf("handling /auth/%s/callback\n", p)

		user, err := util.CompleteUserAuth(res, req)
		if err != nil {
			fmt.Fprintln(res, err)
			return
		}
		j, _ := json.MarshalIndent(user, "", "  ")
		res.Write(j)
	})

	p.Get("/auth/logout/{provider}", func(res http.ResponseWriter, req *http.Request) {
		util.Logout(res, req)
		res.Header().Set("Location", "/")
		res.WriteHeader(http.StatusTemporaryRedirect)
	})

	p.Get("/auth/{provider}/login", func(res http.ResponseWriter, req *http.Request) {
		// try to get the user without re-authenticating
		if gothUser, err := util.CompleteUserAuth(res, req); err == nil {
			j, _ := json.MarshalIndent(gothUser, "", "  ")
			res.Write(j)
		} else {
			util.BeginAuthHandler(res, req)
		}
	})

	p.Get("/auth/login.html", sso.handleLoginPage)

}

func (sso *SSO) handleLoginPage(w http.ResponseWriter, r *http.Request) {

	log.Println("handling login page")
	t, err := template.ParseFiles("www/login.gohtml")
	if err != nil {
		log.Println(err)
	}
	t.Execute(w,
		SiteConfig{
			Name:      "BeardFromFargo",
			Title:     "Beard Wifi",
			SubTitle:  "Log In for Specials",
			Providers: []string{"facebook", "gplus", "twitter", "email"},
		})
}
