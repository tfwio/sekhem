#!/usr/bin/env bash

# fyi: go tool dist list -json
# sub-commands: mod linux clean, l-install, 386 amd64, verbose tidy, tosrv tofo
# cp2 uses jsoniter

# set a copyto directory.
export copyto=/c/Users/tfwro/Desktop/DesktopMess/.git-react/react-box
# default[mod:GO111MODULE=on]
export GOARCH=amd64
# 
# This script does a few things, all of which require some form of input.
if [ "${#}" == "0" ]; then
  echo no args so we\'ll do a standard build
  echo sub-commands\: mod linux clean, l-install, 386 amd64, verbose tidy, tosrv tofo
  echo cp, cp2 \(uses jsoniter\)
  exit 0
  # standard_stuff clean build mod
fi
# mod:linux:clean
do_init(){
  for i in "$@"
  do
    case "${i}" in
      mod)
        export GO111MODULE=on
        export buildmod="-mod vendor "
        echo GO111MODULE=${GO111MODULE}
        ;;
      linux)
        export GOOS=linux
        echo GOOS=${GOOS}
        ;;
      clean)
        echo cleaning
        go clean
        ;;
      *)
        ;;
    esac
  done
}
# l-install
do_pre(){
  #   installs linux package/sources
  for i in "$@"
  do
    case "${i}" in
      l-install)
        dist install -v pkg/runtime
        go install -v -a std
        ;;
    esac
  done

}
# vendor
do_vendor(){
  for i in "$@"
  do
    case "${i}" in
      vendor)
        # do_init "mod"
        do_get GO111MODULE on
        if [[ x$CK != xtrue ]]; then do_init mod ; fi
        echo ==\> check/download vendor stuff
        go mod vendor
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
        echo go build -tags=jsoniter -o srv-${GOARCH}.exe ${buildmod} *.go
        go build -tags=jsoniter -o srv-${GOARCH}.exe ${buildmod} *.go
        ;;
    esac
  done
}

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
    esac
  done
}

# verbose:tidy
do_flags(){
  for i in "$@"
  do
    case "${i}" in
      tidy)
        do_get GO111MODULE on
        if [[ x$CK != xtrue ]]; then do_init mod ; fi
        go mod tidy
        ;;
      v | verbose | -v)
        echo verbose
        ;;
    esac
  done
}

# tosrv:tofo
do_copy(){
  for i in "$@"
  do
    case "${i}" in
      m2)
        echo go build -o m2.exe util/m2.go
        go build -o m2.exe util/m2.go
        ;;
      cp)
        echo GO111MODULE=on go build -o srv.exe -mod vendor srv.go \&\& cp -f  srv.exe ${copyto}
        GO111MODULE=on go build -o srv.exe -mod vendor srv.go && cp -f  srv.exe ${copyto}
        ;;
      cp2)
        echo GO111MODULE=on go build -tags=jsoniter -o srv.exe -mod vendor srv.go \&\& cp -f  srv.exe ${copyto}
        GO111MODULE=on go build -tags=jsoniter -o srv.exe -mod vendor srv.go && cp -f  srv.exe ${copyto}
        ;;
      tofo)
        echo copying to
        echo destination\: ${copyto}/srv.exe
        echo destination\: ${copyto}
        cp -f  srv-${GOARCH}.exe ${copyto}/srv.exe
        cp -f  srv-${GOARCH}.exe ${copyto}
        ;;
      tosrv)
        echo copying to
        echo destination\: ${copyto}/srv.exe
        echo destination\: srv.exe
        cp -f  srv-${GOARCH}.exe ${copyto}/srv.exe
        cp -f  srv-${GOARCH}.exe srv.exe
        ;;
    esac
  done
}

function do_get() {
  varname=${1}
  varvar=${!varname}
  if [[ "x$varvar" == "x${2}" ]]; then
    export CK=true
  fi
}

function dodo(){
  do_init "$@"
  do_vendor "$@"
  do_pre "$@"
  do_flags "$@"
  do_plat "$@"
  do_build "$@"
  do_copy "$@"
}


do_init "$@"
do_vendor "$@"
do_pre "$@"
do_flags "$@"
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