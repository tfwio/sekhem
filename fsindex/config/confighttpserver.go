package config

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/tfwio/session"
	"github.com/tfwio/srv/fsindex"
	"github.com/tfwio/srv/util"
)

var (
	// UseTLS is a cli (console) override for using TLS or not.
	// It can only be put on in the console and if `--use-tls` is not set
	// then this will always remain false.
	UseTLS = false
	// UsePORT is a CLI override.  If set to a value other than -1 it will be used.
	UsePORT uint = 5500
	// UseHost is set by CLI to override the config-file Host setting..
	UseHost = ""
	// DefaultConfigFile — you know.  Default = `./data/conf`.
	DefaultConfigFile = util.Abs("./data/conf.json")
	// DefaultDatabase this is our default database (file-path: `[configfile-dir]/ormus.db`).
	DefaultDatabase = util.CatPath(util.GetDirectory(DefaultConfigFile), "ormus.db")
	// DefaultDatasys default data system or provider ('sqlite3').
	DefaultDatasys = "sqlite3"
	extMap         map[string]*fsindex.FileSpec
	mdlMap         map[string]*fsindex.Model
	models         []fsindex.Model
	xCounter       int32
	fCounter       int32
)
var sessions = session.Service{
	AppID:               "sekhem",
	CookieHTTPOnly:      true,  // hymmm
	CookieSecure:        false, // we want to see em in the browser
	KeySessionIsValid:   "valid",
	KeySessionIsChecked: "checked",
	AdvanceOnKeepYear:   0, // 0
	AdvanceOnKeepMonth:  6, // 6
	AdvanceOnKeepDay:    0, // 0
	URICheck:            session.WrapURIExpression("/json-index/?$,/json/?$"),
	URIEnforce:          []string{},
	URIMatchHandler: func(uri, unsafe string) bool {
		tomatch := fmt.Sprintf("^%s", unsafe)
		// fmt.Fprintf(os.Stderr, "checking: %s\n", tomatch)
		if match, err := regexp.MatchString(tomatch, uri); err == nil {
			return match
		}
		return false
	},
	FormSession: session.FormSession{User: "user", Pass: "pass", Keep: "keep"},
}

// GinConfigure configures gin.Engine.
func (c *Configuration) GinConfigure(andServe bool, router *gin.Engine) {

	DefaultFile := util.Abs(util.Cat(c.Root.Directory, "\\", c.Root.Default))
	if andServe {

		sessions.Port = c.Port

		session.SetupService(
			&sessions, router,
			c.DatabaseType, util.Abs(util.CatPath(util.GetDirectory(util.Abs(DefaultConfigFile)), c.Database)),
			-1, -1)
		fmt.Printf("default\n  > Target = %-18s, Source =  %s\n", c.Root.Path, DefaultFile)

		router.StaticFile(c.Root.Path, DefaultFile)

		println("alias-default")
		for _, rootEntry := range c.Root.AliasDefault {
			target := util.Cat(c.Root.Path, rootEntry)
			router.StaticFile(target, DefaultFile)
			fmt.Printf("  > Target = %-18s, Source = %s\n", target, c.DefaultFile())
		}

		println("root-files")
		for _, rootEntry := range c.Root.Files {
			target := util.Cat(c.Root.Path, rootEntry)
			source := util.Abs(util.Cat(c.Root.Directory, "\\", rootEntry))
			router.StaticFile(target, source)
			fmt.Printf("  > Target = %-18s, Source = %s\n", target, source)
		}

		println("root-files: allowed")
		if c.Root.Allow != "" {
			allowed := strings.Split(c.Root.Allow, ",")
			for i := range allowed {
				allowed[i] = strings.Trim(allowed[i], " ")
				allowed[i] = strings.Trim(allowed[i], "\n")
				target := util.Cat(c.Root.Path, allowed[i])
				source := util.Abs(util.Cat(c.Root.Directory, "\\", allowed[i]))
				fmt.Printf("  > Target = %-18s, Source = %s\n", target, source)
				router.StaticFile(target, source)
			}
		}

		println("locations")
		for _, tgt := range c.Locations {
			router.StaticFS(tgt.Target, gin.Dir(util.Abs(tgt.Source), tgt.Browsable))
			fmt.Printf("  > Target = %-18s, Source = %s\n", tgt.Target, tgt.Source)
		}

		router.Any("/json-index", c.serveJSONIndex)
		router.Any("/pan/:path/*action", func(g *gin.Context) { c.servePandoc(c.Pandoc.HTMLTemplate, pandoctemplate, g) })
		router.Any("/meta/:path/*action", func(g *gin.Context) { c.servePandoc(c.Pandoc.MetaTemplate, pandoctemplate, g) })
	}
	c.initializeModels()

	if andServe {
		c.serveModelIndex(router)
	}
}

const (
	constServerDefaultHost    = "localhost"
	constServerDefaultPort    = ":5500"
	constServerTLSCertDefault = "data\\cert.pem"
	constServerTLSKeyDefault  = "data\\key.pem"
	constDefaultDataPath      = "./data"
	constConfJSONReadSuccess  = "got JSON configuration"
	constConfJSONReadError    = "Error: failed to read JSON configuration. %s\n"
	constMessageWroteJSON     = `
We've exported a default "%s" for your editing.

Modify the file per your preferences.
`
	constRootAliasDefault     = "home,index.htm,index.html,index,default,default.htm"
	constRootFilesDefault     = "json.json,bundle.js,favicon.ico"
	constRootIndexDefault     = "index.html"
	constRootDirectoryDefault = "public"
	constRootPathDefault      = "/"
	constStaticSourceDefault  = "public/static"
	constStaticTargetDefault  = "/static/"
	constImagesSourceDefault  = "public/images"
	constImagesTargetDefault  = "/images/"
	constExtDefaultMedia      = ".mp4,.m4a,.mp3"
	constExtDefaultMMD        = ".md,.mmd"
)
