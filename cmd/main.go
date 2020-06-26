package main

import (
	"fmt"
	"github.com/alivesubstance/zooverseer/core"
	"github.com/alivesubstance/zooverseer/ui"
	"github.com/alivesubstance/zooverseer/util"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"time"
)

const (
	timeFormatLogMsg      = "2006-01-02 15:04:05,999"
	timeFormatLogFileName = "20060102_150405"
)

func main() {
	initLogger()

	app, err := gtk.ApplicationNew(core.Config.AppId, glib.APPLICATION_FLAGS_NONE)
	util.CheckError(err)
	//addAppIcon(app)

	app.Connect("startup", func() { log.Println("Starting zooverseer") })
	app.Connect("activate", ui.OnAppActivate(app))
	app.Connect("shutdown", func() { log.Println("Stopping zooverseer") })

	os.Exit(app.Run(os.Args))
}

func addAppIcon(app *gtk.Application) {
	//https://stackoverflow.com/questions/45162862/how-do-i-set-an-icon-for-the-whole-application-using-pygobject
	//iconTheme, err := gtk.IconThemeGetDefault()
	//util.CheckError(err)
	//iconTheme.LoadIcon()
}

// todo clean old logs
func initLogger() {
	log.SetLevel(log.TraceLevel)
	//log.SetFormatter(&log.TextFormatter{
	//	DisableColors:   true,
	//	TimestampFormat: "2006-01-02 15:04:05,999",
	//	PadLevelText:    true,
	//})

	_, err := os.Stat(core.Config.LogDir)
	if os.IsNotExist(err) {
		err := os.Mkdir(core.Config.LogDir, 0775)
		if err != nil {
			log.WithError(err).Panicf("Failed to create log dir %v", core.Config.LogDir)
		}
	}

	logFileName := core.Config.LogDir + "/zooverseer-" + time.Now().Format(timeFormatLogFileName) + ".log"
	logFile, err := os.Create(logFileName)
	if err == nil {
		log.SetOutput(logFile)
	} else {
		log.WithError(err).Panicf("Failed to create log file, using default stdout")
		log.SetOutput(os.Stdout)
	}

	cleanOldLogs()
}

func cleanOldLogs() {
	files, _ := ioutil.ReadDir(core.Config.LogDir)
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() > files[j].Name()
	})
	if len(files) > core.Config.LogFilesHistorySize {
		files = files[core.Config.LogFilesHistorySize:]
		for _, file := range files {
			os.Remove(path.Join(core.Config.LogDir, file.Name()))
		}
	}
}

type PlainFormatter struct {
}

func (f *PlainFormatter) Format(entry *log.Entry) ([]byte, error) {
	msg := fmt.Sprintf(
		"%v [%-5v] %v %v",
		time.Now().Format(timeFormatLogMsg),
		entry.Level,
		entry.Message,
		entry.Data,
	)

	return append(util.StrToBytes(msg), '\n'), nil
}
