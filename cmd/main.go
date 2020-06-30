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
	core.InitLogger()
	log.Println("Starting zooverseer")

	app, err := gtk.ApplicationNew(core.Config.AppId, glib.APPLICATION_FLAGS_NONE)
	util.CheckError(err)

	app.Connect("activate", ui.OnAppActivate(app))
	app.Connect("shutdown", func() { log.Println("Stop zooverseer") })

	os.Exit(app.Run(os.Args))
}
