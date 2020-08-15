package ui

import "C"
import (
	"fmt"
	"github.com/alivesubstance/zooverseer/core"
	"github.com/alivesubstance/zooverseer/core/zk"
	"github.com/alivesubstance/zooverseer/util"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	gopath "path"
)

var (
	nodeTreeStore    *gtk.TreeStore
	ZkPathByTreePath = make(map[string]string)
)

func initNodeTree() {
	nodesTreeView := getNodesTreeView()
	nodesTreeView.AppendColumn(createTextColumn("Data", core.NodeColumn))
	nodesTreeView.Connect("test-expand-row", onExpandRow)
	nodesTreeView.Connect("button-press-event", onMouseButtonPress)

	newTreeStore, err := gtk.TreeStoreNew(glib.TYPE_STRING)
	util.CheckError(err)

	nodeTreeStore = newTreeStore
	nodesTreeView.SetModel(nodeTreeStore)

	treeSelection, _ := nodesTreeView.GetSelection()
	treeSelection.SetMode(gtk.SELECTION_SINGLE)
	treeSelection.Connect("changed", onTreeRowSelected)
}

func ClearNodeTree() {
	nodeTreeStore.Clear()
	zk.CachingRepo.InvalidateAll()
	ZkPathByTreePath = make(map[string]string)
}

func ShowTreeRootNodes() error {
	rootNode, err := zk.CachingRepo.GetRootNode()
	if err == nil {
		rootTreeIter := addSubRow(nil, rootNode)
		rootTreePath, _ := nodeTreeStore.GetPath(rootTreeIter)
		getNodesTreeView().ExpandToPath(rootTreePath)
	}
	return err
}

func getTreeSelectedNode(treeSelection *gtk.TreeSelection) (*zk.Node, error) {
	zkPath, err := getTreeSelectedZkPath(treeSelection)
	if len(zkPath) == 0 {
		if err != nil {
			return nil, err
		}
		// there is no error. nothing to selected. happens f.i. after node deletion
		return nil, nil
	}

	node, err := zk.CachingRepo.Get(zkPath)
	if err != nil {
		return nil, errors.Wrapf(err, "Value nil for %s", zkPath)
	}
	return node, nil
}

// todo consider not to return error and show error dlg or even panic
func getTreeSelectedZkPath(treeSelection *gtk.TreeSelection) (string, error) {
	model, iter, ok := treeSelection.GetSelected()
	if !ok {
		enableRowActions(false)

		return "", nil
	}

	treePath, err := model.(*gtk.TreeModel).GetPath(iter)
	if err != nil {
		return "", errors.Wrapf(err, "Could not get path from model: %s\n", treePath)
	}

	zkPath := ZkPathByTreePath[treePath.String()]

	return zkPath, nil
}

func onTreeRowSelected(treeSelection *gtk.TreeSelection) {
	// reset node value
	drawNodeValue(&zk.Node{Value: ""})
	enableRowActions(true)

	node, err := getTreeSelectedNode(treeSelection)
	if err != nil {
		log.WithError(err).Error("Failed to get tree selected node")
		dialog := createWarnDialog(mainWindow.gtkWindow, "Unable to fetch node value: "+errors.Cause(err).Error())
		dialog.Run()
		dialog.Hide()
		return
	}

	if node != nil {
		notebook.showPage(node, getNotebook().GetCurrentPage())
	}
}

func enableRowActions(enabled bool) {
	nodeAction.enableButtons(enabled)
	contextMenu.enableMenu(enabled)
	mainWindow.enableEditActions(enabled)
}

//todo test on 100+ children. works slowly
func onExpandRow(treeView *gtk.TreeView, parentIter *gtk.TreeIter, treePath *gtk.TreePath) {
	removeRowChildren(parentIter)

	parentValue, _ := nodeTreeStore.GetValue(parentIter, core.NodeColumn)
	parentGoValue, _ := parentValue.GoValue()

	zkPath := ZkPathByTreePath[treePath.String()]
	node, err := zk.CachingRepo.Get(zkPath)
	if err != nil {
		//todo show error dlg
		log.WithError(err).Errorf("Failed to get data and children for %s", zkPath)
	}

	log.Tracef("Add %v children at %s", len(node.Children), parentGoValue)
	for i := range node.Children {
		addSubRow(parentIter, node.Children[i])
	}
}

func removeRowChildren(parentIter *gtk.TreeIter) {
	parentTreePath, _ := nodeTreeStore.GetPath(parentIter)
	children := nodeTreeStore.IterNChildren(parentIter)

	parentZkPath := ZkPathByTreePath[parentTreePath.String()]
	log.Tracef("Remove %v children at %s", children, parentZkPath)

	hasChildren := nodeTreeStore.IterHasChild(parentIter)
	if hasChildren {
		for {
			var child gtk.TreeIter
			ok := nodeTreeStore.IterChildren(parentIter, &child)
			if ok {
				childrenRemoved := nodeTreeStore.Remove(&child)
				if childrenRemoved {
					childTreePath, _ := nodeTreeStore.GetPath(&child)
					childTreePathStr := childTreePath.String()
					delete(ZkPathByTreePath, childTreePathStr)
				} else {
					break
				}
			}
		}
	}
}

