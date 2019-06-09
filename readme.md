This is a little golang web-server app for servicing NPM sandboxes
for react or the like.

It generally creates a web-server and will likely use statik for
compiling once you might be happy with a simple react app distro.

The primary purpose of the app is to create a file-system index
wrapper that can serve JSON file-system indexes.

This was written to support go-mmd-fs or whatever its called at this point.

Configuration
----------------

You may setup as many paths as you like to be served however
they must not conflict with one another!

All of our configuration takes place in a configure function...

```go
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
		AliasDefault: []string{"home", "index.htm", "index.html", "index", "default", "default.htm"},
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
```

Root
-------------

The way this works is explicitly telling the server what files to serve in the root.
For this we use something like: `[]string{"public\\index.html", "public\\bundle.js", "public\\bundle.js"}`

