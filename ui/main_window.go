package ui

import "C"
import (
	"github.com/gotk3/gotk3/gtk"
)

var createNodeDlg *CreateNodeDlg

func InitMainWindow(mainWindow *gtk.Window) {
	createNodeDlg = NewCreateNodeDlg(mainWindow)

	initNodeTree()
	notebook.init()
	initMainMenu()
	initContextMenu()
	initNodeActionSignal()

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

func initNodeActionSignal() {
	getObject("nodeCreateBtn").(*gtk.Button).Connect("clicked", onNodeCreateBtnClicked)
	getObject("nodeRefreshBtn").(*gtk.Button).Connect("clicked", onNodeRefreshBtnClicked)
	getObject("nodeDeleteBtn").(*gtk.Button).Connect("clicked", onNodeDeleteBtnClicked)
}

func onNodeCreateBtnClicked() {
	createNodeDlg.showAll()
}

func onNodeRefreshBtnClicked() {
	selection, _ := getNodesTreeView().GetSelection()
	parentPath, _ := getTreeSelectedZkPath(selection)
	refreshNode(parentPath)
}

func onNodeDeleteBtnClicked() {
	deleteSelectedNode()
}
