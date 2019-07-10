package config

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tfwio/sekhem/fsindex"
	"github.com/tfwio/sekhem/util"
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
	extMap            map[string]*fsindex.FileSpec
	mdlMap            map[string]*fsindex.Model
	models            []fsindex.Model
	xCounter          int32
	fCounter          int32
)

// GinConfig configures gin.Engine.
func (c *Configuration) GinConfig(router *gin.Engine) {
	c.GinConfigure(true, router)
}

// GinConfigure configures gin.Engine.
// if justIndex is set to true, we just rebuild our indexes.
// We currently are not exposing this to http as our "/refresh/:target"
// path allows us to refresh a single index as needed.
func (c *Configuration) GinConfigure(andServe bool, router *gin.Engine) {

	// these files are all stored in the public directory.
	// they are the only files we're serving specifically in
	// that directory.

	DefaultFile := util.Abs(util.Cat(c.Root.Directory, "\\", c.Root.Default))

	fmt.Printf("default\n  > Target = %-18s, Source =  %s\n", c.Root.Path, DefaultFile)
	if andServe {
		router.StaticFile(c.Root.Path, DefaultFile)
	}

	if andServe {
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
	}

	xdata := JSONIndex{} // xdata indexes is just a string array map.
	xdata.Index = []string{}
	println("location indexes #1: string-map")
	for _, path := range c.Indexes {
		jsonpath := util.WReap("/", "json", util.AbsBase(path.Source))
		xdata.Index = append(xdata.Index, jsonpath)
	}
	println("JSON-index Target \"/json-index\"")
	router.GET("/json-index", func(g *gin.Context) {
		g.JSON(http.StatusOK, xdata)
	})

	router.GET("/pan/:path/*action", func(g *gin.Context) {
		c.servePandoc(c.Pandoc.HTMLTemplate, pandoctemplate, g)
	})

	router.GET("/meta/:path/*action", func(g *gin.Context) {
		c.servePandoc(c.Pandoc.MetaTemplate, pandoctemplate, g)
	})

	c.initializeModels()
	c.serveModelIndex(router)
}

func (c *Configuration) serveModelIndex(router *gin.Engine) {
	println("location indexes #2: primary")
	for _, path := range c.Indexes {
		jsonpath := util.WReap("/", "json", util.AbsBase(path.Source))
		modelpath := util.WReap("/", path.Target)
		fmt.Printf("  > Target = %-18s, json = %s,  Source = %s\n", modelpath, c.GetPath(jsonpath), path.Source)
		modelpath = c.getIndexTarget(&path)

		if path.Servable {
			router.StaticFS(modelpath, gin.Dir(util.Abs(path.Source), path.Browsable))
		}
	}
	router.GET("/json/:route", c.serveJSON)
	println("/tag/ handler")
	router.GET("/refresh/:route", c.refreshRouteJSON)
	router.GET("/tag/:route/*action", func(g *gin.Context) { TagHandler(c, g) })
	router.GET("/jtag/:route/*action", func(g *gin.Context) { TagHandlerJSON(c, g) })
}

func (c *Configuration) serveJSON(ctx *gin.Context) {

	mroute := ctx.Param("route")

	if c.hasModel(mroute) {
		mmdl := mdlMap[mroute]
		fmt.Printf("SERVING JSON FOR: %v, %s\n", mroute, mmdl.FullPath)
		ctx.JSON(http.StatusOK, &mmdl.PathEntry)
	} else {
		jsi := JSONIndex{Index: []string{fmt.Sprintf("COULD NOT find model for index: %s", mroute)}}
		ctx.JSON(http.StatusNotFound, &jsi)
		fmt.Printf("--> COULD NOT FIND ROUTE %s\n", mroute)
	}
}

func (c *Configuration) refreshRouteJSON(g *gin.Context) {
	mroute := g.Param("route")
	jsi := JSONIndex{Index: []string{fmt.Sprintf("FOUND model for index: %s", mroute)}}
	if ndx, ok := c.indexFromTarget(mroute), c.hasModel(mroute); ok && ndx != nil {
		c.initializeModel(ndx)
		g.JSON(http.StatusOK, jsi)
		return
	}
	jsi = JSONIndex{Index: []string{fmt.Sprintf("COULD NOT find model for index: %s", mroute)}}
	g.JSON(http.StatusOK, &jsi)
	fmt.Printf("ERROR> COULD NOT find model for index: %s\n", mroute)
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
