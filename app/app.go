package app

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/pat"
	"github.com/gorilla/sessions"
	"github.com/ryanhatfield/connectlabs-login/ap"
	"github.com/ryanhatfield/connectlabs-login/sso"
)

const SessionKey = "connectlabs-login"

type App struct {
	AccessPointHandler  *ap.AP
	Database            *Data
	Debug               bool
	Port                string
	SingleSignOnHandler *sso.SSO
	router              *pat.Router
	SessionStore        sessions.Store
	initialized         bool
}

func (a *App) ListenAndServe() error {
	a.Initialize()
	return http.ListenAndServe(":"+a.Port, a)
}

func (a *App) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	log.Printf("Serving app request: %s\n", req.URL.String())

	// cookieSession, err := a.SessionStore.Get(req, SessionKey)
	// if err != nil {
	// 	http.Error(res, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	s := time.Now()
	a.router.ServeHTTP(res, req)
	log.Printf("Finished serving app request: %s Request took %d nanoseconds", req.URL.String(), time.Now().Sub(s).Nanoseconds())
}

func (a *App) Initialize() error {
	if a.initialized {
		return nil
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
}
