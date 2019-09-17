#!/usr/bin/env bash

# fyi: go tool dist list -json
# sub-commands: mod linux clean, l-install, 386 amd64, verbose tidy, tosrv tofo
# cp2 uses jsoniter (login | crypt : compile ./login.exe)

# set a copyto directory. (this would be my react sandbox)
export copyto=/c/Users/tfwro/Desktop/DesktopMess/.git-react/react-box
# default[mod:GO111MODULE=on]
export GOARCH=amd64
# 
# This script does a few things, all of which require some form of input.
if [ "${#}" == "0" ]; then
  echo no args so we\'ll do nothing
  echo 
  echo flags\:
  echo "  - 386 amd64 v"
  echo 
  echo commands are\:
  echo 
  echo "  - bootstrap | bs"
  echo "  - standard build: just call './do build'"
  echo 
  echo 
  exit 0
  # standard_stuff clean build mod
fi

# 386:amd64
do_plat(){
  # 386:amd64
  for i in "$@"
  do
    case "${i}" in
      386)
        export GOARCH=386
        echo GOARCH=${GOARCH}
        ;;
      amd64)
        export GOARCH=amd64
        echo GOARCH=${GOARCH}
        ;;
      v | verbose | -v)
        echo verbose
        ;;
    esac
  done
}

# build
do_build(){
  for i in "$@"
  do
    case "${i}" in
      build)
        echo building

        echo go clean
        go clean

        echo go build -tags jsoniter -o srv-${GOARCH}.exe ${buildmod} srv.go
        go build -tags jsoniter -o srv-${GOARCH}.exe ${buildmod} srv.go
        ;;
    esac
  done
}
# mod:linux:clean
do_init(){
  for i in "$@"
  do
    case "${i}" in
      bootstrap | bs)
        export GO111MODULE=on
        echo "- go and get github.com/json-iterator/go (a gin-gonic/gin \`-tags=jsoniter\` extension)"
        echo " ** a github.com/gin-gonic/gin \`-tags=jsoniter\` extension **"
        go get github.com/json-iterator/go
        echo "- delete go.mod"
        rm -f go.mod
        echo "- go mod tidy"
        go mod tidy
        echo "- go mod vendor"
        go mod vendor
        echo "- go mod download"
        go mod download
        ;;
    esac
  done
}

# tosrv:tofo
do_copy(){
  for i in "$@"
  do
    case "${i}" in
      cp2)
        echo GO111MODULE=on go build -tags jsoniter -o srv.exe -mod vendor srv.go \&\& cp -f  srv.exe ${copyto}
        GO111MODULE=on go build -tags jsoniter -o srv.exe -mod vendor srv.go && cp -f  srv.exe ${copyto}
        ;;
    esac
  done
}

function dodo(){
  do_init "$@"
  do_plat "$@"
  do_build "$@"
  do_copy "$@"
}


do_init "$@"
do_plat "$@"
do_build "$@"
do_copy "$@"

exit 0
bash_trick() {
  s='Strings:With:Four:Words'
  IFS=: read -r var1 var2 var3 var4 <<< "$s"
  echo "[$var1] [$var2] [$var3 [$var4]"
  exit
}

do_this() {
  todo=$1
  args=$2
  shift 2
  IFS=: read -r var1 var2 var3 var4 var5 var6 <<< "$args"
  echo "[$var1] [$var2] [$var3 [$var4] [$var5] [$va6]"
  exit
}