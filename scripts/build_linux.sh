#!/bin/bash

VERSION="1.0.0"
GOOS="linux"
GOARCH="amd64"
BIN_DIR="/home/mirian/code/zooverseer/bin"

SRC_ASSETS_DIR="assets"
INSTALL_DIR="$BIN_DIR/linux-$VERSION"
INSTALL_ASSETS_DIR=${INSTALL_DIR}/${SRC_ASSETS_DIR}

rm -rf "${INSTALL_DIR}"
mkdir -p "${INSTALL_ASSETS_DIR}"
cp "../${SRC_ASSETS_DIR}/main.glade" "${INSTALL_ASSETS_DIR}/"
cp "../${SRC_ASSETS_DIR}/logo.png" "${INSTALL_ASSETS_DIR}/"
cp "../${SRC_ASSETS_DIR}/logo_small.png" "${INSTALL_ASSETS_DIR}/"
cp "../${SRC_ASSETS_DIR}/style.css" "${INSTALL_ASSETS_DIR}/"

env \
  GOOS="${GOOS}" \
  GOARCH="${GOARCH}" \
  NO_AT_BRIDGE=1 \
  go build -o "${INSTALL_DIR}/zooverseer" github.com/alivesubstance/zooverseer/cmd

cp "${INSTALL_DIR}/zooverseer" /home/mirian/env/zooverseer/bin
cp -r "${INSTALL_ASSETS_DIR}/." /home/mirian/env/zooverseer/assets

cd "$INSTALL_DIR"
tar -zcvf  "$BIN_DIR"/zooverseer-"$VERSION".tar.gz $(ls -A)