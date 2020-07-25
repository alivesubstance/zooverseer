#!/bin/bash

GOOS="windows"
GOARCH="amd64"

SRC_ASSETS_DIR="assets"
INSTALL_DIR="../bin/win-1.0.0b"
INSTALL_ASSETS_DIR=${INSTALL_DIR}/${SRC_ASSETS_DIR}

rm -rf "${INSTALL_DIR}"
mkdir -p "${INSTALL_ASSETS_DIR}"
cp "../${SRC_ASSETS_DIR}/{main.glade,logo.png,logo_small.png,style.css}" "${INSTALL_ASSETS_DIR}/"

GOOS="${GOOS}" GOARCH="${GOARCH}" \
  CGO_ENABLED=1 CXX=x86_64-w64-mingw32-g++ CC=x86_64-w64-mingw32-gcc \
  NO_AT_BRIDGE=1 \
  go build -o "${INSTALL_DIR}/zooverseer.exe" \
  github.com/alivesubstance/zooverseer/cmd

