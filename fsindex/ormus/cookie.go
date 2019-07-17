package ormus

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	cookieSecure   = false
	cookieHTTPOnly = true
	cookieAgeSec   = ConstCookieAge48H
)

//
const (
	ConstCookieAge48H = /*60*60*/ 3600 * 48
	ConstCookieAge24H = /*60*60*/ 3600 * 24
	ConstCookieAge12H = /*60*60*/ 3600 * 12
)

// SetCookieMaxAge will set a cookie with our default settings.
//
// Will use `maxAge` in seconds to tell the browser
// to destroy the cookie as oppsed to using `SetCookieExpires`.
// See `CookieDefaults` in order to override default settings.
//
// Note: *Like `github.com/gogonic/gin`, we are applying `url.QueryEscape`
// `value` stored to the cookie so be sure to UnEscape the value when retrieved.*
func SetCookieMaxAge(cli *gin.Context, name string, value string, maxAge int) {
	http.SetCookie(cli.Writer, &http.Cookie{
		Name:     name,
		Value:    url.QueryEscape(value),
		MaxAge:   maxAge,
		Path:     "/",
		Secure:   cookieSecure,
		HttpOnly: cookieHTTPOnly,
	})
}

// SetCookieSessOnly will set a cookie with our default settings.
// Will expire with the browser session.
//
// See `CookieDefaults` in order to override default settings.
//
// Note: *Like `github.com/gogonic/gin`, we are applying `url.QueryEscape`
// `value` stored to the cookie so be sure to UnEscape the value when retrieved.*
func SetCookieSessOnly(cli *gin.Context, name string, value string) {
	http.SetCookie(cli.Writer, &http.Cookie{
		Name:     name,
		Value:    url.QueryEscape(value),
		Path:     "/",
		Secure:   cookieSecure,
		HttpOnly: cookieHTTPOnly,
	})
}

// SetCookieExpires will set a cookie with our default settings.
//
// Will use Expires (as opposed to `SetCookieMaxAge`).
// See `CookieDefaults` in order to override default settings.
//
// Note: *Like `github.com/gogonic/gin`, we are applying `url.QueryEscape`
// `value` stored to the cookie so be sure to UnEscape the value when retrieved.*
func SetCookieExpires(cli *gin.Context, name string, value string, expire time.Time) {
	http.SetCookie(cli.Writer, &http.Cookie{
		Name:     name,
		Value:    url.QueryEscape(value),
		Expires:  expire,
		Path:     "/",
		Secure:   cookieSecure,
		HttpOnly: cookieHTTPOnly,
	})
}

// CookieDefaults sets default cookie expire age and security.
func CookieDefaults(ageSec int, httpOnly bool, isSecure bool) {
	cookieAgeSec = ageSec
	cookieHTTPOnly = httpOnly
	cookieSecure = isSecure
}

func getCookie(cname string, client *gin.Context) *http.Cookie {
	var result *http.Cookie
	if xid, e := client.Request.Cookie(cname); e == nil {
		result = xid
	}
	return result
}

func getCookieValue(cname string, client *gin.Context) string {
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
func SessionCookieValidate(cookieName string, client *gin.Context) bool {

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
