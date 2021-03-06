package ui

import "github.com/gotk3/gotk3/gtk"

type NodeAction struct {
	createBtn  *gtk.Button
	refreshBtn *gtk.Button
	deleteBtn  *gtk.Button
}

func NewNodeAction() *NodeAction {
	n := &NodeAction{}
	n.createBtn = GetObject("nodeCreateBtn").(*gtk.Button)
	n.createBtn.Connect("clicked", n.onNodeCreateBtnClicked)

	n.refreshBtn = GetObject("nodeRefreshBtn").(*gtk.Button)
	n.refreshBtn.Connect("clicked", n.onNodeRefreshBtnClicked)

	n.deleteBtn = GetObject("nodeDeleteBtn").(*gtk.Button)
	n.deleteBtn.Connect("clicked", n.onNodeDeleteBtnClicked)

	n.enableButtons(false)
	return n
}

func (n *NodeAction) enableButtons(enabled bool) {
	n.createBtn.SetSensitive(enabled)
	n.refreshBtn.SetSensitive(enabled)
	n.deleteBtn.SetSensitive(enabled)
}

func (n *NodeAction) onNodeCreateBtnClicked() {
	createNodeDlg.showAll()
}

func (n *NodeAction) onNodeRefreshBtnClicked() {
	selection, _ := getNodesTreeView().GetSelection()
	parentZkPath, _ := getTreeSelectedZkPath(selection)
	refreshNode(parentZkPath)
}

func (n *NodeAction) onNodeDeleteBtnClicked() {
	deleteSelectedNode()
}
