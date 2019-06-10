package config

import "tfw.io/Go/fsindex/fsindex"

var (
	// UseTLS — you know.
	UseTLS = false
	// DefaultConfigFile — you know.
	DefaultConfigFile = "data/conf.json"
	extMap            map[string]*fsindex.FileSpec
)

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
