package ui

import (
	"github.com/alivesubstance/zooverseer/core"
	"github.com/alivesubstance/zooverseer/util"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	log "github.com/sirupsen/logrus"
	"time"
)

// there are rumors that global variable is evil. why?
var Builder *gtk.Builder

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

func TestExport() {
	// todo exclude, test only
	time.Sleep(1 * time.Second)
	log.Tracef("Eagerly select path")
	selection, _ := getNodesTreeView().GetSelection()
	treeIter, _ := nodeTreeStore.GetIterFromString("0:0")
	selection.SelectIter(treeIter)
	contextMenu.onExportNode()
}

func enableSpinner(enable bool) {
	box := GetObject("mainWindow").(*gtk.Window)
	box.SetSensitive(!enable)

	//if enable {
	//spinner, _ := gtk.SpinnerNew()
	//spinner.SetSensitive(true)
	//spinner.ShowAll()
	//spinner.Start()
	//
	//overlay := GetObject("nodeOverlay").(*gtk.Overlay)
	//overlay.AddOverlay(spinner)
	//m.spinner.Show()
	//m.spinner.Start()
	//} else {
	//overlay := GetObject("nodeOverlay").(*gtk.Overlay)
	//overlay.rem
	//m.spinner.Stop()
	//m.spinner.Hide()
	//}
}
