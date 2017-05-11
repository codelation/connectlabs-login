package app

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/pat"
	"github.com/ryanhatfield/connectlabs-login/ap"
	"github.com/ryanhatfield/connectlabs-login/api"
	"github.com/ryanhatfield/connectlabs-login/server"
	"github.com/ryanhatfield/connectlabs-login/sso"
)

//SessionKey is used as the cookie storage name
const SessionKey = "connectlabs-login"

//App holds all the applicatin logic, and the entry point for the server
type App struct {
	AccessPointHandler  *ap.AP
	API                 *api.API
	Database            *server.Data
	Debug               bool
	Port                string
	SingleSignOnHandler *sso.SSO
	router              *pat.Router
	initialized         bool
}

//ListenAndServe initializes the server and calls ServeHTTP
func (a *App) ListenAndServe() error {
	a.Initialize()
	return http.ListenAndServe(":"+a.Port, a)
}

func (a *App) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	log.Printf("Serving app request: %s\n", req.URL.String())

	s := time.Now()

	a.router.ServeHTTP(res, req)

	log.Printf("Finished serving app request: %s Request took %d nanoseconds", req.URL.String(), time.Now().Sub(s).Nanoseconds())
}

//Initialize sets required settings and dependencies
func (a *App) Initialize() error {
	if a.initialized {
		return nil
	}

	if a.API.AuthorizeToken == "" {
		a.API.AuthorizeToken = "logmein"
	}

	a.router = pat.New()

	a.Database.Debug = a.Debug

	if err := a.Database.InitializeDB(); err != nil {
		return fmt.Errorf("error initializing db:\n%+v", err)
	}

	a.setRoutes()
	a.SingleSignOnHandler.Initialize()

	a.initialized = true

	return nil
}

func (a *App) setRoutes() {
	a.router.Get("/auth/{provider}/callback", a.SingleSignOnHandler.HandleAuthCallback)
	a.router.Get("/auth/logout/{provider}", a.SingleSignOnHandler.HandleAuthLogout)
	a.router.Get("/auth/{provider}/login", a.SingleSignOnHandler.HandleAuthLogin)
	a.router.Get("/auth/login.html", a.SingleSignOnHandler.HandleLoginPage)
	a.router.Get("/ap/auth.html", a.AccessPointHandler.HandleAPRequest)
	a.router.Get("/api/users/find/{mac}", a.API.AuthorizeAPI(a.API.GetUserByMAC))
	a.router.Get("/api/users/{id}", a.API.AuthorizeAPI(a.API.GetUserByID))
}
