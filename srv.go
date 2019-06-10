package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/urfave/cli"

	"tfw.io/Go/fsindex/fsindex"
	"tfw.io/Go/fsindex/util"
)

// Configuration variables

var (
	defaultConfigFile = "data/conf.json"
	useTLS            = false
	configuration     ConfigFile
	extMap            map[string]*fsindex.FileSpec
)
var (
	mCli cli.App

	serverRoot string // e.g. https://tfw.io:5500
	indexPath  string // e.g. https://tfw.io:5500/[path] (or <serverRoot>/<path>)

	pathEntry   fsindex.PathEntry
	pathEntries []fsindex.PathEntry

	xCounter  int32
	fCounter  int32
	xpCounter *int32
)

// Server info for JSON i/o.
type Server struct {
	Host string `json:"host"`
	Port string `json:"port"`
	TLS  bool   `json:"tls"` // default=false unless `os.Args[1] == "tls"` or specified in `[data/]config.json`.
	Key  string `json:"key,omitempty"`
	Crt  string `json:"crt,omitempty"`
	Path string `json:"path"`
}

func (s *Server) info() {
	println("> Server")
	println(fmt.Sprintf("--> Host = %s", s.Host))
	println(fmt.Sprintf("--> Port = %s", s.Port))
	println(fmt.Sprintf("--> TLS  = %v", s.TLS))
	println(fmt.Sprintf("--> Key  = %s", s.Key))
	println(fmt.Sprintf("--> Crt  = %s", s.Crt))
	println(fmt.Sprintf("--> Path  = %s", s.Path))
}
func (s *Server) hasKey() bool {
	return util.FileExists(constServerTLSKeyDefault)
}
func (s *Server) hasCert() bool {
	return util.FileExists(constServerTLSCertDefault)
}

func (s *Server) initServerConfig() {
	s.Host = constServerDefaultHost
	s.Port = constServerDefaultPort
	s.TLS = len(os.Args) == 2 && os.Args[1] == "tls"
	s.Crt = constServerTLSCertDefault
	s.Key = constServerTLSKeyDefault
	s.Path = "v"
}

// StaticPath is a definition for directories we'll
// allow into the app, preferably by way of JSON config.
type StaticPath struct {
	Source string `json:"src"`
	Target string `json:"tgt"`
	// show directory file-listing in browser
	Browsable bool `json:"nav"`
}

// IndexPath — same as StaticPath, however we can call on something like
// `[target].json` to calculate/generate a file-index listing.
type IndexPath struct {
	Alias       string   `json:"alias,omitempty"` // this isn't supported really, but the idea is to use this as the name of our target as opposed to our root-directory name, or default server path (in Server.Path).
	Source      string   `json:"src"`
	Target      string   `json:"tgt"`
	Browsable   bool     `json:"nav"` // show directory file-listing in browser
	Servable    bool     `json:"serve"`
	IgnorePaths []string `json:"ignorePaths"` // absolute paths to ignore
	Extensions  []string `json:"spec"`        // file extensions to recognize; I.E.: the `ConfigFile.Extensions` .Name.
	path        string   // path as used in memory; we'll probably just ignore this guy.
}

// RootConfig is used to tell the server what files are to
// be served in the root directory.
type RootConfig struct {
	Path         string   `json:"path"`
	Directory    string   `json:"dir"`
	Files        []string `json:"files"`
	AliasDefault []string `json:"alias"`
	Default      string   `json:"default"`
}

func (r *RootConfig) info() {
	println("> Root")
	println(fmt.Sprintf("--> Path         = %s", r.Path))
	println(fmt.Sprintf("--> Directory    = %s", r.Directory))
	println(fmt.Sprintf("--> Files        = %s", r.Files))
	println(fmt.Sprintf("--> AliasDefault = %s", r.AliasDefault))
	println(fmt.Sprintf("--> Default      = %s", r.Default))
}

