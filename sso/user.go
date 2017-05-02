package sso

import (
	"github.com/jinzhu/gorm"
	"github.com/ryanhatfield/connectlabs-login/ap"
)

//UserStorage provides an interface for the data adapter to implement user srorage
type UserStorage interface {
	FindUser(user *User, userID uint) error
	AddLoginToUser(userID uint, login UserLogin) error
}

type User struct {
	gorm.Model
	ID       uint
	NickName string
	Email    string
	Sessions []ap.Session `gorm:"ForeignKey:UserID"`
	Logins   []UserLogin  `gorm:"ForeignKey:UserID"`
}