func addSubRow(parentIter *gtk.TreeIter, child *zk.Node) *gtk.TreeIter {
	childIter := nodeTreeStore.Append(parentIter)
	setNodeName(childIter, child.Name)

	if child.Meta.NumChildren > 0 {
		// add dummy node value in order to force GtkTreeView show expander icon
		dummyChildIter := nodeTreeStore.Append(childIter)
		setNodeName(dummyChildIter, core.NodeDummy)
	}

	childZkPath := core.NodeRootName
	childTreePath := getZkTreePath(childIter)
	if parentIter != nil {
		parentTreePath := getZkTreePath(parentIter)
		parentZkPath := ZkPathByTreePath[parentTreePath]
		childZkPath = fmt.Sprintf("%s/%s", parentZkPath, child.Name)
		if parentZkPath == core.NodeRootName {
			// do not add leading '/' in case of top level(direct root child) nodes
			childZkPath = fmt.Sprintf("/%s", child.Name)
		}
	}

	log.Tracef("Add %s", childZkPath)
	ZkPathByTreePath[childTreePath] = childZkPath

	return childIter
}

func drawNodeValue(node *zk.Node) {
	nodeDataTextView := GetObject("nodeDataTextView").(*gtk.TextView)
	textBuffer, err := nodeDataTextView.GetBuffer()
	util.CheckErrorWithMsg("Failed to get text buffer", err)

	textBuffer.SetText(node.Value)
}

func getZkTreePath(iter *gtk.TreeIter) string {
	path, err := nodeTreeStore.GetPath(iter)
	util.CheckErrorWithMsg(fmt.Sprintf("Failed to get path for %s\n", iter), err)
	return path.String()
}

func setNodeName(treeIter *gtk.TreeIter, value string) {
	err := nodeTreeStore.SetValue(treeIter, core.NodeColumn, value)
	if err != nil {
		path, err := nodeTreeStore.GetPath(treeIter)
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

func onMouseButtonPress(b *gtk.TreeView, e *gdk.Event) {
	eventButton := gdk.EventButtonNewFromEvent(e)
	if eventButton.Button() == gdk.BUTTON_SECONDARY {
		contextMenu.gtkMenu.ShowAll()
		contextMenu.gtkMenu.PopupAtPointer(e)
	}
}

func getNodesTreeView() *gtk.TreeView {
	return GetObject("nodesTreeView").(*gtk.TreeView)
}

func refreshNode(zkPath string) {
	treePath := findTreePathByZkPath(zkPath)

	parentTreeIter, _ := nodeTreeStore.GetIterFromString(treePath)
	parentTreePath, _ := nodeTreeStore.GetPath(parentTreeIter)

	parentZkPath := ZkPathByTreePath[parentTreePath.String()]
	// remove cached node value
	zk.CachingRepo.Invalidate(parentZkPath)
	// mimic user click expand row
	onExpandRow(getNodesTreeView(), parentTreeIter, parentTreePath)
	// expand tree to parent path
	getNodesTreeView().ExpandToPath(parentTreePath)

	node, _ := zk.CachingRepo.GetValue(zkPath)
	notebook.showPage(node, PageData)
}

func deleteSelectedNode() {
	treeSelection, _ := getNodesTreeView().GetSelection()
	zkPath, _ := getTreeSelectedZkPath(treeSelection)
	dlg := createConfirmDialog(mainWindow.gtkWindow, "Are you sure you want to delete "+gopath.Base(zkPath)+"?")
	resp := dlg.Run()
	dlg.Hide()

	if resp == gtk.RESPONSE_YES {
		node, _ := zk.CachingRepo.GetValue(zkPath)
		err := zk.CachingRepo.Delete(zkPath, node)
		if err != nil {
			msg := "Unable to delete node"
			log.WithError(err).Error(msg)
			warnDlg := createWarnDialog(mainWindow.gtkWindow, msg+errors.Cause(err).Error())
			warnDlg.Run()
			warnDlg.Hide()
			return
		}

		// remove cached tree path for a given zk path
		delete(ZkPathByTreePath, findTreePathByZkPath(zkPath))

		parentZkPath := gopath.Dir(zkPath)
		refreshNode(parentZkPath)
	}
}

func findTreePathByZkPath(zkPath string) string {
	for treePathKey, cachedZkPath := range ZkPathByTreePath {
		if cachedZkPath == zkPath {
			return treePathKey
		}
	}

	panic("Failed to find tree path by " + zkPath + " zk path")
}
