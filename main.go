package main

import (
	"github.com/gotk3/gotk3/glib"
	"log"
	"os"

	"github.com/gotk3/gotk3/gtk"
)

const appId = "com.github.alivesubstance.zooverseer"

// todo change to relative path
const gladeFilePath = "/home/mirian/code/go/src/github.com/alivesubstance/zooverseer/assets/main.glade"

// there are rumors that global variable is evil. why?
var Builder *gtk.Builder

func main() {
	log.Print("Starting zooverseer")

	app, err := gtk.ApplicationNew(appId, glib.APPLICATION_FLAGS_NONE)
	checkError(err)

	app.Connect("activate", onAppActivate(app))

	os.Exit(app.Run(os.Args))
}

func onAppActivate(app *gtk.Application) func() {
	return func() {
		log.Print("Reading glade file")
		builder, err := gtk.BuilderNewFromFile(gladeFilePath)
		checkError(err)

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

	return connDialog
}

func onConnAddBtnClicked() {
	log.Print("Conn add btn clicked")
	connList := getObject("connList").(*gtk.ListBox)
	connList.SetHAlign(gtk.ALIGN_START)
	l, err := gtk.LabelNew("nightly-pleeco")
	checkError(err)

	connList.Add(l)
	connList.ShowAll()
}

func onConnDialogCancelBtnClicked(connDialog *gtk.Dialog) func() {
	return func() {
		connDialog.Hide()
	}
}
