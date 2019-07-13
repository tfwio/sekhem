package ormus

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

// Session represents users who are logged in.
type Session struct {
	dbx
	ID      int64     `gorm:"auto_increment;unique_index;primary_key;column:id"`
	UserID  int64     `gorm:"column:user_id"` // [users].[id]
	Host    string    `gorm:"column:host"`    // running multiple server instance/port(s)?
	Created time.Time `gorm:"not null;column:created"`
	Expires time.Time `gorm:"not null;column:expires"`
	SessID  string    `gorm:"not null;column:sessid"`
	Client  string    `gorm:"not null;column:cli-key"` // .Request.RemoteAddr
}

// TableName Set User's table name to be `users`
func (Session) TableName() string {
	return "sessions"
}

// EnsureTableSessions creates table [sessions] if not exist.
func EnsureTableSessions() {
	var s Session
	db, _ := iniK("error(ensure-table-sessions) loading db; (expected)\n")
	// if !e {
	defer db.Close()
	if !db.HasTable(s) {
		db.CreateTable(s)
	}
	// }
}

// SessionCLIList returns a  sessions for CLI.
// The method first fetches a list of User elements
// then reports the Sessions with user-data (name).
func SessionCLIList() []Session {
	var users []User
	var sessions []Session
	db, err := iniC("error(session-cli-list) loading db\n")

	defer db.Close()
	if !err {
		usermap := UserGetList()
		for m, x := range users {
			fmt.Printf("--> %04d: %s\n", m, x.Name)
			usermap[x.ID] = x
		}
		// list sessions
		db.Find(&sessions)
		fmt.Printf("--> found %d entries\n", len(sessions))
		for _, x := range sessions {
			fmt.Printf("--> '%s'\n  CRD: %s\n  EXP: %s\n  SID: %s\n",
				usermap[x.UserID].Name,
				x.Created.Format("20060102_1504.005"),
				x.Expires.Format("20060102_1504.005"),
				x.SessID)
		}
	}
	return sessions
}

// SessionFind attempts to find a session from SessionID
func SessionFind(sessid string, host string) (Session, User) {
	var s Session
	var u User
	db, err := gorm.Open(datasys, datasource)
	defer db.Close()
	if err != nil {
		db.Where("[sessid] = ?", sessid).First(&s)
		db.Where("[id] = ?", s.UserID).First(&u)
		// Where("")
	} else {
		println("coundn't find session")
	}
	return s, u
}
