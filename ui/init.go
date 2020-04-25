package ui

import (
	"github.com/alivesubstance/zooverseer/core"
	"github.com/alivesubstance/zooverseer/core/zk"
	"github.com/alivesubstance/zooverseer/util"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	log "github.com/sirupsen/logrus"
)

// there are rumors that global variable is evil. why?
var (
	Builder *gtk.Builder
)

func OnAppActivate(app *gtk.Application) func() {
	return func() {
		log.Info("Reading glade file")
		builder, err := gtk.BuilderNewFromFile(core.GladeFilePath)
		util.CheckError(err)

		Builder = builder

		mainWindow := getObject("mainWindow").(*gtk.Window)
		InitMainWindow(mainWindow)

		InitConnDialog(mainWindow)

		app.AddWindow(mainWindow)
	}
}

func showConfirmDialog(parent gtk.IWindow, text string) gtk.ResponseType {
	dialog := gtk.MessageDialogNew(parent, gtk.DIALOG_MODAL, gtk.MESSAGE_QUESTION, gtk.BUTTONS_YES_NO, text)
	return dialog.Run()
}

func getObject(objectName string) glib.IObject {
	object, err := Builder.GetObject(objectName)
	util.CheckError(err)

	return object
}

func getTreeSelectedValue(treeSelection *gtk.TreeSelection) *zk.Node {
	model, iter, ok := treeSelection.GetSelected()
	if ok {
		treePath, err := model.(*gtk.TreeModel).GetPath(iter)
		if err != nil {
			log.WithError(err).Errorf("Could not get path from model: %s\n", treePath)
			return nil
		}

		zkPath := ZkPathByTreePath[treePath.String()]
		log.Tracef("Selected tree path: %s", zkPath)

		node, _ := ZkRepo.GetValue(zkPath, GetSelectedConn())
		if node == nil {
			log.Errorf("Value nil for %s", zkPath)
		}
		return node
	}

	return nil
}
