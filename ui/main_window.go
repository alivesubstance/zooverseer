package ui

import "C"
import (
	"github.com/gotk3/gotk3/gtk"
)

func InitMainWindow(mainWindow *gtk.Window) {
	InitNodeTree()

	initMenuSignals()

	mainWindow.SetTitle("Zooverseer")
	mainWindow.ShowAll()

	// todo for test purpose only
	mainWindow.Move(1500, 0)
}

func initMenuSignals() {
	menuConnect := getObject("menuConnect").(*gtk.MenuItem)
	menuConnect.Connect("activate", func() {
		connDialog := getObject("connDialog").(*gtk.Dialog)
		connDialog.Show()
	})
}
