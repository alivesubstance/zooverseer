package ui

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
	initMainMenu(mainWindow)

	mainWindow.SetTitle("Zooverseer")
	mainWindow.ShowAll()
}

func initMainMenu(mainWindow *gtk.Window) {
	getObject("menuConnect").(*gtk.MenuItem).Connect("activate", func() {
		connDialog := getObject("connDialog").(*gtk.Dialog)
		connDialog.Show()
	})
	getObject("menuExit").(*gtk.MenuItem).Connect("activate", func() {
		ZkCachingRepo.Close()
		mainWindow.Close()
	})

	getObject("menuDisconnect").(*gtk.MenuItem).Connect("activate", func() {
		ZkCachingRepo.Close()
		ClearNodeTree()
	})
}
