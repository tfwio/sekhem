package config

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"tfw.io/Go/fsindex/fsindex"
	"tfw.io/Go/fsindex/util"
)

var (
	// UseTLS — you know.
	UseTLS = false
	// DefaultConfigFile — you know.
	DefaultConfigFile = "data/conf.json"
	extMap            map[string]*fsindex.FileSpec
)

// GinConfig configures gin.Engine.
func (c *Configuration) GinConfig(mGin *gin.Engine, paths ...*fsindex.Model) {

	// these files are all stored in the public directory.
	// they are the only files we're serving specifically in
	// that directory.

	DefaultFile := util.Abs(util.Cat(c.Root.Directory, "\\", c.Root.Default))
	// c.configInfo()

	// println("==> server setup: c.Root.AliasDefault")
	println("alias-default")
	for _, rootEntry := range c.Root.AliasDefault {
		target := util.Cat(c.Root.Path, rootEntry)
		mGin.StaticFile(target, DefaultFile)
		fmt.Printf("  ? Target = %s, Source = %s\n", target, c.DefaultFile())
	}
	println("root-files")
	for _, rootEntry := range c.Root.Files {
		target := util.Cat(c.Root.Path, rootEntry)
		source := util.Abs(util.Cat(c.Root.Directory, "\\", rootEntry))
		mGin.StaticFile(target, source)
		fmt.Printf("  > Target = %s, Source = %s\n", target, source)
	}
	println("locations")
	for _, tgt := range c.Locations {
		println("- serving path:", tgt.Target, "from", util.Abs(tgt.Source))
		mGin.StaticFS(tgt.Target, gin.Dir(util.Abs(tgt.Source), tgt.Browsable))
		fmt.Printf("  > Target = %s, Source = %s\n", tgt.Target, tgt.Source)
	}

	fmt.Printf("- default: Target = %s, Source =  %s\n", c.Root.Path, DefaultFile)
	mGin.StaticFile(c.Root.Path, DefaultFile)

	for _, path := range paths {
		p := util.Cat("/json/", path.Name)
		println(fmt.Sprintf("--> adding JSON %s", p))
		mGin.GET(p, func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, &path)
		})
	}
}

const (
	constServerDefaultHost    = "localhost"
	constServerDefaultPort    = ":5500"
	constServerTLSCertDefault = "data\\ia.crt"
	constServerTLSKeyDefault  = "data\\ia.key"
	constDefaultDataPath      = "./data"
	constConfJSONReadSuccess  = "got JSON configuration"
	constConfJSONReadError    = "Error: failed to read JSON configuration. %s\n"
	constMessageWroteJSON     = `
We've exported a default data/conf.json for your editing.

Modify the file per your preferences.

[terminating application]
`
	constRootAliasDefault     = "home,index.htm,index.html,index,default,default.htm"
	constRootFilesDefault     = "json.json,bundle.js,favicon.ico"
	constRootIndexDefault     = "index.html"
	constRootDirectoryDefault = ".\\public"
	constRootPathDefault      = "/"
	constStaticSourceDefault  = "public\\static"
	constStaticTargetDefault  = "/static/"
	constImagesSourceDefault  = "public\\images"
	constImagesTargetDefault  = "/images/"
	constExtDefaultMedia      = ".mp4,.m4a,.mp3"
	constExtDefaultMMD        = ".md,.mmd"
)
