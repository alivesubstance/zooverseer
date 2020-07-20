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



#  cp -v $(USR_MINGW)/bin/{libgcc_s_seh-1.dll,libatk-1.0-0.dll,libbz2-1.dll,libcairo-2.dll,libcairo-gobject-2.dll,libepoxy-0.dll,libexpat-1.dll,libffi-6.dll,libfontconfig-1.dll,libfreetype-6.dll,libgdk-3-0.dll,libgdk_pixbuf-2.0-0.dll,libgio-2.0-0.dll,libglib-2.0-0.dll,libgmodule-2.0-0.dll,libgobject-2.0-0.dll,libgraphite2.dll,libgtk-3-0.dll,libharfbuzz-0.dll,libiconv-2.dll,libintl-8.dll,libjasper.dll,libjpeg-8.dll,libpango-1.0-0.dll,libpangocairo-1.0-0.dll,libpangoft2-1.0-0.dll,libpangowin32-1.0-0.dll,libpcre-1.dll,libpixman-1-0.dll,libpng16-16.dll,libstdc++-6.dll,libwinpthread-1.dll,zlib1.dll} $(OUT)/
