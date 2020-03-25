package main

import (
	"github.com/alivesubstance/zooverseer/core"
	"github.com/alivesubstance/zooverseer/util"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"log"
	"os"
)

// there are rumors that global variable is evil. why?
var Builder *gtk.Builder
var ConnRepository core.ConnRepository = core.JsonConnInfo{}

func main() {
	log.Print("Starting zooverseer")

	app, err := gtk.ApplicationNew(core.AppId, glib.APPLICATION_FLAGS_NONE)
	util.CheckError(err)

	app.Connect("activate", onAppActivate(app))

	os.Exit(app.Run(os.Args))
}

func onAppActivate(app *gtk.Application) func() {
	return func() {
		log.Print("Reading glade file")
		builder, err := gtk.BuilderNewFromFile(core.GladeFilePath)
		util.CheckError(err)

		Builder = builder

		mainWindow := getObject("mainWindow").(*gtk.Window)

		connDialog := initConnDialog(mainWindow)
		connDialog.ShowAll()

		mainWindow.ShowAll()
		app.AddWindow(mainWindow)
	}
}

func initConnDialog(mainWindow *gtk.Window) *gtk.Dialog {
	connPortEntry := getObject("connPortEntry").(*gtk.Entry)
	connPortEntry.SetWidthChars(10)

	connDialog := getObject("connDialog").(*gtk.Dialog)
	connDialog.SetTransientFor(mainWindow)

	connDialogCancelBtn := getObject("connDialogCancelBtn").(*gtk.Button)
	connDialogCancelBtn.Connect("clicked", onConnDialogCancelBtnClicked(connDialog))

	connAddBtn := getObject("connAddBtn").(*gtk.Button)
	connAddBtn.Connect("clicked", onConnAddBtnClicked)

	initConnsListBox()

	return connDialog
}

func initConnsListBox() {
	//connNameEntry := getObject("connNameEntry").(*gtk.Entry)
	//connHostEntry := getObject("connHostEntry").(*gtk.Entry)
	//connPortEntry := getObject("connPortEntry").(*gtk.Entry)
	//connUserEntry := getObject("connUserEntry").(*gtk.Entry)
	//connPwdEntry := getObject("connPwdEntry").(*gtk.Entry)

	connListBox := getConnListBox()
	connListBox.Connect("row-selected", func() {
		connName, err := connListBox.GetSelectedRow().GetName()
		util.CheckError(err)

		log.Print("Selected row " + connName)
	})

	connInfos := ConnRepository.FindAll()
	for _, connInfo := range connInfos {
		label, err := gtk.LabelNew(connInfo.Name)
		util.CheckError(err)

		connListBox.Add(label)
	}
	connListBox.SelectRow(connListBox.GetRowAtIndex(0))
	connListBox.ShowAll()
}

func onConnAddBtnClicked() {
	log.Print("Conn add btn clicked")
}

func getConnListBox() *gtk.ListBox {
	return getObject("connList").(*gtk.ListBox)
}

func onConnDialogCancelBtnClicked(connDialog *gtk.Dialog) func() {
	return func() {
		connDialog.Hide()
	}
}

func getObject(objectName string) glib.IObject {
	object, err := Builder.GetObject(objectName)
	util.CheckError(err)

	return object
}
