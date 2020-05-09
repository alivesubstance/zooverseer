package ui

import (
	zk2 "github.com/alivesubstance/zooverseer/core/zk"
	"github.com/alivesubstance/zooverseer/util"
	"github.com/gotk3/gotk3/gtk"
)

const (
	PageData     = 0
	PageMetadata = 1
	PageAcl      = 2
)

// use these pair to address exact func owner
var notebook = Notebook{}

type Notebook struct{}

func (n *Notebook) init() {
	getObject("notebook").(*gtk.Notebook).Connect(
		"switch-page",
		func(notebook *gtk.Notebook, widget *gtk.Widget, page int) {
			n.onNotebookSwitchPage(notebook, widget, page)
		},
	)
}

func (n *Notebook) onNotebookSwitchPage(notebook *gtk.Notebook, widget *gtk.Widget, page int) {
	treeSelection, err := getNodesTreeView().GetSelection()
	if err != nil {
		createAndRunWarnDialog(getMainWindow(), "Unable to get tree selection: "+err.Error())
		return
	}

	node, _ := getTreeSelectedNode(treeSelection)
	if node == nil {
		// nothing to show. no node is selected
		return
	}

	n.showPage(node, page)
}

func (n *Notebook) showPage(node *zk2.Node, page int) {
	switch page {
	case PageData:
		n.showPageData(node)
	case PageMetadata:
		n.showPageMetadata(node)
	case PageAcl:
		n.showPageAcl()
	}
}

func (n *Notebook) showPageData(node *zk2.Node) {
	setNodeValue(node)
}

func (n *Notebook) showPageMetadata(node *zk2.Node) {
	meta := node.Meta
	getObject("czxidEntry").(*gtk.Entry).SetText(util.Int64ToStr(meta.Czxid))
	getObject("mzxidEntry").(*gtk.Entry).SetText(util.Int64ToStr(meta.Mzxid))
	getObject("ctimeEntry").(*gtk.Entry).SetText(util.MillisToTime(meta.Ctime).String())
	getObject("mtimeEntry").(*gtk.Entry).SetText(util.MillisToTime(meta.Mtime).String())
	getObject("versionEntry").(*gtk.Entry).SetText(util.Int32ToStr(meta.Version))
	getObject("cversionEntry").(*gtk.Entry).SetText(util.Int32ToStr(meta.Cversion))
	getObject("aversionEntry").(*gtk.Entry).SetText(util.Int32ToStr(meta.Aversion))
	getObject("ephemeralOwnerEntry").(*gtk.Entry).SetText(util.Int64ToStr(meta.EphemeralOwner))
	getObject("dataLengthEntry").(*gtk.Entry).SetText(util.Int32ToStr(meta.DataLength))
	getObject("numChildrenEntry").(*gtk.Entry).SetText(util.Int32ToStr(meta.NumChildren))
	getObject("pzxidEntry").(*gtk.Entry).SetText(util.Int64ToStr(meta.Pzxid))
}

func (n *Notebook) showPageAcl() {

}

func getNotebook() *gtk.Notebook {
	return getObject("notebook").(*gtk.Notebook)
}
