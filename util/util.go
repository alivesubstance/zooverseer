package util

import (
	log "github.com/sirupsen/logrus"
)

func CheckError(e error) {
	if e != nil {
		log.WithError(e).Panic()
	}
}

func CheckErrorWithMsg(msg string, e error) {
	if e != nil {
		log.WithError(e).Panic(msg)
	}
}

func BytesToString(data []byte) string {
	return string(data[:])
}
