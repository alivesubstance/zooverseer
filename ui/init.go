package ui

import (
	"github.com/alivesubstance/zooverseer/core"
	"github.com/alivesubstance/zooverseer/util"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	log "github.com/sirupsen/logrus"
)

// there are rumors that global variable is evil. why?
var (
	Builder *gtk.Builder
)

func OnAppActivate(app *gtk.Application) func() {
	return func() {
		log.Info("Reading glade file")
		builder, err := gtk.BuilderNewFromFile(core.GladeFilePath)
		util.CheckError(err)

		Builder = builder

		mainWindow := getMainWindow()
		InitMainWindow(mainWindow)

		//InitConnDialog(mainWindow)

		app.AddWindow(mainWindow)
	}
}

func getMainWindow() *gtk.Window {
	return getObject("mainWindow").(*gtk.Window)
}

func createConfirmDialog(parent gtk.IWindow, text string) *gtk.MessageDialog {
	return gtk.MessageDialogNew(parent, gtk.DIALOG_MODAL, gtk.MESSAGE_QUESTION, gtk.BUTTONS_YES_NO, text)
}

func createInfoDialog(parent gtk.IWindow, text string) *gtk.MessageDialog {
	return gtk.MessageDialogNew(parent, gtk.DIALOG_MODAL, gtk.MESSAGE_INFO, gtk.BUTTONS_OK, text)
}

func createWarnDialog(parent gtk.IWindow, text string) *gtk.MessageDialog {
	return gtk.MessageDialogNew(parent, gtk.DIALOG_MODAL, gtk.MESSAGE_WARNING, gtk.BUTTONS_OK, text)
}

func createAndRunWarnDialog(parent gtk.IWindow, text string) {
	dlg := gtk.MessageDialogNew(parent, gtk.DIALOG_MODAL, gtk.MESSAGE_WARNING, gtk.BUTTONS_OK, text)
	dlg.Run()
	dlg.Hide()
}

func createErrorDialog(parent gtk.IWindow, text string) *gtk.MessageDialog {
	return gtk.MessageDialogNew(parent, gtk.DIALOG_MODAL, gtk.MESSAGE_ERROR, gtk.BUTTONS_OK, text)
}

func getObject(objectName string) glib.IObject {
	object, err := Builder.GetObject(objectName)
	util.CheckError(err)

	return object
}
