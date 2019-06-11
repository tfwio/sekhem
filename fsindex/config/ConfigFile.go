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
	Server      Server              `json:"serv"`
	Root        RootConfig          `json:"root"`
	Locations   []StaticPath        `json:"stat"`
	Indexes     []IndexPath         `json:"indx,omitempty"`
	Extensions  []fsindex.FileSpec  `json:"spec,omitempty"`
	PathEntries []fsindex.PathEntry `json:"-"`
	indexPath   string
	pathEntry   fsindex.PathEntry
}

// FromJSON loads JSON configuration.
func (c *Configuration) FromJSON() {
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

// GetBasePath returns the server's base path; e.g. `"http://localhost:5500"`
func (c *Configuration) GetBasePath() string {
	return fmt.Sprintf(`%s://%s%s`, util.IIF(c.Server.TLS, "https", "http"), c.Server.Host, c.Server.Port)
}

// GetPath appends `more` to the default path (see `GetBasePath`).
// e.g. `"http://localhost:5500/<...more>"`
func (c *Configuration) GetPath(more ...string) string {
	return util.Cat(c.GetBasePath(), "/", strings.Join(more, "/"))
}

// InitializeDefaults produces faux configuration settings.
func (c *Configuration) InitializeDefaults(path string, targetPath string) {
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
			Source:    util.Abs(constImagesSourceDefault),
			Target:    constImagesTargetDefault,
			Browsable: true,
		},
		StaticPath{
			Source:    util.Abs(constStaticSourceDefault),
			Target:    constStaticTargetDefault,
			Browsable: true,
		},
	}
	c.Indexes = []IndexPath{
		// FIXME: this particular path is to be associated with indexing.
		IndexPath{
			Source:      util.Abs(path),
			Target:      util.Wrap(targetPath, "/"),
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
	// c.info()
	c.indexPath = c.GetPath(c.Server.Path)
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

func (c *Configuration) info() {
	println("Configuration Information: Root")
	println("=================================")
	c.Server.info()
	c.Root.info()
	c.locationInfo()
}

func (c *Configuration) locationInfo() {
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
	if JSON, E := json.MarshalIndent(c, "", "  "); E == nil {
		ioutil.WriteFile(path, JSON, 0777)
	} else {
		panic(E)
	}
}

// LoadJSON reads JSON from `path`.
func (c *Configuration) LoadJSON(path string) {
	data := util.CacheBytes(path)
	if E := json.Unmarshal(data, c); E != nil {
		panic(E)
	}
}