// ConfigFile is for JSON i/o.
type ConfigFile struct {
	Server     Server             `json:"serv"`
	Root       RootConfig         `json:"root"`
	Locations  []StaticPath       `json:"stat"`
	Indexes    []IndexPath        `json:"indx,omitempty"`
	Extensions []fsindex.FileSpec `json:"spec,omitempty"`
}

func (c *ConfigFile) doTLS() bool {
	if useTLS {
		return c.Server.hasCert() && c.Server.hasKey()
	}
	return c.Server.hasCert() && c.Server.hasKey() && c.Server.TLS
}

func (c *ConfigFile) initializeDefaults() {
	// println("==> Configuring")
	c.Server.initServerConfig()
	c.Root = RootConfig{
		Path:         constRootPathDefault,
		Directory:    constRootDirectoryDefault,
		Files:        strings.Split(constRootFilesDefault, ","),
		Default:      constRootIndexDefault,
		AliasDefault: strings.Split(constRootAliasDefault, ","),
	}
	c.Locations = []StaticPath{
		StaticPath{
			Source:    constImagesSourceDefault,
			Target:    constImagesTargetDefault,
			Browsable: true,
		},
		StaticPath{
			Source:    constStaticSourceDefault,
			Target:    constStaticTargetDefault,
			Browsable: true,
		},
		// FIXME: this particular path is to be associated with indexing.
		StaticPath{
			Source:    "C:\\Users\\tfwro\\Desktop\\DesktopMess\\ytdl_util-0.1.2.1-dotnet-client35-anycpu-win64\\downloads",
			Target:    "/v/",
			Browsable: true,
		},
	}
	c.Indexes = []IndexPath{
		// FIXME: this particular path is to be associated with indexing.
		IndexPath{
			Source:      "C:\\Users\\tfwro\\Desktop\\DesktopMess\\ytdl_util-0.1.2.1-dotnet-client35-anycpu-win64\\downloads",
			Target:      "/v/",
			Browsable:   true,
			Servable:    true,
			Extensions:  []string{"Media"},
			IgnorePaths: []string{},
			path:        "",
		},
	}
	c.Extensions = []fsindex.FileSpec{
		fsindex.FileSpec{
			Name:       "Media",
			Extensions: strings.Split(constExtDefaultMedia, ","),
		},
		fsindex.FileSpec{
			Name:       "Markdown",
			Extensions: strings.Split(constExtDefaultMMD, ","),
		},
	}
	// c.configInfo()
}

func (c *ConfigFile) getFilter(extensions []string) []fsindex.FileSpec {
	var result []fsindex.FileSpec
	for _, ext := range extensions {
		if x, o := extMap[ext]; o {
			result = append(result, *x)
		}
	}
	return result
}

func (c *ConfigFile) configInfo() {
	println("Configuration Information: Root")
	println("=================================")
	c.Server.info()
	c.Root.info()
	for _, loc := range configuration.Locations {
		println("==> Location")
		println(fmt.Sprintf("----> Source    = %s", loc.Source))
		println(fmt.Sprintf("----> Browsable = %v", loc.Browsable))
		println(fmt.Sprintf("----> Target    = %s", loc.Target))
	}
}

func (c *ConfigFile) mapExtensions() {
	extMap = make(map[string]*fsindex.FileSpec)
	for i, x := range c.Extensions {
		extMap[x.Name] = &(c.Extensions[i])
	}
}

func (c *ConfigFile) write(path string) {
	println(fmt.Sprintf("==> Marshal JSON %s", path))
	if JSON, E := json.Marshal(c); E == nil {
		ioutil.WriteFile(path, JSON, 0777)
	} else {
		panic(E)
	}
}

func (c *ConfigFile) read(path string) {
	println(fmt.Sprintf("==> Unmarshal JSON %s", path))
	data := util.CacheBytes(path)
	if E := json.Unmarshal(data, c); E != nil {
		panic(E)
	}
}

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

	if !util.FileExists(defaultConfigFile) {
		if !util.DirectoryExists(constDefaultDataPath) {
			os.Mkdir(constDefaultDataPath, 0777)
		}
		configuration.write(defaultConfigFile)
		println(constMessageWroteJSON)
		os.Exit(1)
	} else {
		configuration.read(defaultConfigFile)
	}

	configuration.mapExtensions()

}

