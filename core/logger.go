package core

import (
	"fmt"
	"github.com/alivesubstance/zooverseer/util"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"time"
)

const (
	timeFormatLogMsg      = "2006-01-02 15:04:05.999"
	timeFormatLogFileName = "20060102_150405"
)

func InitLogger() {
	log.SetLevel(Config.Log.Level)
	log.SetFormatter(&PlainFormatter{})
	//log.SetFormatter(&log.TextFormatter{
	//	DisableColors:   true,
	//	TimestampFormat: "2006-01-02 15:04:05,999",
	//	PadLevelText:    true,
	//})

	cleanOldLogs()
	initLogFile()
}

func initLogFile() {
	_, err := os.Stat(Config.Log.Dir)
	if os.IsNotExist(err) {
		err := os.Mkdir(Config.Log.Dir, 0775)
		if err != nil {
			log.WithError(err).Panicf("Failed to create log dir %v", Config.Log.Dir)
		}
	}

	logFileName := Config.Log.Dir + "/zooverseer-" + time.Now().Format(timeFormatLogFileName) + ".log"
	logFile, err := os.Create(logFileName)
	if err == nil {
		log.SetOutput(logFile)
	} else {
		log.WithError(err).Panicf("Failed to create log file, using default stdout")
		log.SetOutput(os.Stdout)
	}
}

func cleanOldLogs() {
	files, _ := ioutil.ReadDir(Config.Log.Dir)
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() > files[j].Name()
	})
	if len(files) > Config.Log.FilesHistorySize {
		files = files[Config.Log.FilesHistorySize:]
		for _, file := range files {
			filePath := path.Join(Config.Log.Dir, file.Name())
			err := os.Remove(filePath)
			if err != nil {
				log.WithError(err).Errorf("Failed to delete file %v", filePath)
			}

		}
	}
}

type PlainFormatter struct {
}

func (f *PlainFormatter) Format(entry *log.Entry) ([]byte, error) {
	msg := fmt.Sprintf(
		"[%-23v] [%-5v] %v",
		time.Now().Format(timeFormatLogMsg),
		entry.Level,
		entry.Message,
	)
	if entry.Data != nil && len(entry.Data) != 0 {
		msg = fmt.Sprintf("%v %v", msg, entry.Data)
	}
	if entry.Buffer != nil {
		msg = fmt.Sprintf("%v %v", msg, entry.Buffer)
	}

	return append(util.StrToBytes(msg), '\n'), nil
}
