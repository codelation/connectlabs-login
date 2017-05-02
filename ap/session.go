package ap

import (
	"time"

	"github.com/jinzhu/gorm"
)

//Session holds information about a WiFi user's session
type Session struct {
	gorm.Model
	Session   string
	Node      string
	IPv4      string
	Device    string
	Seconds   uint
	Download  uint
	Upload    uint
	ExpiresAt time.Time
	UserID    uint
}
