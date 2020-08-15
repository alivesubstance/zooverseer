package core

import (
	"container/list"
	"fmt"
	"github.com/alivesubstance/zooverseer/util"
	log "github.com/sirupsen/logrus"
	"io"
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

type PlainFormatter struct{}

type CompositeWriter struct {
	writers *list.List
}

func (cw CompositeWriter) Write(p []byte) (n int, err error) {
	for e := cw.writers.Front(); e != nil; e = e.Next() {
		n, err = e.Value.(io.Writer).Write(p)
	}
	return n, err
}

func InitLogger() {
	log.SetLevel(Config.Log.Level)
	log.SetFormatter(&PlainFormatter{})
	//log.SetFormatter(&log.TextFormatter{
	//	DisableColors:   true,
	//	TimestampFormat: "2006-01-02 15:04:05,999",
	//	PadLevelText:    true,
	//})

	cleanOldLogs()

	compositeWriter := CompositeWriter{writers: list.New()}
	compositeWriter.writers.PushFront(os.Stdout)
	logFile := initLogFile()
	if logFile != nil {
		compositeWriter.writers.PushFront(logFile)
	}
	log.SetOutput(compositeWriter)
}

func initLogFile() *os.File {
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
		return logFile
	} else {
		log.WithError(err).Panicf("Failed to create log file, using default stdout")
	}

	return nil
}

func cleanOldLogs() {
	files, _ := ioutil.ReadDir(Config.Log.Dir)
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() > files[j].Name()
	})
	if len(files) > Config.Log.FilesHistorySize {
		files = files[Config.Log.FilesHistorySize-1:]
		for _, file := range files {
			filePath := path.Join(Config.Log.Dir, file.Name())
			err := os.Remove(filePath)
			if err != nil {
				log.WithError(err).Errorf("Failed to delete file %v", filePath)
			}

		}
	}
}

func (f *PlainFormatter) Format(entry *log.Entry) ([]byte, error) {
	msg := fmt.Sprintf(
		"%-23v [%-5v] %v",
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
