package sso

import (
	"github.com/jinzhu/gorm"
	"github.com/ryanhatfield/connectlabs-login/ap"
)

//UserStorage provides an interface for the data adapter to implement user srorage
type UserStorage interface {
	FindUserByID(user *User, ID uint) error
	AddLoginToUser(userID uint, login UserLogin) error
	FindUserByDevice(user *User, ip string, mac string) error
	FindUserID(session, ip, mac string) string
}

type User struct {
	gorm.Model
	NickName string
	Email    string
	Sessions []ap.Session `gorm:"ForeignKey:UserID"`
	Logins   []UserLogin  `gorm:"ForeignKey:UserID"`
}
