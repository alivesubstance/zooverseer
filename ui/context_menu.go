package ui

import (
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
	//dlgFileSave    *gtk.FileChooserDialog
}

func NewContextMenu() *ContextMenu {
	contextMenu := &ContextMenu{}
	contextMenu.gtkMenu = GetObject("popupMenu").(*gtk.Menu)

	contextMenu.addItem = GetObject("contextMenuAdd").(*gtk.MenuItem)
	contextMenu.addItem.Connect("activate", contextMenu.onAddNewNode)

	contextMenu.copyValueItem = GetObject("contextMenuCopyValue").(*gtk.MenuItem)
	contextMenu.copyValueItem.Connect("activate", contextMenu.onCopyValue)

	//contextMenu.copyNameItem = GetObject("contextMenuCopyName").(*gtk.MenuItem)
	//contextMenu.copyNameItem.Connect("activate", onCopyValue)
	//
	//contextMenu.copyPathItem = GetObject("contextMenuCopyPath").(*gtk.MenuItem)
	//contextMenu.copyPathItem.Connect("activate", onCopyValue)
	//
	//contextMenu.pasteValueItem = GetObject("contextMenuPasteValue").(*gtk.MenuItem)
	//contextMenu.pasteValueItem.Connect("activate", onCopyValue)
	//
	//contextMenu.renameItem = GetObject("contextMenuRename").(*gtk.MenuItem)
	//contextMenu.renameItem.Connect("activate", onCopyValue)
	//
	//contextMenu.deleteItem = GetObject("contextMenuDeleteNode").(*gtk.MenuItem)
	//contextMenu.deleteItem.Connect("activate", onCopyValue)

	contextMenu.exportItem = GetObject("contextMenuExportNode").(*gtk.MenuItem)
	contextMenu.exportItem.Connect("activate", contextMenu.onExportNode)

	//contextMenu.spinner = GetObject("spinner").(*gtk.Spinner)

	contextMenu.enableMenu(false)
	enableSpinner(false)

	return contextMenu
}

func (m *ContextMenu) onAddNewNode() {
	nodeAction.onNodeCreateBtnClicked()
}

func (m *ContextMenu) onExportNode() {
	ExportSelectedNode()
}

func (m *ContextMenu) onCopyValue() {
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

func (m *ContextMenu) enableMenu(enabled bool) {
	m.gtkMenu.SetSensitive(enabled)
}
