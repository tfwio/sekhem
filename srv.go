package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/urfave/cli"

	"tfw.io/Go/fsindex/fsindex"
	"tfw.io/Go/fsindex/fsindex/config"
	"tfw.io/Go/fsindex/util"
)

// Configuration variables

var (
	configuration config.Configuration
	mCli          cli.App

	serverRoot string // e.g. https://tfw.io:5500
	indexPath  string // e.g. https://tfw.io:5500/[path] (or <serverRoot>/<path>)

	pathEntry   fsindex.PathEntry
	pathEntries []fsindex.PathEntry

	xCounter  int32
	fCounter  int32
	xpCounter *int32
)

// SimpleModel collects our indexes
type SimpleModel struct {
	File     map[string]*fsindex.FileEntry
	FileSHA1 map[string]*fsindex.FileEntry
	Path     map[string]*fsindex.PathEntry
	PathSHA1 map[string]*fsindex.PathEntry
}

func (m *SimpleModel) create() {
	m.File = make(map[string]*fsindex.FileEntry)
	m.FileSHA1 = make(map[string]*fsindex.FileEntry)
	m.Path = make(map[string]*fsindex.PathEntry)
	m.PathSHA1 = make(map[string]*fsindex.PathEntry)
}

func loadConfig() {

	if !util.FileExists(config.DefaultConfigFile) {
		if !util.DirectoryExists(constDefaultDataPath) {
			os.Mkdir(constDefaultDataPath, 0777)
		}
		configuration.SaveJSON(config.DefaultConfigFile)
		println(constMessageWroteJSON)
		os.Exit(1)
	} else {
		configuration.LoadJSON(config.DefaultConfigFile)
	}

	configuration.MapExtensions()

}

// FIXME: `pathIndex[0]` is used (solely).
func configure(pathIndex ...string) {
	configuration.InitializeDefaults()

	// some default values
	configuration.Server.TLS = len(os.Args) == 2 && os.Args[1] == "tls"

	serverRoot = fmt.Sprintf(`%s://%s%s`, util.IIF(configuration.Server.TLS, "https", "http"), constServerDefaultPort, "tfw.io")
	// indexPath = fmt.Sprintf(`%s/%s`, serverRoot, "v")
	indexPath = fmt.Sprintf(`%s/%s`, serverRoot, "v")

	loadConfig() // loads (or creates conf.json and terminates application)
	indexPath = fmt.Sprintf(`%s://%s%s/%s`, util.IIF(configuration.Server.TLS, "https", "http"), configuration.Server.Host, configuration.Server.Port, configuration.Server.Path)

	// TODO: remove this
	println("- path for indexed files: ", indexPath)
	pathEntry = fsindex.PathEntry{
		PathSpec: fsindex.PathSpec{
			FileEntry: fsindex.FileEntry{
				Parent:   nil,
				Name:     util.AbsBase(pathIndex[0]),
				FullPath: util.Abs(pathIndex[0]),
				SHA1:     util.Sha1String(pathIndex[0]),
			},
			IsRoot: true},
		FauxPath:    indexPath,
		FileFilter:  configuration.Extensions,
		IgnorePaths: []string{},
	}
	println("- check name:", pathEntry.Name)
	println("- check sha1:", pathEntry.SHA1)
	println("- check path:", pathEntry.FauxPath)

	// TODO: check each Source path in configuration.Indexes
	pathEntries = make([]fsindex.PathEntry, len(configuration.Indexes))
	for i, p := range configuration.Indexes {
		pathEntries[i].PathSpec = fsindex.PathSpec{
			FileEntry: fsindex.FileEntry{
				Parent:   nil,
				Name:     util.AbsBase(util.Abs(p.Source)),
				FullPath: util.Abs(p.Source),
				SHA1:     util.Sha1String(util.Abs(p.Source)),
			},
			IsRoot: true}
		pathEntries[i].FileFilter = configuration.GetFilter(p.Extensions)
		pathEntries[i].FauxPath = util.Cat(serverRoot, p.Alias)
		pathEntries[i].IgnorePaths = p.IgnorePaths
		print(fmt.Sprintf("- JSON index: %s\n", util.AbsBase(util.Abs(p.Source))))
		for _, x := range pathEntries[i].FileFilter {
			print(fmt.Sprintf("  - got extension: %s\n", x.Name))
		}
	}
}

func main() {
	initializeCli()
}

func serveJSONPathEntry(pContext *gin.Context) {
	pContext.JSON(http.StatusOK, &pathEntry)
}

func makeMdl() SimpleModel {
	var s SimpleModel
	if s.File != nil {
		for k := range s.File {
			delete(s.File, k)
		}
	}
	if s.FileSHA1 != nil {
		for k := range s.FileSHA1 {
			delete(s.FileSHA1, k)
		}
	}
	if s.Path != nil {
		for k := range s.Path {
			delete(s.Path, k)
		}
	}
	if s.PathSHA1 != nil {
		for k := range s.PathSHA1 {
			delete(s.PathSHA1, k)
		}
	}
	s.create()
	return s
}

// AddPath is a callback per PathEntry.
// It adds each PathEntry to a flat (non-hierarchical) map (dictionary).
func (m *SimpleModel) AddPath(p *fsindex.PathEntry, c *fsindex.PathEntry) {
	m.Path[c.Rooted(p)] = c
	m.PathSHA1[c.SHA1] = c
}

