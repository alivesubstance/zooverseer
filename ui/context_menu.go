package ui

import (
	"github.com/atotto/clipboard"
	"github.com/gotk3/gotk3/gtk"
	log "github.com/sirupsen/logrus"
)

func initContextMenu() {
	getObject("popupMenu").(*gtk.Menu).Connect("popped-up", func() {

	})

	getObject("createNodeDlgOkBtn").(*gtk.Button).Connect("activate", onNewNodeDlgOkBtn)
	getObject("createNodeDlgCancelBtn").(*gtk.Button).Connect("activate", onNewNodeDlgCancelBtn)

	getObject("popupMenuAdd").(*gtk.MenuItem).Connect("activate", onAddNewNode)
	getObject("popupMenuCopyValue").(*gtk.MenuItem).Connect("activate", onCopyValue)
	getObject("popupMenuCopyValue").(*gtk.MenuItem).Connect("activate", onCopyValue)
}

func onAddNewNode() {
	newNodeDlg := getNewNodeDlg()
	newNodeDlg.ShowAll()
}

func onNewNodeDlgOkBtn() {
	newNodeDlgEntry := getObject("newNodeDlgEntry").(*gtk.Entry)
	text, _ := newNodeDlgEntry.GetText()
	if len(text) == 0 {
		dialog := createWarnDialog(getNewNodeDlg(), "New node name should be provided")
		dialog.Run()
		dialog.Hide()
	} else {
		//node := getNodesTreeSelectedValue()

		getNewNodeDlg().Close()
	}
}

func onNewNodeDlgCancelBtn() {
	getNewNodeDlg().Close()
}

func onCopyValue() {
	treeSelection, _ := getNodesTreeView().GetSelection()
	node, _ := getTreeSelectedNode(treeSelection)
	if node != nil {
		log.Tracef("Add text value to clipboard")
		err := clipboard.WriteAll(node.Value)
		if err != nil {
			log.WithError(err).Errorf("Failed to copy value to clipboard")
		}
	}
}

func getNewNodeDlg() *gtk.Dialog {
	return getObject("newNodeDlg").(*gtk.Dialog)
}
