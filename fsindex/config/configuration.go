package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/tfwio/sekhem/fsindex"
	"github.com/tfwio/sekhem/util"
	"github.com/tfwio/sekhem/util/pandoc"
)

// Configuration is for JSON i/o.
type Configuration struct {
	Server       `json:"serv"`
	Pandoc       pandoc.Settings    `json:"pandoc"`
	Root         RootConfig         `json:"root"`
	Locations    []StaticPath       `json:"stat,omitempty"`
	Indexes      []IndexPath        `json:"indx,omitempty"`
	Extensions   []fsindex.FileSpec `json:"spec,omitempty"`
	IndexCfg     fsindex.Settings   `json:"index-cfg,omitempty"`
	Database     string             `json:"db.src,omitempty"` // relative to the configuration file
	DatabaseType string             `json:"db.sys,omitempty"`
}

// SessionHost gets a simple string that is used in our sessions db.
func (c *Configuration) SessionHost(input string) string {
	return fmt.Sprintf("%s%s", input, strings.TrimLeft(c.Port, ":"))
}

// GetFilePath only checks to see if we have indexed (configuration.Index)
// the path in order to obtain/return it.
func (c *Configuration) GetFilePath(route string, action string) (string, error) {
	urlpath := strings.Replace(action, "/tag/", "", 1)
	result := ""
	for _, index := range c.Indexes {
		if route == strings.Trim(index.Target, "/") {
			result += fmt.Sprintf("%s == %s\n", route, strings.Trim(index.Target, "/"))
			return fmt.Sprintf("%s%s", util.UnixSlash(filepath.Dir(index.Source)), urlpath), nil
		}
	}
	return result, errors.New("file not found")
}

// ToJSON writes JSON configuration.
func (c *Configuration) ToJSON(jsonPath string) {
	containingDirectory := filepath.Dir(jsonPath)
	if !util.DirectoryExists(containingDirectory) {
		panic(fmt.Sprintf("Target directory does not exist: %s\n", containingDirectory))
	}
	c.SaveJSON(jsonPath)
	fmt.Printf(constMessageWroteJSON, jsonPath)
}

// FromJSON loads JSON configuration.
func (c *Configuration) FromJSON(json string) {
	if !util.FileExists(json) {
		c.ToJSON(json)
		println("[terminating application]")
		os.Exit(1)
	} else {
		c.LoadJSON(json)
	}
	c.MapExtensions()
}

// HasTLS checks if we have certificate files pointed to by our configuration file.
//
// It does not validate the certificates.
func (c *Configuration) HasTLS() bool {
	return c.Server.hasCert() && c.Server.hasKey()
}

// DoTLS returns a boolean value that tells wether or not
// the client and configuration is going to serve via TLS.
func (c *Configuration) DoTLS() bool {
	if UseTLS {
		return c.HasTLS()
	} else if c.TLS {
		return c.HasTLS()
	}
	return false
}

// DefaultFile provides a absolute file-system path to the
// default (e.g. "<http|https>://<host>:<port>/index.html") file that is served.
func (c *Configuration) DefaultFile() string {
	return util.Abs(util.Cat(c.Root.Directory, "\\", c.Root.Default))
}

// GetBasePath returns the server's base path; e.g. `"http://localhost:5500"`
func (c *Configuration) GetBasePath() string {
	return fmt.Sprintf(`%s://%s%s`, util.IIF(c.Server.TLS, "https", "http"), c.Server.Host, c.Port)
}

// GetPath appends `more` to the default path (see `GetBasePath`).
// e.g. `"<http|https>://<host>:<port>/<...more>"`
func (c *Configuration) GetPath(more ...string) string {
	return util.Cat(c.GetBasePath(), "/", util.TrimJoin("/", more...))
}

// InitializeDefaults produces faux configuration settings.
func (c *Configuration) InitializeDefaults(path string, targetPath string) {
	c.IndexCfg = fsindex.DefaultSettings

	// initialize model array and map.
	models = []fsindex.Model{}
	mdlMap = make(map[string]*fsindex.Model)

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
			Source: constImagesSourceDefault,
			Target: constImagesTargetDefault,
		},
		StaticPath{
			Source:    constStaticSourceDefault,
			Target:    constStaticTargetDefault,
			Browsable: true,
		},
	}
	c.Indexes = []IndexPath{
		// FIXME: this particular path is to be associated with indexing.
		IndexPath{
			Source:      path,
			Target:      util.Wrapper("/", targetPath),
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
	c.Pandoc = pandoc.Settings{
		Executable:   "data/pandoc/pandoc.exe",
		HTMLTemplate: "data/pandoc/md-tpl.htm",
		MetaTemplate: "data/pandoc/md-meta.htm",
		Flags:        "-N", // numbered headers
		Extensions:   "+abbreviations+auto_identifiers+autolink_bare_uris+backtick_code_blocks+bracketed_spans+definition_lists+emoji+escaped_line_breaks+example_lists+fancy_lists+fenced_code_attributes+fenced_divs+footnotes+header_attributes+inline_code_attributes+implicit_figures+implicit_header_references+inline_notes+link_attributes+mmd_title_block+multiline_tables+raw_tex+simple_tables+smart+startnum+strikeout+table_captions+yaml_metadata_block",
	}
	// c.info()
	c.Prepare()
}

// Prepare sorts out ugly configuration settings.
func (c *Configuration) Prepare() {
	for _, l := range c.Locations {
		l.Target = util.WReap("/", l.Target)
	}
	for _, i := range c.Indexes {
		i.Target = util.WReap("/", i.Target)
	}
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
	if JSON, E := json.MarshalIndent(c, "", "\t"); E == nil {
		ioutil.WriteFile(path, JSON, constIOAccess)
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

const constIOAccess = 0600 // 0777
