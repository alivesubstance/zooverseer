package ui

import (
	"github.com/alivesubstance/zooverseer/core"
	"github.com/alivesubstance/zooverseer/core/zk"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	log "github.com/sirupsen/logrus"
)

var (
	createNodeDlg *CreateNodeDlg
	nodeExportDlg *ExportResultDlg
	nodeAction    *NodeAction
	contextMenu   *ContextMenu
	notebook      *Notebook
)

type MainWindow struct {
	gtkWindow      *gtk.Window
	menuAdd        *gtk.MenuItem
	menuCopyValue  *gtk.MenuItem
	menuCopyName   *gtk.MenuItem
	menuExportNode *gtk.MenuItem
	menuDeleteNode *gtk.MenuItem
}

func NewMainWindow() *MainWindow {
	mainWindow := &MainWindow{}
	mainWindow.gtkWindow = GetObject("mainWindow").(*gtk.Window)

	createNodeDlg = NewCreateNodeDlg(mainWindow.gtkWindow)
	nodeExportDlg = NewNodeExportDlg(mainWindow.gtkWindow)
	nodeAction = NewNodeAction()
	contextMenu = NewContextMenu()
	notebook = NewNotebook()

	initNodeTree()
	initMainMenu(mainWindow)
	initCssProvider()

	mainWindow.gtkWindow.SetTitle(core.Config.AppTitle)
	mainWindow.gtkWindow.Show()

	mainWindow.enableEditActions(false)

	return mainWindow
}

func initMainMenu(mainWindow *MainWindow) {
	GetObject("menuConnect").(*gtk.MenuItem).Connect("activate", func() {
		connDlg.dlg.ShowAll()
	})
	GetObject("menuExit").(*gtk.MenuItem).Connect("activate", func() {
		zk.CachingRepo.Close()
		mainWindow.gtkWindow.Close()
	})

	GetObject("menuDisconnect").(*gtk.MenuItem).Connect("activate", func() {
		zk.Reset()
		ClearNodeTree()
		mainWindow.gtkWindow.SetTitle(core.Config.AppTitle)
	})

	mainWindow.menuAdd = GetObject("menuAdd").(*gtk.MenuItem)
	mainWindow.menuAdd.Connect("activate", contextMenu.onAddNewNode)

	mainWindow.menuCopyValue = GetObject("menuCopyValue").(*gtk.MenuItem)
	mainWindow.menuCopyValue.Connect("activate", contextMenu.onCopyValue)

	mainWindow.menuCopyName = GetObject("menuCopyName").(*gtk.MenuItem)
	mainWindow.menuCopyName.Connect("activate", contextMenu.onCopyName)

	mainWindow.menuExportNode = GetObject("menuExportNode").(*gtk.MenuItem)
	mainWindow.menuExportNode.Connect("activate", contextMenu.onExportNode)

	mainWindow.menuDeleteNode = GetObject("menuDeleteNode").(*gtk.MenuItem)
	mainWindow.menuDeleteNode.Connect("activate", contextMenu.onDeleteNode)
}

func initCssProvider() {
	providerNew, err := gtk.CssProviderNew()
	if err != nil {
		log.WithError(err).Errorf("Failed to create CSS provider")
	}

	err = providerNew.LoadFromPath(core.Config.CssStyleFilePath)
	if err != nil {
		log.WithError(err).Errorf("Failed to load CSS styles")
	}

	defaultScreen, err := gdk.ScreenGetDefault()
	if err != nil {
		log.WithError(err).Errorf("Failed to get default screen")
	}

	gtk.AddProviderForScreen(defaultScreen, providerNew, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
}

func (w *MainWindow) enableEditActions(enabled bool) {
	w.menuAdd.SetSensitive(enabled)
	w.menuCopyValue.SetSensitive(enabled)
	w.menuCopyValue.SetSensitive(enabled)
	w.menuExportNode.SetSensitive(enabled)
	w.menuDeleteNode.SetSensitive(enabled)
}
