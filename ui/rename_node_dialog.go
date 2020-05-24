package ui

import (
	"github.com/gotk3/gotk3/gtk"
)

type RenameNodeDlg struct {
	nameEntry *gtk.Entry
	dlg       *gtk.Dialog
}

func NewRenameNodeDlg(mainWindow *gtk.Window) *RenameNodeDlg {
	renameNodeDlg := RenameNodeDlg{}
	renameNodeDlg.nameEntry = getObject("renameNodeDlgNameEntry").(*gtk.Entry)
	renameNodeDlg.dlg = getObject("createNodeDlg").(*gtk.Dialog)
	createNodeDlg.dlg.SetTransientFor(mainWindow)

	getObject("renameNodeDlgOkBtn").(*gtk.Button).Connect("clicked", renameNodeDlg.onOkBtnClicked)
	getObject("renameNodeDlgCancelBtn").(*gtk.Button).Connect("clicked", renameNodeDlg.onCancelBtnClicked)

	return &renameNodeDlg
}

func (c *RenameNodeDlg) onOkBtnClicked() {
	//newName, _ := c.nameEntry.GetText()
	//if len(newName) == 0 {
	//	return
	//}
	//
	//selection, _ := getNodesTreeView().GetSelection()
	//parentZkPath, _ := getTreeSelectedZkPath(selection)
	//
	//c.hide()
	//err := ZkCachingRepo.SaveChild(parentZkPath, newName, c.getAcl())
	//if err != nil {
	//	msg := "Unable to create node: " + newName
	//	log.WithError(err).Warn(msg)
	//	dlg := CreateWarnDialog(GetMainWindow(), msg+"\n"+err.Error())
	//	dlg.Run()
	//	dlg.Hide()
	//} else {
	//	refreshNode(parentZkPath)
	//}
}

func (c *RenameNodeDlg) onCancelBtnClicked() {
	c.hide()
}

func (c *RenameNodeDlg) showAll() {
	c.nameEntry.SetText("")
	c.nameEntry.GrabFocus()
	c.dlg.ShowAll()
}

func (c *RenameNodeDlg) hide() {
	c.dlg.Hide()
}
