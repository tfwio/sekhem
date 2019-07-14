package ormus

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

func getCookie(cname string, client *http.Request) *http.Cookie {
	var result *http.Cookie
	if xid, e := client.Cookie(cname); e == nil {
		result = xid
	}
	return result
}

func getCookieValue(cname string, client *http.Request) string {
	cookie := getCookie(cname, client)
	cookieValue := ""
	if cookie != nil {
		if sessid, x := url.QueryUnescape(cookie.Value); x == nil {
			cookieValue = sessid
		}
	}
	return cookieValue
}

func cookieValue(cookie *http.Cookie) string {
	cookieValue := ""
	if cookie != nil {
		if sessid, x := url.QueryUnescape(cookie.Value); x == nil {
			cookieValue = sessid
		}
	}
	return cookieValue
}

// SessionCookieValidate checks against a provided salt and hash.
// BUT FIRST, it checks for a valid session?
func SessionCookieValidate(cookieName string, client *http.Request) bool {

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

// SessionCookie retrieves a cookie provided client-host.
func SessionCookie(host string, client interface{}) (Session, bool) {
	clistr := getClientString(client)
	sess := Session{}
	db, err := iniC("error(validate-session) loading database\n")
	if err {
		return sess, false
	}
	db.LogMode(true)
	defer db.Close()
	db.First(&sess, "[cli-key] = ? AND [host] = ?", clistr, host)
	return sess, db.RowsAffected != 0
}
