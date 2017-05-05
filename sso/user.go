package sso

import (
	"github.com/jinzhu/gorm"
	"github.com/ryanhatfield/connectlabs-login/ap"
)

//UserStorage provides an interface for the data adapter to implement user srorage
type UserStorage interface {
	FindByID(id uint, out interface{}) error
	FindUserByDevice(mac string, node string, out *User) error
	FindUserIDByDevice(token, mac, node string) string
	AddLoginToUser(id uint, login UserLogin) error
}

//User holds information about a single wifi user
type User struct {
	gorm.Model
	NickName string
	Email    string
	Sessions []ap.Session `gorm:"ForeignKey:UserID"`
	Logins   []UserLogin  `gorm:"ForeignKey:UserID"`
}
