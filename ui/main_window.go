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
)

func InitMainWindow(mainWindow *gtk.Window) {
	createNodeDlg = NewCreateNodeDlg(mainWindow)
	nodeExportDlg = NewNodeExportDlg(mainWindow)
	nodeAction = NewNodeAction()
	contextMenu = NewContextMenu()

	initNodeTree()
	notebook.init()
	initMainMenu(mainWindow)
	initCssProvider()

	mainWindow.SetTitle("Zooverseer")
	mainWindow.ShowAll()
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
		zk.CachingRepo.Close()
		ClearNodeTree()
	})

	GetObject("menuExport").(*gtk.MenuItem).Connect("activate", func() {
		contextMenu.onExportNode()
	})
}

func initCssProvider() {
	providerNew, err := gtk.CssProviderNew()
	if err != nil {
		log.WithError(err).Errorf("Failed to create CSS provider")
	}

	err = providerNew.LoadFromPath(core.CssStyleFilePath)
	if err != nil {
		log.WithError(err).Errorf("Failed to load CSS styles")
	}

	defaultScreen, err := gdk.ScreenGetDefault()
	if err != nil {
		log.WithError(err).Errorf("Failed to get default screen")
	}

	gtk.AddProviderForScreen(defaultScreen, providerNew, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
}
