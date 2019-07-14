package ormus

import (
	"fmt"
	"net/http"
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

// SessionValidateCookie checks against a provided salt and hash.
// BUT FIRST, it checks for a valid session?
func SessionValidateCookie(cookieName string, client *http.Request) bool {

	clistr := getClientString(client)
	cookie := getCookie(cookieName, client)
	sessid := cookieValue(cookie)

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
	db.First(&sess, "[cli-key] = ? AND [host] = ? AND [sessid] = ?", clistr, cookieName, sessid)
	fmt.Printf("SESS\nsess: %s\ncook: %s\n", sess.SessID, sessid)
	fmt.Printf("EXPR\nsess: %v\ncook: %v\n", sess.Expires, cookie.Expires)

	if sess.SessID == sessid {
		result = time.Now().Before(sess.Expires)
	}

	return result
}

// Save session data to db.
func (s *Session) Save() bool {
	db, err := s.iniC("error(validate-session) loading database\n")
	if err {
		return false
	}
	defer db.Close()
	db.Save(s)
	return db.RowsAffected > 0
}
