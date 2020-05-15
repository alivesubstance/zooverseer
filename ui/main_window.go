package ui

import "C"
import (
	"github.com/gotk3/gotk3/gtk"
)

var createNodeDlg *CreateNodeDlg
var nodeAction *NodeAction
var contextMenu *ContextMenu

func InitMainWindow(mainWindow *gtk.Window) {
	createNodeDlg = NewCreateNodeDlg(mainWindow)
	nodeAction = NewNodeAction()
	contextMenu = NewContextMenu()

	initNodeTree()
	notebook.init()
	initMainMenu()

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
