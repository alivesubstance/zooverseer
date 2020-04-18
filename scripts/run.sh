#!/bin/bash

ASSETS_DIR="../assets"
INSTALL_DIR="${HOME}/env/zooverseer"

mkdir -p "${INSTALL_DIR}"
cp -a "${ASSETS_DIR}/." "${INSTALL_DIR}/config/"

go build -o "${INSTALL_DIR}/zooverseer" github.com/alivesubstance/zooverseer/cmd

cd "${INSTALL_DIR}"
zooverseer