// AddFile is a callback per FileEntry.
// It adds each FileEntry to a flat (non-hierarchical) map (dictionary).
func (m *SimpleModel) AddFile(p *fsindex.PathEntry, c *fsindex.FileEntry) {
	m.File[c.Rooted(p)] = c
	m.FileSHA1[c.SHA1] = c
}

func createIndex() {

	xCounter, fCounter = 0, 0

	mdl := makeMdl()

	handler := fsindex.Handlers{
		ChildPath: func(root *fsindex.PathEntry, child *fsindex.PathEntry) bool {
			mdl.AddPath(root, child)
			return false
		},
		ChildFile: func(root *fsindex.PathEntry, child *fsindex.FileEntry) bool {
			ext := strings.ToLower(filepath.Ext(child.FullPath))
			if ext == ".md" {
				// datestring := checkDateString(child.Base())
				mdl.AddFile(root, child)
				// println(fmt.Sprintf("  - f%-4d, %s -- %s", fCounter, (*child).Path, datestring))
				// return true
			}
			fCounter++
			return false
		},
	}
	pathEntry.Refresh(nil, &xCounter, &handler)

	checkSimpleModel(&mdl)
}

func checkSimpleModel(mdl *SimpleModel) {
	// map counters don't yield adequate length
	println("File map Count: ", len(mdl.File))
	println("Path map Count: ", len(mdl.Path))
	//
	println("File Count: ", fCounter)
	println("Path Count: ", xCounter)
	// println("some model: ", (*mdl.Path[`.mmd\THIRD PARTY\relisoft-windows-api-tut\12 olerant.md`]).FullPath)
	ref1 := &pathEntry.Paths[0].Files[0]
	println("some model: ", (*ref1).FullPath)
	println("parent:", ref1.Parent.FauxPath)
	// mf := mdl.File[`.mmd\THIRD PARTY\relisoft-windows-api-tut\12 olerant.md`]
	fmt.Printf("looking in \"%s\" for files...\n", ref1.Parent.Base())
	for _, x := range ref1.Parent.Files {
		println("  -->", x.Path)
	}
}

func initializeCli() {
	mCli.Name = filepath.Base(os.Args[0])
	mCli.Authors = []cli.Author{cli.Author{Name: "tfw; et alia" /*, Email: "tfwroble@gmail.com"}, cli.Author{Name: "Et al."*/}}
	mCli.Version = "v0.0.0"
	mCli.Copyright = "tfwio.github.com/go-fsindex\n\n   This is free, open-source software.\n   disclaimer: use at own risk."
	mCli.Action = func(*cli.Context) { initializeApp() }
	mCli.Commands = []cli.Command{cli.Command{
		Name:        "run",
		Action:      func(*cli.Context) { initializeApp() },
		Usage:       "Runs the server.",
		Description: "Default operation.",
		Aliases:     []string{"go"},
		Flags:       []cli.Flag{},
	}}
	mCli.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "tls",
			Destination: &config.UseTLS,
			Usage:       "Wether or not to use TLS.\n\tNote: if set, overrides (JSON) conf settings.",
		},
		cli.StringFlag{
			Name:        "conf",
			Usage:       "Points to a custom configuration file.",
			Value:       config.DefaultConfigFile,
			Destination: &config.DefaultConfigFile,
		},
	}
	mCli.Run(os.Args)
}

func initializeApp() {

	gin.SetMode(gin.ReleaseMode)
	// should be using

	println("==> configure")
	configure(`C:\Users\tfwro\Desktop\DesktopMess\ytdl_util-0.1.2.1-dotnet-client35-anycpu-win64\downloads`)

	println("==> server setup")
	mGin := gin.Default()

	// these files are all stored in the public directory.
	// they are the only files we're serving specifically in
	// that directory.

	DefaultFile := util.Abs(util.Cat(configuration.Root.Directory, "\\", configuration.Root.Default))
	// configuration.configInfo()

	// println("==> server setup: configuration.Root.AliasDefault")
	println("alias-default")
	for _, rootEntry := range configuration.Root.AliasDefault {
		target := util.Cat(configuration.Root.Path, rootEntry)
		mGin.StaticFile(target, DefaultFile)
		fmt.Printf("  ? Target = %s, Source = %s\n", target, configuration.DefaultFile())
	}
	println("root-files")
	for _, rootEntry := range configuration.Root.Files {
		target := util.Cat(configuration.Root.Path, rootEntry)
		source := util.Abs(util.Cat(configuration.Root.Directory, "\\", rootEntry))
		mGin.StaticFile(target, source)
		fmt.Printf("  > Target = %s, Source = %s\n", target, source)
	}
	println("locations")
	for _, tgt := range configuration.Locations {
		println("- serving path:", tgt.Target, "from", util.Abs(tgt.Source))
		mGin.StaticFS(tgt.Target, gin.Dir(util.Abs(tgt.Source), tgt.Browsable))
		fmt.Printf("  > Target = %s, Source = %s\n", tgt.Target, tgt.Source)
	}
	fmt.Printf("- default: Target = %s, Source =  %s\n", configuration.Root.Path, DefaultFile)
	mGin.StaticFile(constRootPathDefault, DefaultFile)

	mGin.GET("/json/", serveJSONPathEntry)

	createIndex()

	if configuration.DoTLS() {
		println("- use tls")
		if err := mGin.RunTLS(configuration.Server.Port, configuration.Server.Crt, configuration.Server.Key); err != nil {
			panic(fmt.Sprintf("router error: %s\n", err))
		}
	} else {
		println("- no tls")
		if err := mGin.Run(configuration.Server.Port); err != nil {
			panic(fmt.Sprintf("router error: %s\n", err))
		}
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
