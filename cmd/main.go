package main

import (
	"github.com/alivesubstance/zooverseer/core"
	"github.com/alivesubstance/zooverseer/ui"
	"github.com/alivesubstance/zooverseer/util"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

func main() {
	log.Info("Starting zooverseer")

	initLogger()

	app, err := gtk.ApplicationNew(core.Config.AppId, glib.APPLICATION_FLAGS_NONE)
	util.CheckError(err)

	app.Connect("activate", ui.OnAppActivate(app))

	os.Exit(app.Run(os.Args))
}

// todo clean old logs
func initLogger() {
	log.SetReportCaller(true)
	log.SetLevel(log.TraceLevel)
	log.SetFormatter(&log.TextFormatter{
		DisableColors:   true,
		TimestampFormat: "2006-01-02 15:04:05,999",
		PadLevelText:    true,
	})

	_, err := os.Stat(core.Config.LogDir)
	if os.IsNotExist(err) {
		err := os.Mkdir(core.Config.LogDir, 0775)
		if err != nil {
			log.WithError(err).Fatalf("Failed to create log dir %v", core.Config.LogDir)
		}
	}

	logFileName := core.Config.LogDir + "/zooverseer-" + time.Now().Format("20060102_150405") + ".log"
	logFile, err := os.Create(logFileName)
	if err == nil {
		log.SetOutput(logFile)
	} else {
		log.WithError(err).Info("Failed to create log file, using default stdout")
		log.SetOutput(os.Stdout)
	}
}
