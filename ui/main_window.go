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

func InitMainWindow(gtkWindow *gtk.Window) {
	createNodeDlg = NewCreateNodeDlg(gtkWindow)
	nodeExportDlg = NewNodeExportDlg(gtkWindow)
	nodeAction = NewNodeAction()
	contextMenu = NewContextMenu()
	notebook = NewNotebook()

	initNodeTree()
	initMainMenu(gtkWindow)
	initCssProvider()

	gtkWindow.SetTitle("Zooverseer")
	gtkWindow.ShowAll()
}

func initMainMenu(mainWindow *gtk.Window) {
	GetObject("menuConnect").(*gtk.MenuItem).Connect("activate", func() {
		connDialog := GetObject("connDialog").(*gtk.Dialog)
		connDialog.Show()
	})
	GetObject("menuExit").(*gtk.MenuItem).Connect("activate", func() {
		zk.CachingRepo.Close()
		mainWindow.Close()
	})

	GetObject("menuDisconnect").(*gtk.MenuItem).Connect("activate", func() {
		zk.Reset()
		ClearNodeTree()
	})

	GetObject("menuAdd").(*gtk.MenuItem).Connect("activate", contextMenu.onAddNewNode)
	GetObject("menuCopyValue").(*gtk.MenuItem).Connect("activate", contextMenu.onCopyValue)
	GetObject("menuExportNode").(*gtk.MenuItem).Connect("activate", contextMenu.onExportNode)
	GetObject("menuDeleteNode").(*gtk.MenuItem).Connect("activate", contextMenu.onDeleteNode)
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