// FIXME: `pathIndex[0]` is used (solely).
func configure(pathIndex ...string) {
	configuration.initializeDefaults()

	// some default values
	configuration.Server.TLS = len(os.Args) == 2 && os.Args[1] == "tls"

	serverRoot = fmt.Sprintf(`%s://%s%s`, util.IIF(configuration.Server.TLS, "https", "http"), constServerDefaultPort, "tfw.io")
	indexPath = fmt.Sprintf(`%s/%s`, serverRoot, "v")

	loadConfig() // loads (or creates conf.json and terminates application)
	indexPath = fmt.Sprintf(`%s://%s%s/%s`, util.IIF(configuration.Server.TLS, "https", "http"), configuration.Server.Port, configuration.Server.Host, configuration.Server.Path)

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
		pathEntries[i].FileFilter = configuration.getFilter(p.Extensions)
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

func loadModel() {

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
	pathEntry.Refresh1(nil, &xCounter, &handler)

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

// not used
func cliInfo() {
}

func initializeCli() {
	mCli.Name = filepath.Base(os.Args[0])
	mCli.Authors = []cli.Author{cli.Author{Name: "tfw, et alia" /*, Email: "tfwroble@gmail.com"}, cli.Author{Name: "Et al."*/}}
	mCli.Version = "0.0.a"
	mCli.Copyright = "2019 tfwio.github.com/go-fsindex — this software warrants no license.\n\tdisclaimer: use at own risk."
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
			Destination: &useTLS,
			Usage:       "Wether or not to use TLS.\n\tNote: if set, overrides (JSON) conf settings.",
		},
		cli.StringFlag{
			Name:        "conf",
			Usage:       "Points to a custom configuration file.",
			Value:       defaultConfigFile,
			Destination: &defaultConfigFile,
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

	// FIXME: what?
	defaultFile := util.Abs(util.Cat(configuration.Root.Directory, "\\", configuration.Root.Default))
	// configuration.configInfo()

	println("==> server setup: configuration.Root.AliasDefault")
	for _, rootEntry := range configuration.Root.AliasDefault {
		fmt.Printf("  > Target = %s\n", rootEntry)
		target := util.Cat(configuration.Root.Path, rootEntry)
		mGin.StaticFile(target, defaultFile)
	}
	for _, rootEntry := range configuration.Root.Files {
		target := util.Cat(configuration.Root.Path, rootEntry)
		source := util.Abs(util.Cat(configuration.Root.Directory, "\\", rootEntry))
		fmt.Printf("- serving file: \"%s\" from \"%s\"\n", target, source)
		mGin.StaticFile(target, source)
	}
	for _, tgt := range configuration.Locations {
		println("- serving path:", tgt.Target, "from", util.Abs(tgt.Source))
		mGin.StaticFS(tgt.Target, gin.Dir(util.Abs(tgt.Source), tgt.Browsable))
	}
	fmt.Printf("- default: \"%s\" from \"%s\"\n", configuration.Root.Path, defaultFile)
	mGin.StaticFile(constRootPathDefault, defaultFile)

	mGin.GET("/json/", serveJSONPathEntry)

	loadModel()

	println("running")

	if configuration.doTLS() {
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
	constMessageErrorWritingJSON = "Error: writing configuration. %s"
	constRootAliasDefault        = "home,index.htm,index.html,index,default,default.htm"
	constRootFilesDefault        = "json.json,bundle.js,favicon.ico"
	constRootIndexDefault        = "index.html"
	constRootDirectoryDefault    = ".\\public"
	constRootPathDefault         = "/"
	constStaticSourceDefault     = "public\\static"
	constStaticTargetDefault     = "/static/"
	constImagesSourceDefault     = "public\\images"
	constImagesTargetDefault     = "/images/"
	constExtDefaultMedia         = ".mp4,.m4a,.mp3"
	constExtDefaultMMD           = ".md,.mmd"
)
