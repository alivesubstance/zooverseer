#!/bin/bash

GOOS="linux"
GOARCH="amd64"

export NO_AT_BRIDGE=1

#PATH="${PATH}:./bin"
ASSETS_DIR="../assets"
#BIN_DIR="../bin"
INSTALL_DIR="../build/linux-1.0.0b"

rm -rf "${INSTALL_DIR}"
mkdir -p "${INSTALL_DIR}/conf"
cp "${ASSETS_DIR}/main.glade" "${INSTALL_DIR}/conf/main.glade"
cp "${ASSETS_DIR}/style.css" "${INSTALL_DIR}/conf/style.css"

env GOOS="${GOOS}" GOARCH="${GOARCH}" go build -o "${INSTALL_DIR}/zooverseer" github.com/alivesubstance/zooverseer/cmd
