#!/bin/bash

GOOS="linux"
GOARCH="amd64"

export NO_AT_BRIDGE=1

#PATH="${PATH}:./bin"
SRC_ASSETS_DIR="assets"
#BIN_DIR="../bin"
INSTALL_DIR="../bin/linux-1.0.0b"
INSTALL_ASSETS_DIR=${INSTALL_DIR}/${SRC_ASSETS_DIR}

rm -rf "${INSTALL_DIR}"
mkdir -p "${INSTALL_ASSETS_DIR}"
cp "../${SRC_ASSETS_DIR}/main.glade" "${INSTALL_ASSETS_DIR}/main.glade"
cp "../${SRC_ASSETS_DIR}/style.css" "${INSTALL_ASSETS_DIR}/style.css"

env GOOS="${GOOS}" GOARCH="${GOARCH}" go build -o "${INSTALL_DIR}/zooverseer" github.com/alivesubstance/zooverseer/cmd
