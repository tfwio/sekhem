package ormus

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
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

// SessionValidateCookie checks against a provided salt and hash.
// BUT FIRST, it checks for a valid session?
func SessionValidateCookie(host string, client *http.Request) bool {

	clistr := getClientString(client)
	sessid := ""
	if xid, e := client.Cookie(host); e == nil {
		sessid, _ = url.QueryUnescape(xid.Value)
	}

	if sessid == "" {
		return false
	}

	result := false

	db, err := iniC("error(validate-session) loading database\n")
	if err {
		return false
	}
	db.LogMode(true)
	sess := Session{}
	defer db.Close()
	db.First(&sess, "[cli-key] = ? AND [host] = ? AND [sessid] = ?", clistr, host, sessid)
	fmt.Printf("%s\n%s\n", sessid, sess.SessID)
	if sess.SessID == sessid {
		result = true
	}

	return result
}
