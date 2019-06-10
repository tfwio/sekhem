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
	configFilePath = "data/conf.json"
	localMedia     = fsindex.FileSpec{
		Name:       "Media (images)",
		Extensions: []string{".mp4", ".m4a", ".mp3"},
	}
	localMarkdown = fsindex.FileSpec{
		Name:       "Markdown (hyper-text)",
		Extensions: []string{".md", ".mmd"},
	}
	// config template
	configuration = ConfigFile{
		Server: Server{
			Host: "tfw.io",
			Port: ":5500",
			TLS:  len(os.Args) == 2 && os.Args[1] == "tls",
			Crt:  util.IIF(util.FileExists(tlsCrt), tlsCrt, ""),
			Key:  util.IIF(util.FileExists(tlsKey), tlsKey, ""),
			Path: "v",
		},
		Root: RootConfig{
			Path:         "/",
			Directory:    ".\\public",
			Files:        []string{`json.json`, `bundle.js`, `favicon.ico`},
			Default:      "index.html",
			AliasDefault: []string{"home", "index.htm", "index.html", "index", "default", "default.htm"},
		},
		Locations: []StaticPath{
			StaticPath{
				Source:    "public\\images",
				Target:    "/images/",
				Browsable: true,
			},
			StaticPath{
				Source:    "public\\static",
				Target:    "/static/",
				Browsable: true,
			},
			// FIXME: this particular path is to be associated with indexing.
			StaticPath{
				Source:    "C:\\Users\\tfwro\\Desktop\\DesktopMess\\ytdl_util-0.1.2.1-dotnet-client35-anycpu-win64\\downloads",
				Target:    "/v/",
				Browsable: true,
			},
		},
		Extensions: []fsindex.FileSpec{
			fsindex.FileSpec{
				Name:       "Media",
				Extensions: []string{".mp4", ".m4a", ".mp3"},
			},
			fsindex.FileSpec{
				Name:       "Markdown",
				Extensions: []string{".md", ".mmd"},
			},
		},
	}
	locations  []StaticPath
	rootConfig RootConfig
)
var (
	mCli cli.App

	serveHost string
	servePort string
	serveProt = "http" // default="http" unless `os.Args[1]` is "tls".
	serveTLS  bool     // default=false unless `os.Args[1] == "tls"`
	tlsKey    string
	tlsCrt    string

	indexRoot string // sub-path
	indexPath string // basis for generated URLs for the indexer.

	pathEntry fsindex.PathEntry

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

// StaticPath is a definition for directories we'll
// allow into the app, preferably by way of JSON config.
type StaticPath struct {
	Source    string `json:"src"`
	Target    string `json:"tgt"`
	Browsable bool   `json:"nav"` // show directory file-listing in browser
}

// IndexPath — same as StaticPath, however we can call on something like
// `[target].json` to calculate/generate a file-index listing.
type IndexPath struct {
	Source      string   `json:"src"`
	Target      string   `json:"tgt"`
	Browsable   bool     `json:"nav"`            // show directory file-listing in browser
	IgnorePaths []string `json:"ignore-paths"`   // absolute paths to ignore
	Extensions  []string `json:"spec,omitempty"` // file extensions to recognize; I.E.: the `ConfigFile.Extensions` .Name.
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

// ConfigFile is for JSON i/o.
type ConfigFile struct {
	Server     Server             `json:"serv"`
	Root       RootConfig         `json:"root"`
	Locations  []StaticPath       `json:"stat"`
	Indexes    []IndexPath        `json:"indx,omitempty"`
	Extensions []fsindex.FileSpec `json:"spec,omitempty"`
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

func writeConfig(cfg *ConfigFile) error {
	if JSON, E := json.Marshal(cfg); E == nil {
		ioutil.WriteFile("data/conf.json", JSON, 0777)
		return nil
	} else {
		println(E)
		return E
	}
}
func readConfig(path string) (ConfigFile, error) {

	var cfg ConfigFile
	if E := json.Unmarshal([]byte(util.CacheFile(path)), &cfg); E == nil {
		return cfg, E
	}
	return cfg, nil
}

func refreshConfig(tls string, port string, host string, path string) {
	serveProt, servePort, serveHost, indexRoot = tls, port, host, path
	indexPath = fmt.Sprintf(`%s://%s%s/%s`, serveProt, serveHost, servePort, indexRoot)
}

// FIXME: `pathIndex[0]` is used (solely).
func configure(pathIndex ...string) {

	serveTLS = len(os.Args) == 2 && os.Args[1] == "tls"
	refreshConfig(util.IIF(serveTLS, "https", "http"), ":5500", "tfw.io", "v")

	println("- path for indexed files: ", indexPath)

	tlsCrt, tlsKey = "data\\ia.crt", "data\\ia.key"

	if !util.FileExists(configFilePath) {
		if !util.DirectoryExists("./data") {
			os.Mkdir("data", 0777)
		}
		if E := writeConfig(&configuration); E != nil {
			panic(fmt.Sprintf("Error: writing configuration. %s", E))
		} else {
			panic("- wrote default data/conf.json\n  for your editing.")
		}
	} else {
		if conf, R := readConfig(configFilePath); R == nil {
			configuration = conf
			println("- got data/conf.json")
		} else {
			panic(fmt.Sprintf("Error: reading configuration. %s", R))
		}
	}

	pathEntry = fsindex.PathEntry{
		PathSpec: fsindex.PathSpec{
			FileEntry: fsindex.FileEntry{
				Parent:   nil,
				Name:     filepath.Base(pathIndex[0]),
				FullPath: util.Abs(pathIndex[0]),
				SHA1:     util.Sha1String(pathIndex[0]),
			},
			IsRoot: true},
		FauxPath:    indexPath,
		FileFilter:  []fsindex.FileSpec{localMedia, localMarkdown},
		IgnorePaths: []string{},
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

func (m *SimpleModel) AddPath(p *fsindex.PathEntry, c *fsindex.PathEntry) {
	m.Path[c.Rooted(p)] = c
	m.PathSHA1[c.Rooted(p)] = c
}

func (m *SimpleModel) AddFile(p *fsindex.PathEntry, c *fsindex.FileEntry) {
	m.File[c.Rooted(p)] = c
	m.FileSHA1[c.Rooted(p)] = c
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
	mCli.Name = "fsindex - indexer and simple utility web-server."
	mCli.UsageText = `this guy is primarily for use with NPM, React Webpack.`
	mCli.Authors = []cli.Author{cli.Author{Name: "tfw, et alia" /*, Email: "tfwroble@gmail.com"}, cli.Author{Name: "Et al."*/}}
	mCli.Description = `A simple little prototyping web-server.`
	mCli.Version = "0.0.0.A"
	mCli.Copyright = "2019 tfwio.github.com/go-fsindex"
}

func initializeCli() {
	mCli.Commands = []cli.Command{cli.Command{
		Name:    "run",
		Action:  func(*cli.Context) { initializeApp() },
		Usage:   "Deault operation — run the server.",
		Aliases: []string{"go"},
		Hidden:  true,
	}}
	mCli.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "conf",
			Value:       "data/conf.json",
			Destination: &configFilePath,
		},
	}
	if len(os.Args) == 1 {
		mCli.Run([]string{os.Args[0], "run"})
	} else {
		mCli.Run(os.Args)
	}
}

func initializeApp() {

	gin.SetMode(gin.ReleaseMode)
	// should be using
	configure(`C:\Users\tfwro\Desktop\DesktopMess\ytdl_util-0.1.2.1-dotnet-client35-anycpu-win64\downloads`)

	mGin := gin.Default()

	// these files are all stored in the public directory.
	// they are the only files we're serving specifically in
	// that directory.

	defaultFile := util.Abs(util.Cat(rootConfig.Directory, "\\", rootConfig.Default))
	for _, rootEntry := range rootConfig.AliasDefault {
		target := util.Cat(rootConfig.Path, rootEntry)
		fmt.Printf("- default alias: \"%s\" from \"%s\"\n", target, defaultFile)
		mGin.StaticFile(target, defaultFile)
	}
	for _, rootEntry := range rootConfig.Files {
		target := util.Cat(rootConfig.Path, rootEntry)
		source := util.Abs(util.Cat(rootConfig.Directory, "\\", rootEntry))
		fmt.Printf("- serving file: \"%s\" from \"%s\"\n", target, source)
		mGin.StaticFile(target, source)
	}
	for _, tgt := range locations {
		println("- serving path:", tgt.Target, "from", util.Abs(tgt.Source))
		mGin.StaticFS(tgt.Target, gin.Dir(util.Abs(tgt.Source), tgt.Browsable))
	}
	fmt.Printf("- default: \"%s\" from \"%s\"\n", rootConfig.Path, defaultFile)
	mGin.StaticFile("/", defaultFile)

	mGin.GET("/json/", serveJSONPathEntry)

	loadModel()

	println("running")

	if len(os.Args) == 2 && os.Args[1] == "tls" {
		println("- Using TLS")
		if err := mGin.RunTLS(servePort, tlsCrt, tlsKey); err != nil {
			fmt.Println("router error:", err)
		}
	} else {
		println("- Not using TLS")
		if err := mGin.Run(servePort); err != nil {
			fmt.Println("router error:", err)
		}
	}
}
