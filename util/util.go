package util

import (
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
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

func Int64ToStr(value int64) string {
	return strconv.FormatInt(value, 10)
}

func Int32ToStr(value int32) string {
	return Int64ToStr(int64(value))
}

func MillisToTime(millis int64) time.Time {
	return time.Unix(0, millis*int64(time.Millisecond))
}
