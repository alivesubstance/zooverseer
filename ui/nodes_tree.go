package ui

import (
	"fmt"
	"github.com/alivesubstance/zooverseer/core"
	"github.com/alivesubstance/zooverseer/core/zk"
	"github.com/alivesubstance/zooverseer/util"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	log "github.com/sirupsen/logrus"
)

var (
	NodeTreeStore    *gtk.TreeStore
	ZkPathByTreePath = make(map[string]string)
)

func InitNodesTree() {
	treeView := getObject("nodesTreeView").(*gtk.TreeView)
	treeView.AppendColumn(createTextColumn("Node", core.NodeColumn))
	treeView.Connect("row-expanded", onTreeRowExpanded)

	treeSelection, _ := treeView.GetSelection()
	treeSelection.SetMode(gtk.SELECTION_SINGLE)
	treeSelection.Connect("changed", onTreeRowChanged)

	nodesTreeView := getObject("nodesTreeView").(*gtk.TreeView)
	newTreeStore, err := gtk.TreeStoreNew(glib.TYPE_STRING)
	util.CheckError(err)

	NodeTreeStore = newTreeStore
	nodesTreeView.SetModel(NodeTreeStore)

	//TODO test. remove once conn dialog will be used
	ShowTreeRootNodes()
}

func ClearNodeTree() {
	NodeTreeStore.Clear()
}

func ShowTreeRootNodes() {
	rootChildren, err := zk.GetRootNodeChildren(GetSelectedConn())
	if err != nil {
		log.WithError(err).Error("Failed to get read ZK root node")
	}

	// add root children to tree
	for _, rootChild := range rootChildren {
		addSubRow(nil, &rootChild)
	}
}

func onTreeRowChanged(treeSelection *gtk.TreeSelection) {
	model, iter, ok := treeSelection.GetSelected()
	if ok {
		selectedPath, err := model.(*gtk.TreeModel).GetPath(iter)
		if err != nil {
			util.CheckErrorWithMsg(fmt.Sprintf("Could not get path from model: %s\n", selectedPath), err)
			return
		}
		log.Infof("Selected path: %s\n", selectedPath)

		value, _ := model.(*gtk.TreeModel).GetValue(iter, core.NodeColumn)
		valueStr, _ := value.GetString()
		log.Info("Selected value " + valueStr)
	}
}

func onTreeRowExpanded(treeView *gtk.TreeView, treeIter *gtk.TreeIter, treePath *gtk.TreePath) {
	//TODO use go subroutine with channel in order not to freeze UI
	//TODO add spinner in case of long running function

	treePathStr := treePath.String()
	zkPath := ZkPathByTreePath[treePathStr]
	if zkPath == core.NodeRootName {
		return
	}

	node, err := zk.Get(zkPath, GetSelectedConn())
	if err != nil {
		//TODO show error dialog
		log.WithError(err).Errorf("Failed to get data and children for %s", zkPath)
	}

	setNodeName(treeIter, node.Name)
	for _, child := range node.Children {
		addSubRow(treeIter, &child)
	}
}

func addSubRow(parentIter *gtk.TreeIter, child *zk.Node) {
	childIter := NodeTreeStore.Append(parentIter)
	setNodeName(childIter, child.Name)
	setNodeValue(child)

	if child.Meta.NumChildren > 0 {
		// add dummy node value in order to force GtkTreeView show expander icon
		dummyChildIter := NodeTreeStore.Append(childIter)
		setNodeName(dummyChildIter, core.NodeDummy)
	}

	parentZkPath := ""
	if parentIter != nil {
		parentTreePath := getTreePath(parentIter)
		parentZkPath = ZkPathByTreePath[parentTreePath]
	}

	childTreePath := getTreePath(childIter)
	childZkPath := fmt.Sprintf("%s/%s", parentZkPath, child.Name)
	ZkPathByTreePath[childTreePath] = childZkPath
}

func setNodeValue(child *zk.Node) {
	nodeDataTextView := getObject("nodeDataTextView").(*gtk.TextView)
	textBuffer, err := nodeDataTextView.GetBuffer()
	util.CheckErrorWithMsg("Faield to get text buffer", err)

	textBuffer.SetText(child.Value)
}

func getTreePath(iter *gtk.TreeIter) string {
	path, err := NodeTreeStore.GetPath(iter)
	util.CheckErrorWithMsg(fmt.Sprintf("Failed to get path for %s\n", iter), err)
	return path.String()
}

func setNodeName(treeIter *gtk.TreeIter, value string) {
	err := NodeTreeStore.SetValue(treeIter, core.NodeColumn, value)
	if err != nil {
		path, err := NodeTreeStore.GetPath(treeIter)
		util.CheckError(err)

		log.Panic("Unable set value ["+value+"] for ["+path.String()+"]", err)
	}
}

func createTextColumn(title string, id int) *gtk.TreeViewColumn {
	// In this column we want to show text, hence create a text renderer
	cellRenderer, err := gtk.CellRendererTextNew()
	util.CheckErrorWithMsg("Unable to create text cell renderer", err)

	// Tell the renderer where to pick input from. Text renderer understands
	// the "text" property.
	column, err := gtk.TreeViewColumnNewWithAttribute(title, cellRenderer, "text", id)
	util.CheckErrorWithMsg("Unable to create cell column:", err)

	return column
}
