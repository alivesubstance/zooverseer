#!/bin/bash

GOOS="windows"
GOARCH="amd64"

#PATH="${PATH}:./bin"
SRC_ASSETS_DIR="assets"
#BIN_DIR="../bin"
INSTALL_DIR="../bin/win-1.0.0b"
INSTALL_ASSETS_DIR=${INSTALL_DIR}/${SRC_ASSETS_DIR}

rm -rf "${INSTALL_DIR}"
mkdir -p "${INSTALL_ASSETS_DIR}"
cp "../${SRC_ASSETS_DIR}/main.glade" "${INSTALL_ASSETS_DIR}/main.glade"
cp "../${SRC_ASSETS_DIR}/logo.png" "${INSTALL_ASSETS_DIR}/logo.png"
cp "../${SRC_ASSETS_DIR}/logo_small.png" "${INSTALL_ASSETS_DIR}/logo_small.png"
cp "../${SRC_ASSETS_DIR}/style.css" "${INSTALL_ASSETS_DIR}/style.css"

GOOS="${GOOS}" GOARCH="${GOARCH}" \
  CGO_ENABLED=1 CXX=i686-w64-mingw32-g++ CC=i686-w64-mingw32-gcc \
  go build -o "${INSTALL_DIR}/zooverseer.exe" \
  github.com/alivesubstance/zooverseer/cmd

