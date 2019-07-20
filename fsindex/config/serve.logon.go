// +build session

package config

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/tfwio/sekhem/util"

	"github.com/gin-gonic/gin"
	"github.com/tfwio/sekhem/fsindex/session"
)

// LogonModel responds to a login action such as "/login/" or (perhaps) "/login-refresh/"
type LogonModel struct {
	Action string `json:"action"`
	Status bool   `json:"status"`
	Detail string `json:"detail"`
	Data   string `json:"data,omitempty"`
}

var (
	_safeHandlers   = wrapup(strings.Split("login,register,logout", ",")...)
	_unSafeHandlers = wrapup(strings.Split("json,json-index,pan,meta,json,refresh,tag,jtag", ",")...)
)

func (c *Configuration) initServerLogin(router *gin.Engine) bool {
	// fmt.Println("--> LOGON SESSIONS SUPPORTED")
	router.Use(c.sessMiddleware)
	router.Any("/logout/", c.serveLogout)
	router.Any("/login/", c.serveLogin)
	router.Any("/register/", c.serveRegister)
	router.Any("/stat/", c.serveUserStatus)
	return true
}

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
	// ck := false

	if isunsafe(g.Request.RequestURI) {
		yn = session.SessionCookieValidate(c.SessionHost("sekhem"), g)
		// ck = true
		// fmt.Printf("--> session? %v %s\n", yn, g.Request.RequestURI)
		// fmt.Printf("==> CHECK: %v, VALID: %v\n", ck, yn)
	}
	g.Set("valid", yn)
	g.Next() // after request
	// fmt.Printf("--> URI: %s\n", g.Request.RequestURI)
}

// serveUserStatus serves JSON checking if a session exists,
// persists, and a user exists.
// `{status: true,  detail: "found", data: <username>}` if all checks out,
// `{status: false, detail: "exists"}` if not logged in and
// `{status: false, detail: "none"}` if no user was found.
func (c *Configuration) serveUserStatus(g *gin.Context) {
	sh := c.SessionHost("sekhem")
	fmt.Println("==> CHECKING USER STATUS")
	if sess, success := session.SessionCookie(sh, g); success {
		fmt.Printf("  ==> CLIENT COOKIE EXISTS; USER=%d\n", sess.UserID)
		if u, success := sess.GetUser(); success && sess.IsValid() {
			g.JSON(
				http.StatusOK,
				&LogonModel{
					Action: "user-info",
					Detail: "found",
					Status: true,
					Data:   u.Name,
				})
		} else {
			g.JSON(
				http.StatusOK,
				&LogonModel{
					Action: "user-info",
					Detail: "exists",
					Status: false})
		}
	} else {
		g.JSON(
			http.StatusOK,
			&LogonModel{
				Action: "user-info",
				Detail: "none",
				Status: false})
	}
}

func (c *Configuration) serveLogout(g *gin.Context) {
	sh := c.SessionHost("sekhem")
	fmt.Println("==> LOGOUT ATTEMPT")
	sess, success := session.SessionCookie(sh, g)
	if success {
		fmt.Printf("  ==> CLIENT COOKIE EXISTS; USER=%d\n", sess.UserID)
		session.SetCookieMaxAge(g, sh, sess.SessID, 0)
		sess.Expires = time.Now()
		if time.Now().Before(sess.Expires) {
			// Found matching session from browser-cookie
			fmt.Printf("  --> NOT EXPIRED; USER=%d\n", sess.UserID)
			sess.Save()
			g.JSON(
				http.StatusOK,
				&LogonModel{
					Action: "logout",
					Detail: "Session exists; logged out.",
					Status: true})
		} else {
			fmt.Printf("  --> SESSION EXP: %v\n", sess.Expires)
			sess.Save()
			g.JSON(
				http.StatusOK,
				&LogonModel{
					Action: "logout",
					Detail: fmt.Sprintf("User was logged out prior; Logout re-enforced."),
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

	fmt.Println("==> LOGIN REQUEST")

	usr := g.Request.FormValue("user")
	j := LogonModel{Action: "login", Detail: "session creation failed.", Status: false}
	sh := c.SessionHost("sekhem")

	u := session.User{}
	if !u.ByName(usr) {
		println("  --> USER NOT FOUND!")
		j.Detail = "No user record."
		j.Status = false
	} else {
		// We have a valid user;
		sess, success := u.UserSession(sh, g)
		if success {
			fmt.Println("  ==> FOUND USER, VALIDATING PW")
			if fpass := g.Request.FormValue("pass"); fpass != "" {
				if u.ValidatePassword(fpass) {
					fmt.Println("  ==> PW:GOOD")
					sess.Refresh(true)
					session.SetCookieMaxAge(g, sh, sess.SessID, session.ConstCookieAge12H)
					session.SetCookieSessOnly(g, sh+"_xo", u.Name)
					j.Detail = "Logged in."
					j.Status = true
				} else {
					fmt.Println("  ==> PW:FAIL")
					j.Detail = "Password did not match."
					j.Status = true
				}
			}
		} else {
			// There is no session for the user.
			// this use case shouldn't exist since a session is created when a user is created!
			// this should report success to spite the fact that its a failure.
			fmt.Println("  ==> DESTROY SESSION")
			sess.Destroy(true)
			session.SetCookieMaxAge(g, sh, sess.SessID, -1)
			session.SetCookieMaxAge(g, sh+"_xo", "", -1)
			j.Detail = "Session destroyed."
			j.Status = true
		}
	}
	g.JSON(http.StatusOK, j)
}

func (c *Configuration) serveRegister(g *gin.Context) {

	j := LogonModel{Action: "register", Detail: "user creation failed.", Status: false}
	usr := g.Request.FormValue("user")
	pas := g.Request.FormValue("pass")

	u := session.User{}
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
			session.SetCookieMaxAge(g, sh, sess.SessID, session.ConstCookieAge12H)
			session.SetCookieSessOnly(g, sh+"_xo", u.Name)
			j.Status = true
			j.Detail = "User and Session created."
		} else {
			j.Status = false
			j.Detail = "User created; session failed."
		}
	}
	g.JSON(http.StatusOK, j)
}
