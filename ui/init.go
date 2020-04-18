package ui

import (
	"github.com/alivesubstance/zooverseer/core"
	"github.com/alivesubstance/zooverseer/util"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"log"
)

// there are rumors that global variable is evil. why?
var (
	Builder *gtk.Builder
)

func OnAppActivate(app *gtk.Application) func() {
	return func() {
		log.Print("Reading glade file")
		builder, err := gtk.BuilderNewFromFile(core.GladeFilePath)
		util.CheckError(err)

		Builder = builder

		mainWindow := getObject("mainWindow").(*gtk.Window)
		InitMainWindow(mainWindow)

		InitConnDialog(mainWindow)

		app.AddWindow(mainWindow)
	}
}

func getObject(objectName string) glib.IObject {
	object, err := Builder.GetObject(objectName)
	util.CheckError(err)

	return object
}
