package main

import (
	"github.com/alivesubstance/zooverseer/core"
	"github.com/alivesubstance/zooverseer/ui"
	"github.com/alivesubstance/zooverseer/util"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	log.Info("Starting zooverseer")

	initLogger()

	app, err := gtk.ApplicationNew(core.AppId, glib.APPLICATION_FLAGS_NONE)
	util.CheckError(err)

	app.Connect("activate", ui.OnAppActivate(app))

	os.Exit(app.Run(os.Args))
}

func initLogger() {
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.TraceLevel)
}
