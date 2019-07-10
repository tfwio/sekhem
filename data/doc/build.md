[github.com/dhowden/tag]:         https://github.com/dhowden/tag "github.com/gogonic/gin"
[github.com/gin-gonic/gin]:       https://github.com/gin-gonic/gin
[github.com/jinzhu/gorm]:         https://github.com/jinzhu/gorm
[github.com/urfave/cli]:          https://github.com/urfave/cli
[github.com/mattn/go-sqlite3]:    https://github.com/mattn/go-sqlite3
[github.com/jgm/pandoc]:          https://github.com/jgm/pandoc
[github.com/jgm/pandoc/releases]: https://github.com/jgm/pandoc/releases
[home]: ../../readme.md "github.com/tfwio/sekhm/readme.md"
[features]: features.md
[configuration]: configuration.md
[build]: build.md
[usage]: usage.md
<!-- []:  -->

- [home]
    - [features]
    - [configuration]
    - [usage]
    - [build]


Building
=======================

go get it

```bash
go get github.com/tfwio/sekhm
```
go to the project directory
```bash
pushd ${GOPATH}/src/github.com/tfwio/sekhm
```
run the bootstrap script
```bash
./do bootstrap
```
environment variable ${GOARCH} is targeted by default  
build it targeting your native OS.
```bash
./do build mod
```
build it targeting native os amd64
```bash
./do build mod amd64
```
build it targeting linux amd64 (from another architecture)
```bash
./do build mod linux amd64
```
build it targeting 
```bash
./do build mod 386
```
or manually typing everything out...
```bash
GO111MODULE=on GOARCH=windows go build -tags=jsoniter -o srv.exe -mod vendor srv.go
```


*This hadn't been tested on \*nix or mac work-station(s) just yet.  
if you experience an OS related issue, PRs and bug-report/discussion are welcome*

Golang Development Libs Used
-----------------

*direct dependencies*

- [github.com/urfave/cli]
- [github.com/gin-gonic/gin]
- [github.com/mattn/go-sqlite3]
- [github.com/jinzhu/gorm]
- [github.com/dhowden/tag]

