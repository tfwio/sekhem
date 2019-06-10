package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"tfw.io/Go/fsindex/fsindex"
	"tfw.io/Go/fsindex/util"
)

// Configuration is for JSON i/o.
type Configuration struct {
	Server     Server             `json:"serv"`
	Root       RootConfig         `json:"root"`
	Locations  []StaticPath       `json:"stat"`
	Indexes    []IndexPath        `json:"indx,omitempty"`
	Extensions []fsindex.FileSpec `json:"spec,omitempty"`
	indexPath  string
}

// LoadConfig loads JSON configuration.
func (c *Configuration) LoadConfig() {
	if !util.FileExists(DefaultConfigFile) {
		if !util.DirectoryExists(constDefaultDataPath) {
			os.Mkdir(constDefaultDataPath, 0777)
		}
		c.SaveJSON(DefaultConfigFile)
		println(constMessageWroteJSON)
		os.Exit(1)
	} else {
		c.LoadJSON(DefaultConfigFile)
	}

	c.MapExtensions()
}

// DoTLS returns a boolean value that tells wether or not
// the client and configuration is going to serve via TLS.
func (c *Configuration) DoTLS() bool {
	if UseTLS {
		return c.Server.hasCert() && c.Server.hasKey()
	}
	return c.Server.hasCert() && c.Server.hasKey() && c.Server.TLS
}

// DefaultFile provides a absolute file-system path to the
// default (e.g. "http://host:80/index.html") file that is served.
func (c *Configuration) DefaultFile() string {
	return util.Abs(util.Cat(c.Root.Directory, "\\", c.Root.Default))
}

func (c *Configuration) getBasePath() string {
	return fmt.Sprintf(`%s://%s%s`, util.IIF(c.Server.TLS, "https", "http"), c.Server.Host, c.Server.Port)
}

func (c *Configuration) getPath(more ...string) string {
	return util.Cat(c.getBasePath(), "/", strings.Join(more, "/"))
}

// InitializeDefaults produces faux configuration settings.
func (c *Configuration) InitializeDefaults() {
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
	c.indexPath = c.getPath(c.Server.Path)
}

// GetFilter uses the dictionary map to return the applicable `[]FileSpec`.
func (c *Configuration) GetFilter(extensions []string) []fsindex.FileSpec {
	var result []fsindex.FileSpec
	for _, ext := range extensions {
		if x, o := extMap[ext]; o {
			result = append(result, *x)
		}
	}
	return result
}

func (c *Configuration) configInfo() {
	println("Configuration Information: Root")
	println("=================================")
	c.Server.info()
	c.Root.info()
	for _, loc := range c.Locations {
		println("==> Location")
		println(fmt.Sprintf("----> Source    = %s", loc.Source))
		println(fmt.Sprintf("----> Browsable = %v", loc.Browsable))
		println(fmt.Sprintf("----> Target    = %s", loc.Target))
	}
}

// MapExtensions builds a dictionary map of extension info.
func (c *Configuration) MapExtensions() {
	extMap = make(map[string]*fsindex.FileSpec)
	for i, x := range c.Extensions {
		extMap[x.Name] = &(c.Extensions[i])
	}
}

// SaveJSON writes JSON to `path`.
func (c *Configuration) SaveJSON(path string) {
	println(fmt.Sprintf("==> Marshal JSON %s", path))
	if JSON, E := json.Marshal(c); E == nil {
		ioutil.WriteFile(path, JSON, 0777)
	} else {
		panic(E)
	}
}

// LoadJSON reads JSON from `path`.
func (c *Configuration) LoadJSON(path string) {
	println(fmt.Sprintf("==> Unmarshal JSON %s", path))
	data := util.CacheBytes(path)
	if E := json.Unmarshal(data, c); E != nil {
		panic(E)
	}
}
