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
	// UseTLS is a cli (console) override for using TLS or not.
	// It can only be put on in the console and if `--use-tls` is not set
	// then this will always remain false.
	UseTLS = false
	// OverrideHost overrides the default `Configuration.Server.Host`  with `UseHost` if true.
	OverrideHost = false
	// UseHost is identifies the host to use during `OverrideHost` use.
	UseHost = ""
	// DefaultConfigFile — you know.  Default = `./data/conf`.
	DefaultConfigFile = "data/conf.json"
	extMap            map[string]*fsindex.FileSpec
	mdlMap            map[string]*fsindex.Model
	models            []fsindex.Model
	xCounter          int32
	fCounter          int32
)

// JSONIndex — a simple container for JSON.
type JSONIndex struct {
	Index []string `json:"index"`
}

// GinConfig configures gin.Engine.
func (c *Configuration) GinConfig(router *gin.Engine) {
	c.GinConfigure(true, router)
}

// GinConfigure configures gin.Engine.
// if justIndex is set to true, we just rebuild our indexes.
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
	c.initializeModels()
	c.serveModelIndex(router)
}

func (c *Configuration) initializeModels() {
	for _, path := range c.Indexes {
		c.initializeModel(&path)
	}
}
func (c *Configuration) getModelIndex(mdl *fsindex.Model) (int, bool) {
	for index, mMdl := range models {
		if mMdl.FullPath == mdl.FullPath {
			return index, true
		}
	}
	return -1, false
}
func (c *Configuration) initializeModel(path *IndexPath) {

	fmt.Printf("--> indexing: %s\n", path.Target)
	model := c.createEntry(*path, c.IndexCfg)
	if _, ok := mdlMap[util.AbsBase(path.Source)]; !ok {
		models = append(models, model)
	} else {
		if index, ok := c.getModelIndex(&model); ok {
			models[index] = model
			println("Injecting memory-Model %s at index %d", mdlMap[model.Name].Name, index)
		}
	}
	if index, ok := c.getModelIndex(&model); ok {
		mdlMap[model.Name] = &models[index]
	} else {
		panic("Could not find memory-Model")
	}
}

func (c *Configuration) getSimpleIndexTarget(path *IndexPath) string {
	return util.WReap("/", util.AbsBase(path.Source))
}

func (c *Configuration) getIndexTarget(path *IndexPath) string {
	modelpath := util.WReap("/", path.Target)
	if !c.IndexCfg.OmitRootNameFromPath {
		modelpath = util.WReap("/", path.Target, util.AbsBase(path.Source))
	}
	return modelpath
}
func (c *Configuration) hasModel(route string) bool {

	if _, ok := mdlMap[route]; ok {
		return true
	}
	return false
}
func (c *Configuration) indexFromTarget(route string) *IndexPath {
	inputTarget := util.WReap("/", route)
	for _, x := range c.Indexes {
		simpleIndexTarget := c.getSimpleIndexTarget(&x)
		if inputTarget == simpleIndexTarget {
			return &x
		}
	}
	return nil
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

func (c *Configuration) createEntry(path IndexPath, settings fsindex.Settings) fsindex.Model {

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
