package ui

import (
	"github.com/alivesubstance/zooverseer/core"
	"github.com/alivesubstance/zooverseer/core/zk"
	"github.com/gotk3/gotk3/gtk"
	goZk "github.com/samuel/go-zookeeper/zk"
	log "github.com/sirupsen/logrus"
)

type CreateNodeDlg struct {
	dlg                *gtk.Dialog
	publicNodeCheckBtn *gtk.CheckButton
	nameEntry          *gtk.Entry
	valueEntry         *gtk.Entry
}

func NewCreateNodeDlg(mainWindow *gtk.Window) *CreateNodeDlg {
	createNodeDlg := CreateNodeDlg{}
	createNodeDlg.dlg = GetObject("createNodeDlg").(*gtk.Dialog)
	createNodeDlg.dlg.SetTransientFor(mainWindow)

	createNodeDlg.nameEntry = GetObject("createNodeDlgNameEntry").(*gtk.Entry)
	createNodeDlg.valueEntry = GetObject("createNodeDlgValueEntry").(*gtk.Entry)
	createNodeDlg.publicNodeCheckBtn = GetObject("createNodeDlgPublicNodeCheckBtn").(*gtk.CheckButton)

	GetObject("createNodeDlgOkBtn").(*gtk.Button).Connect("clicked", createNodeDlg.onOkBtnClicked)
	GetObject("createNodeDlgCancelBtn").(*gtk.Button).Connect("clicked", createNodeDlg.onCancelBtnClicked)

	return &createNodeDlg
}

func (c *CreateNodeDlg) onOkBtnClicked() {
	nodeName, _ := c.nameEntry.GetText()
	if len(nodeName) == 0 {
		return
	}

	selection, _ := getNodesTreeView().GetSelection()
	parentZkPath, _ := getTreeSelectedZkPath(selection)
	nodeValue, _ := c.valueEntry.GetText()

	c.hide()
	child := &zk.Node{
		Name:  nodeName,
		Value: nodeValue,
		Acl:   c.getAcl(),
	}
	err := zk.CachingRepo.SaveChild(parentZkPath, child)
	if err != nil {
		msg := "Unable to create child: " + nodeName
		log.WithError(err).Warn(msg)
		dialog := createWarnDialog(mainWindow.gtkWindow, msg+"\n"+err.Error())
		dialog.Run()
		dialog.Hide()
	} else {
		refreshNode(parentZkPath)
	}
}

func (c *CreateNodeDlg) onCancelBtnClicked() {
	c.hide()
}

func (c *CreateNodeDlg) getAcl() []goZk.ACL {
	//digest:someuser:hashedpw:crdwa
	connInfo := getSelectedConn()
	if c.publicNodeCheckBtn.GetActive() {
		return core.AclWorldAnyone
	}
	return goZk.DigestACL(goZk.PermAll, connInfo.User, connInfo.Password)
}

func (c *CreateNodeDlg) showAll() {
	c.nameEntry.SetText("")
	c.nameEntry.GrabFocus()
	c.valueEntry.SetText("")

	connInfo := getSelectedConn()
	c.publicNodeCheckBtn.SetSensitive(len(connInfo.User) != 0 && len(connInfo.Password) != 0)
	c.publicNodeCheckBtn.SetActive(true)

	c.dlg.ShowAll()
}

func (c *CreateNodeDlg) hide() {
	c.dlg.Hide()
}
