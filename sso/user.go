package sso

import (
	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
	"github.com/ryanhatfield/connectlabs-login/ap"
)

//UserStorage provides an interface for the data adapter to implement user srorage
type UserStorage interface {
	FindUserByID(user *User, ID uint) error
	AddLoginToUser(userID uint, login UserLogin) error
	FindUserByDevice(user *User, mac string, node string) error
	FindUserID(session, ip, mac string) string
	SessionStore() sessions.Store
	AddSessionToUser(token string, userID uint) error
}

type User struct {
	gorm.Model
	NickName string
	Email    string
	Sessions []ap.Session `gorm:"ForeignKey:UserID"`
	Logins   []UserLogin  `gorm:"ForeignKey:UserID"`
}
