package ui

import (
	"github.com/alivesubstance/zooverseer/core/zk"
	"github.com/atotto/clipboard"
	"github.com/gotk3/gotk3/gtk"
	log "github.com/sirupsen/logrus"
)

type ContextMenu struct {
	gtkMenu        *gtk.Menu
	addItem        *gtk.MenuItem
	copyValueItem  *gtk.MenuItem
	copyNameItem   *gtk.MenuItem
	copyPathItem   *gtk.MenuItem
	pasteValueItem *gtk.MenuItem
	renameItem     *gtk.MenuItem
	deleteItem     *gtk.MenuItem
	exportItem     *gtk.MenuItem
	spinner        *gtk.Spinner
}

func NewContextMenu() *ContextMenu {
	contextMenu := &ContextMenu{}
	contextMenu.gtkMenu = GetObject("contextMenu").(*gtk.Menu)

	contextMenu.addItem = GetObject("contextMenuAdd").(*gtk.MenuItem)
	contextMenu.addItem.Connect("activate", contextMenu.onAddNewNode)

	contextMenu.copyValueItem = GetObject("contextMenuCopyValue").(*gtk.MenuItem)
	contextMenu.copyValueItem.Connect("activate", contextMenu.onCopyValue)

	contextMenu.copyNameItem = GetObject("contextMenuCopyName").(*gtk.MenuItem)
	contextMenu.copyNameItem.Connect("activate", contextMenu.onCopyName)

	//contextMenu.copyPathItem = GetObject("contextMenuCopyPath").(*gtk.MenuItem)
	//contextMenu.copyPathItem.Connect("activate", onCopyValue)
	//
	//contextMenu.pasteValueItem = GetObject("contextMenuPasteValue").(*gtk.MenuItem)
	//contextMenu.pasteValueItem.Connect("activate", onCopyValue)
	//
	//contextMenu.renameItem = GetObject("contextMenuRename").(*gtk.MenuItem)
	//contextMenu.renameItem.Connect("activate", onCopyValue)
	//
	contextMenu.deleteItem = GetObject("contextMenuDeleteNode").(*gtk.MenuItem)
	contextMenu.deleteItem.Connect("activate", contextMenu.onDeleteNode)

	contextMenu.exportItem = GetObject("contextMenuExportNode").(*gtk.MenuItem)
	contextMenu.exportItem.Connect("activate", contextMenu.onExportNode)

	contextMenu.enableMenu(false)

	return contextMenu
}

func (m *ContextMenu) onAddNewNode() {
	nodeAction.onNodeCreateBtnClicked()
}

func (m *ContextMenu) onExportNode() {
	exportSelectedNode()
}

func (m *ContextMenu) onDeleteNode() {
	deleteSelectedNode()
}

func (m *ContextMenu) onCopyValue() {
	m.copyToClipboard(func(node *zk.Node) string { return node.Value }, "node value")
}

func (m *ContextMenu) onCopyName() {
	m.copyToClipboard(func(node *zk.Node) string { return node.Name }, "node name")
}

func (m *ContextMenu) copyToClipboard(producer func(node *zk.Node) string, errMsg string) {
	treeSelection, _ := getNodesTreeView().GetSelection()
	node, _ := getTreeSelectedNode(treeSelection)
	if node != nil {
		err := clipboard.WriteAll(producer(node))
		if err != nil {
			log.WithError(err).Errorf("Failed to copy %v to clipboard", errMsg)
		}
	}
}

func (m *ContextMenu) enableMenu(enabled bool) {
	GetObject("contextMenu").(*gtk.Menu).SetSensitive(enabled)
}
