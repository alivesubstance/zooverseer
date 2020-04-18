#!/bin/bash

GOOS="linux"
GOARCH="amd64"

env GOOS="${GOOS}" GOARCH="${GOARCH}" go build -o /home/mirian/bin/zooverseer github.com/alivesubstance/zooverseer/cmd
