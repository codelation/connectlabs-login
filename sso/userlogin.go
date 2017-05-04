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
	ul.AccessToken = user.AccessToken
	ul.AccessTokenSecret = user.AccessTokenSecret
	ul.AvatarURL = user.AvatarURL
	ul.Description = user.Description
	ul.Email = user.Email
	ul.ExpiresAt = user.ExpiresAt
	ul.FirstName = user.FirstName
	ul.LastName = user.LastName
	ul.Location = user.Location
	ul.Name = user.Name
	ul.NickName = user.NickName
	ul.Provider = user.Provider
	ul.RefreshToken = user.RefreshToken
}
