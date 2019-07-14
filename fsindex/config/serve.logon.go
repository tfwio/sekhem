package config

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/tfwio/sekhem/util"

	"github.com/gin-gonic/gin"
	"github.com/tfwio/sekhem/fsindex/ormus"
)

var (
	_safeHandlers   = wrapup(strings.Split("login,register,logout", ",")...)
	_unSafeHandlers = wrapup(strings.Split("json,json-index,pan,meta,json,refresh,tag,jtag", ",")...)
	cookieSecure    = false
	cookieHTTPOnly  = true
)

func isunsafe(input string) bool {
	for _, unsafe := range _unSafeHandlers {
		if strings.Contains(input, unsafe) {
			return true
		}
	}
	return false
}
func wrapup(inputs ...string) []string {
	data := inputs
	for i, hander := range data {
		data[i] = strings.TrimRight(util.WReapLeft("/", hander), "/")
		// println(data[i])
	}
	return data
}

func (c *Configuration) sessMiddleware(context *gin.Context) {
	yn := false
	ck := false
	if isunsafe(context.Request.RequestURI) {
		yn = ormus.SessionCookieValidate(c.SessionHost("sekhem"), context.Request)
		ck = true
		fmt.Printf("--> session? %v %s\n", yn, context.Request.RequestURI)
		fmt.Printf("==> CHECK: %v, VALID: %v\n", ck, yn)
	}
	context.Set("valid", yn)
	context.Next() // after request
	fmt.Printf("--> URI: %s\n", context.Request.RequestURI)
}

func (c *Configuration) serveLogout(g *gin.Context) {
	sh := c.SessionHost("sekhem")

	if sess, err := ormus.SessionCookie(sh, g.Request); !err {
		if time.Now().Before(sess.Expires) {
			fmt.Printf("==> SESSION EXISTS; LOGGIN OUT")
			g.SetCookie(sh, sess.SessID, 0, "/", "", cookieSecure, cookieHTTPOnly)
			sess.Expires = time.Now()
			sess.Save()
			g.JSON(
				http.StatusOK,
				&LogonModel{
					Action: "logout",
					Detail: "Session exists; logged out.",
					Status: true})
		} else {
			fmt.Printf("==> SESSION EXP: %v\n", sess.Expires)
			g.SetCookie(sh, sess.SessID, 0, "/", "", cookieSecure, cookieHTTPOnly)
			sess.Expires = time.Now()
			sess.Save()
			g.JSON(
				http.StatusOK,
				&LogonModel{
					Action: "logout",
					Detail: fmt.Sprintf("Session expired; cookie stamped as expred now (%v).", sess.Expires),
					Status: false})
		}
	} else {
		fmt.Printf("==> SESSION NOT EXIST; NOTHING TO DO")
		g.JSON(
			http.StatusOK,
			&LogonModel{
				Action: "logout",
				Detail: "Session not exist; nothing to do.",
				Status: false})
	}
}

func (c *Configuration) serveLogin(g *gin.Context) {

	usr := g.Request.FormValue("user")
	j := LogonModel{Action: "login", Detail: "session creation failed.", Status: false}
	sh := c.SessionHost("sekhem")

	u := ormus.User{}
	if !u.ByName(usr) {
		println("--> NO user found!")
		j.Detail = "No user record."
		j.Status = false
	} else if u.ValidateSessionByUserID(sh, g.Request) {
		xs := u.SessionRefresh(sh, g.Request)
		if xs != nil {
			sh := c.SessionHost("sekhem")
			ss := xs.(ormus.Session) // d, _ := time.ParseDuration("12h") // exp := time.Now().Add(d)
			g.SetCookie(sh, ss.SessID, 12*3600, "/", "", cookieSecure, cookieHTTPOnly)
			g.SetCookie(sh+"_xo", u.Name, 12*3600, "/", "", cookieSecure, cookieHTTPOnly)
		}
		j.Detail = "Prior session exists; updated."
		j.Status = true
	} else {
		hr := 12
		fmt.Printf("--> Create session for user: %v\n", u.Name)
		if e, sess := u.CreateSession(g.Request, hr, sh, -1); !e {
			g.SetCookie(sh, sess.SessID, 12*3600, "/", "", cookieSecure, cookieHTTPOnly)
			g.SetCookie(sh+"_xo", u.Name, 12*3600, "/", "", cookieSecure, cookieHTTPOnly)
			j.Detail = "Session created."
			j.Status = true
		} else {
			j.Detail = "Session creation failed"
		}
	}
	g.JSON(http.StatusOK, j)
}

func (c *Configuration) serveRegister(g *gin.Context) {

	j := LogonModel{Action: "register", Detail: "user creation failed.", Status: false}
	usr := g.Request.FormValue("user")
	pas := g.Request.FormValue("pass")

	u := ormus.User{}
	fmt.Printf("!-> %s, %s, %v\n", usr, pas, -1)
	if e := u.Create(usr, pas, -1); e != 0 {
		switch e {
		case -1:
			j.Detail = "Failed to load db."
		case 1:
			j.Detail = "User record already exists."
		}
	} else {
		hr := 12
		sh := c.SessionHost("sekhem")
		e, sess := u.CreateSession(g.Request, hr, sh, -1)
		if !e {
			g.SetCookie(sh, sess.SessID, 12*3600, "/", "", cookieSecure, cookieHTTPOnly)
			g.SetCookie(sh+"_xo", sess.SessID, 12*3600, "/", "", cookieSecure, cookieHTTPOnly)
			j.Status = true
			j.Detail = "User and Session created."
		} else {
			j.Status = false
			j.Detail = "User created; session failed."
		}
	}
	g.JSON(http.StatusOK, j)
}
