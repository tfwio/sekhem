[github.com/dhowden/tag]:         https://github.com/dhowden/tag "github.com/gogonic/gin"
[github.com/gin-gonic/gin]:       https://github.com/gin-gonic/gin
[github.com/jinzhu/gorm]:         https://github.com/jinzhu/gorm
[github.com/urfave/cli]:          https://github.com/urfave/cli
[github.com/mattn/go-sqlite3]:    https://github.com/mattn/go-sqlite3
[github.com/jgm/pandoc]:          https://github.com/jgm/pandoc
[github.com/jgm/pandoc/releases]: https://github.com/jgm/pandoc/releases
[pandoc]:                         https://pandoc.org


Written to be customized, a golang web-server app for servicing NPM
sandboxes for react and the like.

Sekhm
===================

A powerful, customizable sandbox/server which allows one to hone/nurture
their HTML/Javascript/UX/React designs and programming around some
XHR/JSON file-indexes.

> This project is fairly new, documentation and demos are forthcoming 😁

Features
-------------

CURRENT


- **JSON server configuration**.
- **XHR/JSON file-system index**
    - with GET|POST requests with manual refresh capability
    - file-extention and directory filters
    - XHR/JSON Tag requests (audio/video metadata) for multi-media (using [github.com/dhowden/tag])
- smart CLI interface for overriding config settings like
    - `--port <number>`
    - `--tls`: supply this flag to use tls when the config
      file has it off by default.
- More XHR request/data integrations to come perhaps including Calibre EBOOK
  data, plex (meta-info) and Chrome bookmarks/favicons, although
  is yet to be determined exactly how and when at this point.

IN PROGRESS

- Logon sessions (only sqlite3 data backend for now) are nearly complete.
- Separate demo sandbox projects (soon) 

KNOWN BUGS for expected fixes

- #1 **I'd like to see file time-stamps (CRD)**  
  This may only be implemented in windows since thats the main dev
  workstation.  PRs and discussions (bug section) are welcome.
- #2 [bug] *if file-extensions are poorly configured ATM:*  
  Multiple/duplicate files are returned in XHR/JSON due to extension definitions sharing the same extension.  
  (will be fixed soon)
- #3 [feature] **remove long-empty path entries**  
  The idea is to add a post indexing filter that strips out all empty directories and to provide a JSON option to apply such a filter.

Project Status: alpha (development) phase v0.0.0 has not changed yet for 90 revisions ATM.


### XHR: MMD2HTML (Pandoc)

Pandoc is required for XHR address `[http|https]://<host:port>/pan/:path/`

- [pandoc] (home)  
  [github.com/jgm/pandoc] on github  
  get it here: [github.com/jgm/pandoc/releases]

### conf.json (usage)

See: [data/doc/configuration](data/doc/configuration.md)

Develop
=======================

First steps

```bash
# go get it
go get github.com/tfwio/sekhm

# go to the project directory
pushd ${GOPATH}/src/github.com/tfwio/sekhm

# run the bootstrap script
./do bootstrap

# environment variable ${GOARCH} is targeted by default
# build it targeting your native OS.
./do build mod

# build it targeting native os amd64
./do build mod amd64

# build it targeting linux amd64 (from another architecture)
./do build mod linux amd64

# build it targeting 
./do build mod 386

# or manually typing everything out...
GO111MODULE=on GOARCH=windows go build -tags=jsoniter -o srv.exe -mod vendor srv.go
# 
```
*This hadn't been tested on \*nix or mac work-station(s) just yet.  
if you experience an OS related issue, PRs and bug-report/discussion are welcome*

The first time you run `./srv-${GOARCH}`, it will generate a example configuration file, `data/conf.json`, and exit.

Running the application once more will run the server on your localhost (http://localhost:5500/) and serve a demo html file that's sitting in the `./public` directory.

You can run the server on a differt port by calling

```bash
./srv-${GOARCH} --port <number>
```

Though there are the sub-commands `run` | `go` which do the same thing, you don't need to use them as its the default context — to run the server.

Running Help:
```bash
NAME:
   srv-amd64.exe

USAGE:
    [global options] command [command options] [arguments...]

VERSION:
   v0.0.0

AUTHOR:
   tfw; et alia

COMMANDS:
     run, go    Runs the server.
     make-conf  srv-amd64.exe make-conf <[file-path].json>
     help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --tls          Wether or not to use TLS.
   --port value   UseHost is identifies the host to use to over-ride JSON config. (default: 0)
   --conf value   Points to a custom configuration file. (default: "data/conf.json")
   --version, -v  print the version

COPYRIGHT:
   github.com/tfwio/sekhm

   This is free, open-source software.
   disclaimer: use at own risk.

```


Golang Development Libs Used
-----------------

*direct dependencies*

- [github.com/urfave/cli]
- [github.com/gin-gonic/gin]
- [github.com/mattn/go-sqlite3]
- [github.com/jinzhu/gorm]
- [github.com/dhowden/tag]

