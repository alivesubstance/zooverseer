package ui

import (
	"fmt"
	"github.com/alivesubstance/zooverseer/core"
	"github.com/alivesubstance/zooverseer/core/zk"
	"github.com/alivesubstance/zooverseer/util"
	"github.com/gotk3/gotk3/gtk"
	goZk "github.com/samuel/go-zookeeper/zk"
	log "github.com/sirupsen/logrus"
)

type CreateNodeDlg struct {
	dlg *gtk.Dialog
}

func NewCreateNodeDlg(mainWindow *gtk.Window) *CreateNodeDlg {
	createNodeDlg := CreateNodeDlg{}
	createNodeDlg.dlg = GetObject("createNodeDlg").(*gtk.Dialog)
	createNodeDlg.dlg.SetTransientFor(mainWindow)

	GetObject("createNodeDlgOkBtn").(*gtk.Button).Connect("clicked", createNodeDlg.onOkBtnClicked)
	GetObject("createNodeDlgCancelBtn").(*gtk.Button).Connect("clicked", createNodeDlg.onCancelBtnClicked)

	return &createNodeDlg
}

func (c *CreateNodeDlg) onOkBtnClicked() {
	nodeName, _ := c.getNameEntry().GetText()
	if len(nodeName) == 0 {
		return
	}

	selection, _ := getNodesTreeView().GetSelection()
	parentZkPath, _ := getTreeSelectedZkPath(selection)

	c.hide()
	err := zk.CachingRepo.SaveChild(parentZkPath, nodeName, c.getAcl())
	if err != nil {
		msg := "Unable to create node: " + nodeName
		log.WithError(err).Warn(msg)
		dialog := createWarnDialog(GetMainWindow(), msg+"\n"+err.Error())
		dialog.Run()
		dialog.Hide()
	} else {
		refreshNode(parentZkPath)
	}
}

func (c *CreateNodeDlg) onCancelBtnClicked() {
	c.hide()
}

func (c *CreateNodeDlg) setAcl(connInfo *core.ConnInfo) {
	//digest:someuser:hashedpw:crdwa
	//todo remove once ACL support will be full
	aclStr := ""
	if len(connInfo.User) != 0 && len(connInfo.Password) != 0 {
		aclStr = fmt.Sprintf("digest:%s:%s:%s", connInfo.User, util.Encrypt(connInfo.Password), "a")
	}
	c.getAclEntry().SetText(aclStr)
}

func (c *CreateNodeDlg) getAcl() []goZk.ACL {
	//digest:someuser:hashedpw:crdwa
	//todo remove once ACL support will be full
	connInfo := getSelectedConn()
	if len(connInfo.User) == 0 && len(connInfo.Password) == 0 {
		return core.AclWorldAnyone
	}
	return goZk.DigestACL(goZk.PermAll, connInfo.User, connInfo.Password)
}

func (c *CreateNodeDlg) showAll() {
	c.getNameEntry().SetText("")
	c.getNameEntry().GrabFocus()
	c.setAcl(getSelectedConn())
	c.dlg.ShowAll()
}

func (c *CreateNodeDlg) hide() {
	c.dlg.Hide()
}

func (c *CreateNodeDlg) getAclEntry() *gtk.Entry {
	return GetObject("createNodeDlgAclEntry").(*gtk.Entry)
}

func (c *CreateNodeDlg) getNameEntry() *gtk.Entry {
	return GetObject("createNodeDlgNameEntry").(*gtk.Entry)
}