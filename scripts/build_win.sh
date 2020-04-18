#!/bin/bash

GOOS="windows"
GOARCH="amd64"

env GOOS="${GOOS}" GOARCH="${GOARCH}" go build -o /home/mirian/bin/zooverseer.exe github.com/alivesubstance/zooverseer/cmd
