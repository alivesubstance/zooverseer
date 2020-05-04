package ui

import "C"
import (
	"github.com/gotk3/gotk3/gtk"
)

func InitMainWindow(mainWindow *gtk.Window) {
	// todo for test purpose only
	mainWindow.Move(0, 0)

	initNodeTree()
	notebook.init()
	initMainMenu()
	initContextMenu()

	mainWindow.SetTitle("Zooverseer")
	mainWindow.ShowAll()
}

func initMainMenu() {
	menuConnect := getObject("menuConnect").(*gtk.MenuItem)
	menuConnect.Connect("activate", func() {
		connDialog := getObject("connDialog").(*gtk.Dialog)
		connDialog.Show()
	})
}
