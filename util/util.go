package util

import (
	"log"
)

func CheckError(e error) {
	if e != nil {
		log.Panic(e)
	}
}

func CheckErrorWithMsg(msg string, e error) {
	if e != nil {
		log.Panic(msg, e)
	}
}

func BytesToString(data []byte) string {
	return string(data[:])
}
