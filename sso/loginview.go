package sso

import "github.com/markbates/goth"

type LoginView struct {
	TwitterProvider              bool
	TwitterUser                  goth.User
	FacebookProvider             bool
	FacebookUser                 goth.User
	GPlusProvider                bool
	GPlusUser                    goth.User
	SSID                         string
	Title                        string
	SubTitle                     string
	Email                        string
	Message                      string
	BackgroundURL                string
	UserAccountManagementIP      string
	UserAccountManagementPort    string
	UserAccountManagementURL     string
	UserAccountManagementSession string
	MacAddress                   string
}
