package sso

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/markbates/goth"
)

type UserLogin struct {
	gorm.Model
	RawData           map[string]interface{} `gorm:"-"`
	Provider          string
	Email             string
	Name              string
	FirstName         string
	LastName          string
	NickName          string
	Description       string
	UserID            uint
	AvatarURL         string
	Location          string
	AccessToken       string
	AccessTokenSecret string
	RefreshToken      string
	ExpiresAt         time.Time
}

func (ul *UserLogin) Parse(user *goth.User) {
	//TODO: write parse code to save user
}
