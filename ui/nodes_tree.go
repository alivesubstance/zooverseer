package ui

// #cgo pkg-config: gdk-3.0 glib-2.0 gobject-2.0
// #include <gdk/gdk.h>
// #include "/home/mirian/code/go/src/github.com/gotk3/gotk3/gdk/gdk.go.h"
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
	"unsafe"
)

var (
	nodeTreeStore    *gtk.TreeStore
	ZkPathByTreePath = make(map[string]string)
)

func initNodeTree() {
	nodesTreeView := getNodesTreeView()
	nodesTreeView.AppendColumn(createTextColumn("Node", core.NodeColumn))
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
	ZkCachingRepo.InvalidateAll()
	ZkPathByTreePath = make(map[string]string)
}

func ShowTreeRootNodes() error {
	rootNode, err := ZkCachingRepo.GetRootNode()
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
		log.Tracef("Nothing selected")
		// there is no error. nothing to selected. happens f.i. after node deletion
		return nil, nil
	}

	node, err := ZkCachingRepo.GetValue(zkPath)
	if err != nil {
		return nil, errors.Wrapf(err, "Value nil for %s", zkPath)
	}
	return node, nil
}

func getTreeSelectedZkPath(treeSelection *gtk.TreeSelection) (string, error) {
	model, iter, ok := treeSelection.GetSelected()
	if !ok {
		nodeAction.enableButtons(false)
		contextMenu.enableButtons(false)
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
	nodeAction.enableButtons(true)
	contextMenu.enableButtons(true)

	node, err := getTreeSelectedNode(treeSelection)
	if err != nil {
		log.WithError(err).Error("Failed to get tree selected node")
		dialog := createWarnDialog(getMainWindow(), "Unable to fetch node value: "+errors.Cause(err).Error())
		dialog.Run()
		dialog.Hide()
		return
	}

	if node != nil {
		notebook.showPage(node, getNotebook().GetCurrentPage())
	}
}

//todo test on 100+ children. works slowly
func onExpandRow(treeView *gtk.TreeView, parentIter *gtk.TreeIter, treePath *gtk.TreePath) {
	removeRowChildren(parentIter)

	//todo use go subroutine with channel in order not to freeze UI
	//todo add spinner in case of long running function
	parentValue, _ := nodeTreeStore.GetValue(parentIter, core.NodeColumn)
	parentGoValue, _ := parentValue.GoValue()

	zkPath := ZkPathByTreePath[treePath.String()]
	children, err := ZkCachingRepo.GetChildren(zkPath)
	if err != nil {
		//todo show error dialog
		log.WithError(err).Errorf("Failed to get data and children for %s", zkPath)
	}

	log.Tracef("Add %v children at %s", len(children), parentGoValue)
	for i := range children {
		addSubRow(parentIter, children[i])
	}
}

func removeRowChildren(parentIter *gtk.TreeIter) {
	parentTreePath, _ := nodeTreeStore.GetPath(parentIter)
	children := nodeTreeStore.IterNChildren(parentIter)

	parentZkPath := ZkPathByTreePath[parentTreePath.String()]
	log.Tracef("Remove %v children at %s", children, parentZkPath)
	ZkCachingRepo.Invalidate(parentZkPath)

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
	//setNodeValue(child)

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

func setNodeValue(node *zk.Node) {
	nodeDataTextView := getObject("nodeDataTextView").(*gtk.TextView)
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
	if isMouse2ButtonClicked(e) {
		menu := getObject("popupMenu").(*gtk.Menu)
		menu.ShowAll()
		menu.PopupAtPointer(e)
	}
}

// Go GTK3 adapter not fully support GTK3 API.
// But it expose native GDK objects to be used to extend API.
// This method determine when second mouse button clicked by
// analyzing button field of native C GdkEventButton struct.
func isMouse2ButtonClicked(e *gdk.Event) bool {
	event := &gdk.EventKey{Event: e}
	mouseButton := (*C.GdkEventButton)(unsafe.Pointer(event.Event.GdkEvent)).button
	return uint(mouseButton) == 3
}

func getNodesTreeView() *gtk.TreeView {
	return getObject("nodesTreeView").(*gtk.TreeView)
}

func refreshNode(zkPath string) {
	var treePath string
	for treePathKey, cachedZkPath := range ZkPathByTreePath {
		if cachedZkPath == zkPath {
			treePath = treePathKey
			break
		}
	}

	parentTreeIter, _ := nodeTreeStore.GetIterFromString(treePath)
	parentTreePath, _ := nodeTreeStore.GetPath(parentTreeIter)
	onExpandRow(getNodesTreeView(), parentTreeIter, parentTreePath)
	getNodesTreeView().ExpandToPath(parentTreePath)

	node, _ := ZkCachingRepo.GetValue(zkPath)
	notebook.showPage(node, PageData)
}

func deleteSelectedNode() {
	treeSelection, _ := getNodesTreeView().GetSelection()
	zkPath, _ := getTreeSelectedZkPath(treeSelection)
	//dialog := createConfirmDialog(getConnDialog(), "Are you sure you want to delete "+gopath.Base(zkPath)+"?")
	//resp := dialog.Run()
	//dialog.Hide()
	//if resp == gtk.RESPONSE_YES {
	node, _ := ZkCachingRepo.GetValue(zkPath)
	err := ZkCachingRepo.Delete(zkPath, node)
	if err != nil {
		msg := "Unable to delete node"
		log.WithError(err).Error(msg)
		warnDlg := createWarnDialog(getMainWindow(), msg+errors.Cause(err).Error())
		warnDlg.Run()
		warnDlg.Hide()
		return
	}

	parentZkPath := gopath.Dir(zkPath)
	refreshNode(parentZkPath)
	//}
}
