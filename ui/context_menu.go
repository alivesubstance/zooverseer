package ui

import (
	"github.com/alivesubstance/zooverseer/core/zk"
	"github.com/alivesubstance/zooverseer/task"
	"github.com/atotto/clipboard"
	"github.com/gotk3/gotk3/gtk"
	log "github.com/sirupsen/logrus"
	"time"
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
	dlgFileSave    *gtk.FileChooserDialog
}

func NewContextMenu() *ContextMenu {
	contextMenu := &ContextMenu{}
	contextMenu.gtkMenu = getObject("popupMenu").(*gtk.Menu)

	contextMenu.addItem = getObject("contextMenuAdd").(*gtk.MenuItem)
	contextMenu.addItem.Connect("activate", contextMenu.onAddNewNode)

	contextMenu.copyValueItem = getObject("contextMenuCopyValue").(*gtk.MenuItem)
	contextMenu.copyValueItem.Connect("activate", contextMenu.onCopyValue)

	//contextMenu.copyNameItem = getObject("contextMenuCopyName").(*gtk.MenuItem)
	//contextMenu.copyNameItem.Connect("activate", onCopyValue)
	//
	//contextMenu.copyPathItem = getObject("contextMenuCopyPath").(*gtk.MenuItem)
	//contextMenu.copyPathItem.Connect("activate", onCopyValue)
	//
	//contextMenu.pasteValueItem = getObject("contextMenuPasteValue").(*gtk.MenuItem)
	//contextMenu.pasteValueItem.Connect("activate", onCopyValue)
	//
	//contextMenu.renameItem = getObject("contextMenuRename").(*gtk.MenuItem)
	//contextMenu.renameItem.Connect("activate", onCopyValue)
	//
	//contextMenu.deleteItem = getObject("contextMenuDeleteNode").(*gtk.MenuItem)
	//contextMenu.deleteItem.Connect("activate", onCopyValue)

	contextMenu.exportItem = getObject("contextMenuExportNode").(*gtk.MenuItem)
	contextMenu.exportItem.Connect("activate", contextMenu.onExportNode)

	//contextMenu.spinner = getObject("spinner").(*gtk.Spinner)

	contextMenu.dlgFileSave = getObject("dlgFileSave").(*gtk.FileChooserDialog)

	//overlay := getObject("nodeOverlay").(*gtk.Overlay)
	//overlay.AddOverlay()

	contextMenu.enableMenu(false)
	contextMenu.enableSpinner(false)

	return contextMenu
}

func (m *ContextMenu) onAddNewNode() {
	nodeAction.onNodeCreateBtnClicked()
}

func (m *ContextMenu) onExportNode() {
	treeSelection, _ := getNodesTreeView().GetSelection()
	zkPath, _ := getTreeSelectedZkPath(treeSelection)

	onError := func(err error) {
		CreateErrorDialog(GetMainWindow(), "Export from "+zkPath+" failed: "+err.Error())
	}
	onSuccess := func(tree interface{}) {
		log.Infof("Exported " + tree.(*zk.Node).Name)
		m.enableSpinner(false)
		//responseType := m.dlgFileSave.Run()
		//log.Infof("%v", responseType)
	}

	go task.CreateExportTask(zkPath, m.onStartTask, onError, onSuccess)
}

func (m *ContextMenu) onStartTask() {
	m.enableSpinner(true)

	time.Sleep(2 * time.Second)
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

func (m *ContextMenu) enableSpinner(enable bool) {
	box := getObject("mainWindow").(*gtk.Window)
	box.SetSensitive(!enable)

	if enable {
		//spinner, _ := gtk.SpinnerNew()
		//spinner.SetSensitive(true)
		//spinner.ShowAll()
		//spinner.Start()
		//
		//overlay := getObject("nodeOverlay").(*gtk.Overlay)
		//overlay.AddOverlay(spinner)
		//m.spinner.Show()
		//m.spinner.Start()
	} else {
		//overlay := getObject("nodeOverlay").(*gtk.Overlay)
		//overlay.rem
		//m.spinner.Stop()
		//m.spinner.Hide()
	}
}
