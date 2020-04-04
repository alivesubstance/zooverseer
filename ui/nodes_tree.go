package ui

import (
	"fmt"
	"github.com/alivesubstance/zooverseer/core/zk"
	"github.com/alivesubstance/zooverseer/util"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	log "github.com/sirupsen/logrus"
)

// ID to access the tree view columns by
const (
	NodeColumn    = 0
	NodeRootPath  = "0"
	NodeRootValue = "/"
)

var (
	TreeStore        *gtk.TreeStore
	NodeDummy        = "__dummy"
	ZkPathByTreePath = make(map[string]string)
)

func InitNodesTree() {
	treeView := getObject("nodesTreeView").(*gtk.TreeView)
	treeView.AppendColumn(createTextColumn("Node", NodeColumn))
	treeView.Connect("row-expanded", onTreeRowExpanded)

	treeSelection, _ := treeView.GetSelection()
	treeSelection.SetMode(gtk.SELECTION_SINGLE)
	treeSelection.Connect("changed", onTreeRowChanged)

	nodesTreeView := getObject("nodesTreeView").(*gtk.TreeView)
	newTreeStore, err := gtk.TreeStoreNew(glib.TYPE_STRING)
	util.CheckError(err)

	TreeStore = newTreeStore
	nodesTreeView.SetModel(TreeStore)

	// ---- TEST ----
	// append root node
	//rootIter := TreeStore.Append(nil)
	//// set value for root node(value it's a node name but not a value from ZK root node)
	//err = TreeStore.SetValue(rootIter, NodeColumn, NodeRootValue)
	//util.CheckError(err)
	//addChild(TreeStore, rootIter, "1")
	//ZkPathByTreePath["0"] = NodeRootValue

	ShowTreeRootNodes()
	//rootTreePath, err := gtk.TreePathNewFromString(NodeRootPath)
	//util.CheckErrorWithMsg("Failed to get root tree path", err)
	//onTreeRowExpanded(nodesTreeView, rootIter, rootTreePath)
	// --------------
}

func ClearNodesTree() {
	nodesStore := getObject("nodesStore").(*gtk.TreeStore)
	nodesStore.Clear()
}

func ShowTreeRootNodes() {
	// for tree path see
	//https://developer.gnome.org/gtk3/stable/GtkTreeModel.html#gtk-tree-path-new-from-string
	rootTreePath, err := gtk.TreePathNewFromString(NodeRootPath)
	util.CheckErrorWithMsg("Failed to get root tree path", err)

	// append root node
	rootIter := TreeStore.Append(nil)
	addChild(TreeStore, rootIter, NodeDummy)

	// set value for root node(value it's a node name but not a value from ZK root node)
	err = TreeStore.SetValue(rootIter, NodeColumn, NodeRootValue)
	util.CheckError(err)

	nodesTreeView := getObject("nodesTreeView").(*gtk.TreeView)
	onTreeRowExpanded(nodesTreeView, rootIter, rootTreePath)
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

		value, _ := model.(*gtk.TreeModel).GetValue(iter, NodeColumn)
		valueStr, _ := value.GetString()
		log.Info("Selected value " + valueStr)
	}
}

func onTreeRowExpanded(treeView *gtk.TreeView, treeIter *gtk.TreeIter, treePath *gtk.TreePath) {
	//TODO use go subroutine with channel in order not to freeze UI
	//TODO add spinner in case of long running function

	//nodeChannel := make(chan zk.Node)

	//TODO FUCK! tree path has form 0:1 but zk need /env/sandbox-pleeco
	treePathStr := treePath.String()
	node, err := zk.Get("/", GetSelectedConn() /*, nodeChannel*/)
	if err != nil {
		//TODO show error dialog
		util.CheckErrorWithMsg("Failed to get data and children for ["+treePathStr+"]", err)
	}
	//node := <-nodeChannel

	log.Info("Node " + node.Name + " has " + string(len(node.Children)) + " children")

	//setNodeValue(TreeStore, treeIter, node.Value)
	//
	//childIter := TreeStore.Append(treeIter)
	for _, child := range node.Children {
		addChild(TreeStore, treeIter, child.Name)
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

func addChild(treeStore *gtk.TreeStore, parent *gtk.TreeIter, childValue string) {
	childIter := treeStore.Append(parent)
	setNodeValue(treeStore, childIter, childValue)

	childPath := getTreePath(treeStore, childIter).String()
	parentPath := getTreePath(treeStore, parent).String()

	parentZkPath := ZkPathByTreePath[parentPath]

	ZkPathByTreePath[childPath] = fmt.Sprintf("%s/%s", parentZkPath, childPath)
}

func getTreePath(treeStore *gtk.TreeStore, iter *gtk.TreeIter) *gtk.TreePath {
	path, err := treeStore.GetPath(iter)
	util.CheckErrorWithMsg(fmt.Sprintf("Failed to get path for %s\n", iter), err)
	return path
}

func setNodeValue(nodeStore *gtk.TreeStore, treeIter *gtk.TreeIter, value string) {
	err := nodeStore.SetValue(treeIter, NodeColumn, value)
	if err != nil {
		path, err := nodeStore.GetPath(treeIter)
		util.CheckError(err)

		log.Panic("Unable set value ["+value+"] for ["+path.String()+"]", err)
	}
}
