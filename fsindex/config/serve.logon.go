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

// LogonModel responds to a login action such as "/login/" or (perhaps) "/login-refresh/"
type LogonModel struct {
	Action string `json:"action"`
	Status bool   `json:"status"`
	Detail string `json:"detail"`
}

var (
	_safeHandlers   = wrapup(strings.Split("login,register,logout", ",")...)
	_unSafeHandlers = wrapup(strings.Split("json,json-index,pan,meta,json,refresh,tag,jtag", ",")...)
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

func (c *Configuration) sessMiddleware(g *gin.Context) {
	yn := false
	ck := false

	xi := g.ClientIP()
	ix := strings.Index(xi, ":")
	fmt.Printf("%d\n", ix)
	if ix != -1 {
		fmt.Printf("%s, %d\n", xi[:ix-1], ix)
	} else {
		fmt.Printf("%s, %d\n", xi, ix)
	}
	if isunsafe(g.Request.RequestURI) {
		yn = ormus.SessionCookieValidate(c.SessionHost("sekhem"), g)
		ck = true
		fmt.Printf("--> session? %v %s\n", yn, g.Request.RequestURI)
		fmt.Printf("==> CHECK: %v, VALID: %v\n", ck, yn)
	}
	g.Set("valid", yn)
	g.Next() // after request
	fmt.Printf("--> URI: %s\n", g.Request.RequestURI)
}

func (c *Configuration) serveLogout(g *gin.Context) {
	sh := c.SessionHost("sekhem")

	if sess, err := ormus.SessionCookie(sh, g.Request); !err {
		if time.Now().Before(sess.Expires) {
			fmt.Printf("==> SESSION EXISTS; LOGGIN OUT")
			ormus.SetCookieMaxAge(g, sh, sess.SessID, 0)
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
			ormus.SetCookieMaxAge(g, sh, sess.SessID, 0)
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
	} else if u.ValidateSessionByUserID(sh, g) {
		// we have a valid session (expired or not)
		xs := u.SessionRefresh(sh, g)
		if xs != nil {
			sh := c.SessionHost("sekhem")
			ss := xs.(ormus.Session) // d, _ := time.ParseDuration("12h") // exp := time.Now().Add(d)
			ormus.SetCookieMaxAge(g, sh, ss.SessID, ormus.ConstCookieAge12H)
			ormus.SetCookieSessOnly(g, sh+"_xo", u.Name)
		}
		j.Detail = "Prior session exists; updated."
		j.Status = true
	} else {
		hr := 12
		fmt.Printf("--> Create session for user: %v\n", u.Name)
		if e, sess := u.CreateSession(g, hr, sh, -1); !e {
			ormus.SetCookieMaxAge(g, sh, sess.SessID, ormus.ConstCookieAge12H)
			ormus.SetCookieSessOnly(g, sh+"_xo", u.Name)
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
		sh := c.SessionHost("sekhem")
		// create a session for the user
		e, sess := u.CreateSession(g, 12, sh, -1)
		if !e {
			ormus.SetCookieMaxAge(g, sh, sess.SessID, ormus.ConstCookieAge12H)
			ormus.SetCookieSessOnly(g, sh+"_xo", u.Name)
			j.Status = true
			j.Detail = "User and Session created."
		} else {
			j.Status = false
			j.Detail = "User created; session failed."
		}
	}
	g.JSON(http.StatusOK, j)
}
