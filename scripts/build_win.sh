#!/bin/bash

GOOS="windows"
GOARCH="amd64"

SRC_ASSETS_DIR="assets"
INSTALL_DIR="../bin/win-1.0.0b"
INSTALL_ASSETS_DIR=${INSTALL_DIR}/${SRC_ASSETS_DIR}

rm -rf "${INSTALL_DIR}"
mkdir -p "${INSTALL_ASSETS_DIR}"
cp "../${SRC_ASSETS_DIR}/main.glade" "${INSTALL_ASSETS_DIR}/main.glade"
cp "../${SRC_ASSETS_DIR}/main.glade" "${INSTALL_ASSETS_DIR}/main.glade"
cp "../${SRC_ASSETS_DIR}/logo.png" "${INSTALL_ASSETS_DIR}/logo.png"
cp "../${SRC_ASSETS_DIR}/logo_small.png" "${INSTALL_ASSETS_DIR}/logo_small.png"
cp "../${SRC_ASSETS_DIR}/style.css" "${INSTALL_ASSETS_DIR}/style.css"

GOOS="${GOOS}" GOARCH="${GOARCH}" \
  CGO_ENABLED=1 CXX=x86_64-w64-mingw32-g++ CC=x86_64-w64-mingw32-gcc \
  NO_AT_BRIDGE=1 \
  go build -o "${INSTALL_DIR}/zooverseer.exe" \
  github.com/alivesubstance/zooverseer/cmd