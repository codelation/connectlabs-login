package ap

import (
	"time"

	"github.com/jinzhu/gorm"
)

//Session holds information about a WiFi user's session
type Session struct {
	gorm.Model
	Token     string
	Node      string
	IPv4      string
	Device    string
	Seconds   uint
	Download  uint
	Upload    uint
	ExpiresAt time.Time
	UserID    uint
}

//SessionStorage defines methods needed to privide session storage
type SessionStorage interface {
	FindSessionByID(id uint, out *Session) error
	FindSessionByToken(token string, out *Session) error
	FindSessionByUserID(id int, out *Session) error
	UpdateSessionFromRequest(*Request) error
}
