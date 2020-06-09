#!/bin/bash

# see https://wiki.archlinux.de/title/GNOME#Zugangshilfen
export NO_AT_BRIDGE=1

#PATH="${PATH}:./bin"
ASSETS_DIR="../assets"
#BIN_DIR="../bin"
INSTALL_DIR="${HOME}/env/zooverseer"

mkdir -p "${INSTALL_DIR}"
cp -a "${ASSETS_DIR}/." "${INSTALL_DIR}/conf"

go build -o "${INSTALL_DIR}/zooverseer" github.com/alivesubstance/zooverseer/cmd

cd "${INSTALL_DIR}"
${INSTALL_DIR}/zooverseer