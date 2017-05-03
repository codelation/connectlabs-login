package sso

import (
	"context"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/ryanhatfield/connectlabs-login/ap"
)

type key int

const userIDKey key = 0

//UserStorage provides an interface for the data adapter to implement user srorage
type UserStorage interface {
	FindUser(user *User, ID uint) error
	AddLoginToUser(userID uint, login UserLogin) error
	FindUserID(session, ip, mac string) string
}

type User struct {
	gorm.Model
	NickName string
	Email    string
	Sessions []ap.Session `gorm:"ForeignKey:UserID"`
	Logins   []UserLogin  `gorm:"ForeignKey:UserID"`
}

func newContextWithUserID(users UserStorage, ctx context.Context, req *http.Request) context.Context {
	getVar := func(key string) string {
		return req.URL.Query().Get(key)
	}

	userID := ""
	if userID == "" {
		users.FindUserID(getVar("session"), getVar("ip"), getVar("mac"))
	}
	return context.WithValue(ctx, userIDKey, userID)
}

func userIDFromContext(users UserStorage, ctx context.Context) string {
	return ctx.Value(userIDKey).(string)
}

func UserMiddleware(users UserStorage, next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		ctx := newContextWithUserID(users, req.Context(), req)
		next.ServeHTTP(res, req.WithContext(ctx))
	})
}
