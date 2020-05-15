package ui

import (
	"github.com/atotto/clipboard"
	"github.com/gotk3/gotk3/gtk"
	log "github.com/sirupsen/logrus"
)

type ContextMenu struct {
	addItem        *gtk.MenuItem
	copyValueItem  *gtk.MenuItem
	copyNameItem   *gtk.MenuItem
	copyPathItem   *gtk.MenuItem
	pasteValueItem *gtk.MenuItem
	renameItem     *gtk.MenuItem
	deleteItem     *gtk.MenuItem
}

func NewContextMenu() *ContextMenu {
	contextMenu := &ContextMenu{}
	contextMenu.addItem = getObject("contextMenuAdd").(*gtk.MenuItem)
	contextMenu.addItem.Connect("activate", onAddNewNode)

	contextMenu.copyValueItem = getObject("contextMenuCopyValue").(*gtk.MenuItem)
	contextMenu.copyValueItem.Connect("activate", onCopyValue)

	contextMenu.copyNameItem = getObject("contextMenuCopyName").(*gtk.MenuItem)
	contextMenu.copyNameItem.Connect("activate", onCopyValue)

	contextMenu.copyPathItem = getObject("contextMenuCopyPath").(*gtk.MenuItem)
	contextMenu.copyPathItem.Connect("activate", onCopyValue)

	contextMenu.pasteValueItem = getObject("contextMenuPasteValue").(*gtk.MenuItem)
	contextMenu.pasteValueItem.Connect("activate", onCopyValue)

	contextMenu.renameItem = getObject("contextMenuRename").(*gtk.MenuItem)
	contextMenu.renameItem.Connect("activate", onCopyValue)

	contextMenu.deleteItem = getObject("contextMenuDeleteNode").(*gtk.MenuItem)
	contextMenu.deleteItem.Connect("activate", onCopyValue)

	contextMenu.enableButtons(false)

	return contextMenu
}

func onAddNewNode() {
	nodeAction.onNodeCreateBtnClicked()
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

func (n *ContextMenu) enableButtons(enabled bool) {
	n.addItem.SetSensitive(enabled)
	n.copyValueItem.SetSensitive(enabled)
	n.copyNameItem.SetSensitive(enabled)
	n.copyPathItem.SetSensitive(enabled)
	n.pasteValueItem.SetSensitive(enabled)
	n.renameItem.SetSensitive(enabled)
	n.deleteItem.SetSensitive(enabled)
}
