package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"

	"tfw.io/Go/fsindex/fsindex"
	"tfw.io/Go/fsindex/util"
)

const unknownString = "unknown date"

var (
	locations []StaticPath

	rootConfig RootConfig
	serveFiles []string

	serveHost string
	servePath string
	servePort string
	serveTLS  bool
	// (<serveTls>?"HTTPS:HTTP")://<serveHost>:<servePort>/<servePath>

	tlsKey string
	tlsCrt string

	pathEntry fsindex.PathEntry

	xCounter   int32
	fCounter   int32
	xpCounter  *int32
	localMedia = fsindex.FileSpec{
		Name: "Media (images)",
		Extensions: []string{
			".mp4",
			".m4a", // do these work on iphones/tablets? probably no.
			".mp3",
		},
	}
	localMarkdown = fsindex.FileSpec{
		Name: "Markdown (hyper-text)",
		Extensions: []string{
			".md",
			".mmd",
		},
	}
)

// StaticPath is a definition for directories we'll
// allow into the app, preferably by way of JSON config.
type StaticPath struct {
	Source    string
	Target    string
	Browsable bool
}

// RootConfig is used to tell the server what files are to
// be served in the root directory.
type RootConfig struct {
	Path         string
	Directory    string
	Files        []string
	AliasDefault []string
	Default      string
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

// FIXME: `pathIndex[0]` is used (solely).
func configure(pathIndex ...string) {

	serveFiles = []string{`index.html`, `json.json`, `bundle.js`, `favicon.ico`}
	serveTLS = len(os.Args) == 2 && os.Args[1] == "tls"
	servePort = ":5500"
	serveHost = "tfw.io"
	servePath = "v"
	serveProt := "http"

	if serveTLS == true {
		serveProt = "https"
	}

	faux := fmt.Sprintf(`%s://%s%s/%s`, serveProt, serveHost, servePort, servePath)
	println("- path for indexed files: ", faux)

	rootConfig = RootConfig{
		Path:         "/",
		Directory:    ".\\public",
		Files:        []string{`json.json`, `bundle.js`, `favicon.ico`},
		Default:      "index.html",
		AliasDefault: []string{"index.htm", "index.html", "index"},
	}

	tlsCrt = "data\\ia.crt"
	tlsKey = "data\\ia.key"

	locations = []StaticPath{
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
		StaticPath{
			Source:    "C:\\Users\\tfwro\\Desktop\\DesktopMess\\ytdl_util-0.1.2.1-dotnet-client35-anycpu-win64\\downloads",
			Target:    "/v/",
			Browsable: true,
		},
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
		FauxPath:    faux,
		FileFilter:  []fsindex.FileSpec{localMedia, localMarkdown},
		IgnorePaths: []string{},
	}
}

func main() {

	gin.SetMode(gin.ReleaseMode)

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

func serveJSONPathEntry(pContext *gin.Context) {
	pContext.JSON(http.StatusOK, &pathEntry)
}

func charIsNumber(input string) bool {
	for _, b := range []byte(input) {
		if !(b >= 48 && b <= 57) {
			return false
		}
	}
	return true
}

func checkDateString(input string) string {
	result := strings.Index(input, " ")
	// println(fmt.Sprintf("first-index:  %d", result))
	if result >= 0 && result == 8 && charIsNumber(input[:8]) {
		return input[:8]
	}
	return unknownString
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
