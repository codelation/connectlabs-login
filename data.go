package main

import (
	"errors"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	//This is required for the postgres driver within gorm
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/ryanhatfield/connectlabs-login/ap"
	"github.com/ryanhatfield/connectlabs-login/sso"
)

type Data interface {
	ap.SessionStorage
	sso.UserStorage
}

type data struct {
	DatabaseURL      string
	MaxDBConnections int
	db               *gorm.DB
}

func (d *data) InitializeDB() error {
	db, err := gorm.Open("postgres", d.DatabaseURL)
	if err != nil {
		return err
	}

	db.DB().SetMaxOpenConns(d.MaxDBConnections)
	db.DB().SetMaxIdleConns(d.MaxDBConnections)
	db.AutoMigrate(&ap.Session{})
	db.AutoMigrate(&sso.User{})
	db.AutoMigrate(&sso.UserLogin{})

	d.db = db
	return nil
}

func (d *data) FindSession(session *ap.Session, sessionID string) error {
	if sessionID == "" {
		return errors.New("session id empty")
	}
	d.db.First(session, ap.Session{Session: sessionID})

	return nil
}

func (d *data) UpdateSession(session ap.Session, req *ap.Request) error {
	s := &ap.Session{}

	d.FindSession(s, req.Session)

	if s.Session != req.Session {
		//session not found, save a new instance
		s.Session = req.Session
		s.Device = req.MacAddress
		s.Node = req.NodeAddress
		s.ExpiresAt = time.Now()
		defer d.db.Save(&s)
	} else {
		defer d.db.Model(&s).Updates(ap.Session{
			Download: s.Download,
			Upload:   s.Upload,
			IPv4:     s.IPv4,
			Seconds:  s.Seconds,
		})
	}

	s.IPv4 = req.IPV4Address
	if req.RequestType == ap.RequestTypeAccounting {
		getint := func(s string) uint {
			u, _ := strconv.ParseUint(req.Download, 0, 32)
			return uint(u)
		}
		s.Download = getint(req.Download)
		s.Upload = getint(req.Upload)
		s.Seconds = getint(req.Seconds)
	}

	return nil
}

func (d *data) FindUser(user *sso.User, userID uint) error {

	d.db.First(user, sso.User{ID: userID})

	return nil
}

func (d *data) AddLoginToUser(userID uint, login sso.UserLogin) error {
	user := &sso.User{}

	d.FindUser(user, userID)
	d.db.Create(login)
	//make sure it saved
	if !d.db.NewRecord(login) {
		return errors.New("error associating login with user.")
	}

	user.Logins = append(user.Logins, login)
	d.db.Save(&user)
	return nil
}
