[github.com/dhowden/tag]:         https://github.com/dhowden/tag "github.com/gogonic/gin"
[github.com/gin-gonic/gin]:       https://github.com/gin-gonic/gin
[github.com/jinzhu/gorm]:         https://github.com/jinzhu/gorm
[github.com/urfave/cli]:          https://github.com/urfave/cli
[github.com/mattn/go-sqlite3]:    https://github.com/mattn/go-sqlite3
[github.com/jgm/pandoc]:          https://github.com/jgm/pandoc
[github.com/jgm/pandoc/releases]: https://github.com/jgm/pandoc/releases
[github.com/json-iterator/go]:    https://github.com/json-iterator/go
[home]: ../../readme.md "github.com/tfwio/sekhem/readme.md"
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
  
*This hadn't been tested on \*nix or mac work-station(s) just yet.  
if you experience an OS related issue, PRs and bug-report/discussion are welcome*

go get it

```bash
go get github.com/tfwio/sekhem
```
go to the project directory
```bash
pushd ${GOPATH}/src/github.com/tfwio/sekhem
```
run the bootstrap script
```bash
./do bootstrap
```
environment variable ${GOARCH} is targeted by default  
build targeting your native OS.
```bash
./do build mod
```
[optional] build it targeting native os amd64
```bash
./do build mod amd64
```
[optional] build it targeting linux amd64 (from another architecture)
```bash
./do build mod linux amd64
```
[optional] build it targeting 
```bash
./do build mod 386
```
[optional] or manually typing everything out...
```bash
GO111MODULE=on GOOS=windows GOARCH=amd64 go build -tags 'jsoniter session' -o srv.exe -mod vendor srv.go
```

Just be sure that the vendor libs are downloaded.  If the build process
complains that you're missing dependencies, you'll have to `go get` any such lib.  
For example, to use gogonic/gin's `-tags=jsoniter` tag we'll provide the
[json-iterator][github.com/json-iterator/go] dependency.  
`go get github.com/json-iterator/go`

Golang Development Libs Used
-----------------

*direct dependencies*

- [github.com/urfave/cli]
- [github.com/gin-gonic/gin]
- [github.com/mattn/go-sqlite3]
- [github.com/jinzhu/gorm]
- [github.com/dhowden/tag]



as depicted by the `go.mod` file
```go
module github.com/tfwio/sekhem

go 1.12

require (
	github.com/dhowden/tag v0.0.0-20190519100835-db0c67e351b1
	github.com/gin-gonic/gin v1.4.0
	github.com/jinzhu/gorm v1.9.10
	github.com/kr/pretty v0.1.0 // indirect
	github.com/mattn/go-sqlite3 v1.10.0
	github.com/urfave/cli v1.20.0
	golang.org/x/crypto v0.0.0-20190605123033-f99c8df09eb5
)
```