#!/usr/bin/env bash

# go tool dist list -json
export copyto=/c/Users/tfwro/Desktop/DesktopMess/.git-react/react-fomantic

# default
export GOARCH=amd64
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

for i in "$@"
do
  case "${i}" in
    linstall)
      dist install -v pkg/runtime
      go install -v -a std
      ;;
  esac
done
for i in "$@"
do
  case "${i}" in
    v | verbose | -v)
      echo verbose
      ;;
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
for i in "$@"
do
  case "${i}" in
    build)
      echo building
      echo go build -o srv-${GOARCH}.exe ${buildmod} srv.go
      go build -o srv-${GOARCH}.exe ${buildmod} srv.go
      ;;
  esac
done
for i in "$@"
do
  case "${i}" in
    build)
      cp -f  srv-${GOARCH}.exe ${copyto}/srv.exe
      cp -f  srv-${GOARCH}.exe srv.exe
      ;;
  esac
done
exit






