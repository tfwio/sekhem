package config

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

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
	xCounter          int32
	fCounter          int32
)

// GinConfig configures gin.Engine.
func (c *Configuration) GinConfig(mGin *gin.Engine) {

	// these files are all stored in the public directory.
	// they are the only files we're serving specifically in
	// that directory.

	DefaultFile := util.Abs(util.Cat(c.Root.Directory, "\\", c.Root.Default))

	fmt.Printf("default\n  > Target = %-18s, Source =  %s\n", c.Root.Path, DefaultFile)
	mGin.StaticFile(c.Root.Path, DefaultFile)

	println("alias-default")
	for _, rootEntry := range c.Root.AliasDefault {
		target := util.Cat(c.Root.Path, rootEntry)
		mGin.StaticFile(target, DefaultFile)
		fmt.Printf("  > Target = %-18s, Source = %s\n", target, c.DefaultFile())
	}

	println("root-files")
	for _, rootEntry := range c.Root.Files {
		target := util.Cat(c.Root.Path, rootEntry)
		source := util.Abs(util.Cat(c.Root.Directory, "\\", rootEntry))
		mGin.StaticFile(target, source)
		fmt.Printf("  > Target = %-18s, Source = %s\n", target, source)
	}

	println("locations")
	for _, tgt := range c.Locations {
		mGin.StaticFS(tgt.Target, gin.Dir(util.Abs(tgt.Source), tgt.Browsable))
		fmt.Printf("  > Target = %-18s, Source = %s\n", tgt.Target, tgt.Source)
	}

	println("location indexes")
	for _, path := range c.Indexes {

		p := util.Wrapper("/", "json", filepath.Base(path.Source))
		m := c.createEntry(path, fsindex.DefaultSettings)

		np := util.Wrapper("/", path.Target, m.Name)
		if m.Settings.OmitRootNameFromPath {
			np = util.Wrap(path.Target, "/")
		}

		fmt.Printf("  > Target = %-18s, Source = %s\n", np, path.Source)
		println(fmt.Sprintf("  > %s = %-18s", "JSON", p))

		mGin.StaticFS(np, gin.Dir(util.Abs(path.Source), path.Browsable))

		buildFileSystemModel(&m)
		mGin.GET(p, func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, &m)
		})
	}
}

func (c *Configuration) createEntry(path IndexPath, settings fsindex.Settings) fsindex.Model {
	println("- path for indexed files: ", c.GetPath(path.Target))
	// configure, createIndex, checkSimpleModel
	pe := fsindex.Model{
		PathEntry: fsindex.PathEntry{
			PathSpec: fsindex.PathSpec{
				FileEntry: fsindex.FileEntry{
					Parent:   nil,
					Name:     util.AbsBase(path.Source),
					FullPath: util.Abs(path.Source),
					SHA1:     util.Sha1String(path.Source),
				},
				IsRoot: true},
			FauxPath:    c.GetPath(path.Target),
			FileFilter:  c.Extensions,
			IgnorePaths: []string{},
		},
		Settings: settings,
	}
	buildFileSystemModel(&pe)
	return pe
}

func buildFileSystemModel(model *fsindex.Model) {

	xCounter, fCounter = 0, 0

	model.SimpleModel = fsindex.SimpleModel{}
	model.CreateMaps()

	handler := fsindex.Handlers{
		ChildPath: func(root *fsindex.Model, child *fsindex.PathEntry) bool {
			model.AddPath(root, child)
			return false
		},
		ChildFile: func(root *fsindex.Model, child *fsindex.FileEntry) bool {
			ext := strings.ToLower(filepath.Ext(child.FullPath))
			if ext == ".md" {
				// datestring := checkDateString(child.Base())
				model.AddFile(root, child)
			}
			fCounter++
			return false
		},
	}

	model.Refresh(model, &xCounter, &handler)

	// checkSimpleModel(&mdl)
}

func checkSimpleModel(mdl *fsindex.SimpleModel, pathEntry *fsindex.Model) {
	// map counters don't yield adequate length
	println("File map Count: ", len(mdl.File))
	println("Path map Count: ", len(mdl.Path))
	//
	println("File Count: ", fCounter)
	println("Path Count: ", xCounter)

	ref1 := &pathEntry.Paths[0].Files[0]
	println("some model: ", ref1.FullPath)
	println("parent:", ref1.Parent.FauxPath)
	fmt.Printf("looking in \"%s\" for files...\n", ref1.Parent.Base())
	for _, x := range ref1.Parent.Files {
		println("  -->", x.Path)
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
