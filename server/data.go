package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	//This is required for the postgres driver within gorm
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/ryanhatfield/connectlabs-login/ap"
	"github.com/ryanhatfield/connectlabs-login/sso"
)

//DataError is a generic error type for reference
type DataError string

func (e DataError) Error() string { return string(e) }

//ErrModelNotFound is returned when the model requested wasn't found in the db
const ErrModelNotFound = DataError("Model Not Found")

//ErrInvalidRequest is returned when the request or request parameters are invalid
const ErrInvalidRequest = DataError("Invalid Parameters when calling Data function")

//Data holds methods to get info to/from the db
type Data struct {
	sso.UserStorage
	ap.SessionStorage
	Context          context.Context
	DatabaseURL      string
	Debug            bool
	MaxDBConnections int
	// Store            sessions.Store

	db *gorm.DB
}

//InitializeDB opens the database connection and attempts to migrate the database
func (d *Data) InitializeDB() error {
	log.Println("Initializing DB adapter")
	db, err := gorm.Open("postgres", d.DatabaseURL)
	if err != nil {
		return err
	}

	db.DB().SetMaxOpenConns(d.MaxDBConnections)
	db.DB().SetMaxIdleConns(d.MaxDBConnections)
	db.AutoMigrate(&ap.Session{})
	db.AutoMigrate(&sso.User{})
	db.AutoMigrate(&sso.UserLogin{})

	if d.Debug {
		d.db = db.Debug()
	} else {
		d.db = db
	}

	return nil
}

/****************************/
/* SessionStorage functions */
/****************************/

func (d *Data) FindSessionByID(id uint, out *ap.Session) error {
	d.db.Find(out, id)
	if d.db.NewRecord(out) {
		return ErrModelNotFound
	}
	return nil
}

func (d *Data) FindSessionByToken(token string, out *ap.Session) error {
	if token == "" {
		return ErrInvalidRequest
	}

	d.db.Where("token = ?", token).Order("expires_at desc").Find(out)

	if d.db.NewRecord(out) {
		return ErrModelNotFound
	}

	return nil
}

func (d *Data) FindSessionByUserID(id int, out *ap.Session) error {

	d.db.Where("user_id = ?", id).Order("expires_at desc").Find(out)

	if d.db.NewRecord(out) {
		return ErrModelNotFound
	}

	return nil
}

func (d *Data) UpdateSessionFromRequest(req *ap.Request) error {
	session := &ap.Session{}
	d.FindSessionByToken(req.Session, session)

	session.IPv4 = req.IPV4Address
	if req.RequestType == ap.RequestTypeAccounting {
		getint := func(s string) uint {
			u, _ := strconv.ParseUint(req.Download, 0, 32)
			return uint(u)
		}
		session.Download = getint(req.Download)
		session.Upload = getint(req.Upload)
		session.Seconds = getint(req.Seconds)
	}

	if d.db.NewRecord(session) {
		//session not found, save a new instance
		session.Token = req.Session
		session.Device = req.MacAddress
		session.Node = req.NodeAddress
		session.ExpiresAt = time.Now()
		d.db.Save(&session)
	} else {
		d.db.Model(&session).Updates(ap.Session{
			Download: session.Download,
			Upload:   session.Upload,
			IPv4:     session.IPv4,
			Seconds:  session.Seconds,
		})
	}

	return nil
}

/****************************/
/* UserStorage functions */
/****************************/

//FindUserByID returns a user object from the user id
func (d *Data) FindUserByID(id uint, out *sso.User) error {

	d.db.First(out, id)

	if d.db.NewRecord(out) {
		return ErrModelNotFound
	}
	return nil
}

//FindUserByDevice attempts to find a session. If a session is found, the associated
//	user is returned. If there isn't a user associated with that session, the user is
//  created and associated with the session before being returned.
func (d *Data) FindUserByDevice(mac string, node string, out *sso.User) error {
	s := &ap.Session{}

	if mac == "" || node == "" {
		return fmt.Errorf("both mac and node are required, mac: %s, node: %s", mac, node)
	}

	d.db.Where("node = ? AND device = ?", node, mac).Order("expires_at desc").First(s)
	if d.db.NewRecord(s) {
		return fmt.Errorf("could not find user with mac: %s, node: %s", mac, node)
	}
	if !d.db.NewRecord(s) {
		d.FindUserByID(s.UserID, out)
		if d.db.NewRecord(out) {
			out = &sso.User{}
			d.db.Create(out)
		}

		if d.db.NewRecord(out) {
			//if it's still a new record, we have an issue
			return fmt.Errorf("could not save user to db, user: %+v", *out)
		}

		d.db.Model(s).UpdateColumn("user_id", out.ID)
		d.db.Model(out).Association("Sessions").Find(&out.Sessions)
		d.db.Model(out).Association("Logins").Find(&out.Logins)

		log.Printf("%+v", out)
	}

	return nil
}

func (d *Data) FindUserIDByDevice(token, mac, node string) string { return "" }

func (d *Data) AddLoginToUser(userID uint, login sso.UserLogin) error {
	user := &sso.User{}

	if err := d.FindUserByID(userID, user); err != nil {
		return err
	}

	d.db.Model(user).Association("Logins").Find(&user.Logins)

	updated := false
	for _, j := range user.Logins {
		if j.Provider == login.Provider {
			d.db.Model(&j).Updates(login)
			updated = true
		}
	}

	if !updated {
		d.db.Model(user).Association("Logins").Append(login)

		//make sure it saved
		if d.db.NewRecord(login) {
			return errors.New("error associating login with user")
		}
	}

	return nil
}

func (d *Data) FindUserID(session, ip, mac string) string {
	s := &ap.Session{}

	d.db.Find(s, struct{ Token string }{Token: session})

	if !d.db.NewRecord(s) {
		return fmt.Sprint(s.UserID)
	}

	d.db.Find(s, struct{ IPv4, Device string }{ip, mac})

	if !d.db.NewRecord(s) {
		return fmt.Sprint(s.UserID)
	}

	return ""
}
