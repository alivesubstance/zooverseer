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
	spinner *gtk.Spinner
)

func OnAppActivate(app *gtk.Application) func() {
	return func() {
		log.Info("Reading glade file")
		builder, err := gtk.BuilderNewFromFile(core.GladeFilePath)
		util.CheckError(err)

		Builder = builder

		mainWindow := GetMainWindow()
		InitMainWindow(mainWindow)
		InitConnDialog(mainWindow)

		app.AddWindow(mainWindow)
	}
}

func GetMainWindow() *gtk.Window {
	return getObject("mainWindow").(*gtk.Window)
}

func CreateErrorDialog(parent gtk.IWindow, text string) *gtk.MessageDialog {
	return gtk.MessageDialogNew(parent, gtk.DIALOG_MODAL, gtk.MESSAGE_ERROR, gtk.BUTTONS_OK, text)
}

func createConfirmDialog(parent gtk.IWindow, text string) *gtk.MessageDialog {
	dlg := gtk.MessageDialogNew(parent, gtk.DIALOG_MODAL, gtk.MESSAGE_QUESTION, gtk.BUTTONS_YES_NO, text)
	//todo doesn't work
	//dlg.SetPosition(gtk.WIN_POS_CENTER_ON_PARENT)
	return dlg
}

func createInfoDialog(parent gtk.IWindow, text string) *gtk.MessageDialog {
	return gtk.MessageDialogNew(parent, gtk.DIALOG_MODAL, gtk.MESSAGE_INFO, gtk.BUTTONS_OK, text)
}

func createWarnDialog(parent gtk.IWindow, text string) *gtk.MessageDialog {
	return gtk.MessageDialogNew(parent, gtk.DIALOG_MODAL, gtk.MESSAGE_WARNING, gtk.BUTTONS_OK, text)
}

func getObject(objectName string) glib.IObject {
	object, err := Builder.GetObject(objectName)
	util.CheckError(err)

	return object
}

// todo replace with busy dialog
func enableSpinner(enable bool) {
	if enable {
		log.Infof("Enable spinner")
		styleContext, _ := getNodesTreeView().GetStyleContext()
		styleContext.AddClass("font-disabled")

		createInfoDialog(GetMainWindow(), "Export in progress").Run()
	} else {
		styleContext, _ := getNodesTreeView().GetStyleContext()
		styleContext.AddClass("font-enabled")
		log.Infof("Disable spinner")
	}

	GetMainWindow().SetSensitive(!enable)
	//
	//if enable {
	//	spinner, _ = gtk.SpinnerNew()
	//	getNodeOverlay().AddOverlay(spinner)
	//	spinner.Start()
	//} else {
	//	if spinner != nil {
	//		getNodeOverlay().Remove(spinner)
	//	}
	//}
}

func getNodeOverlay() *gtk.Overlay {
	return GetObject("nodeOverlay").(*gtk.Overlay)
}
