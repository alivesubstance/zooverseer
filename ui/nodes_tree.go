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
	log "github.com/sirupsen/logrus"
	"unsafe"
)

var (
	NodeTreeStore    *gtk.TreeStore
	ZkPathByTreePath = make(map[string]string)
	ZkRepo           = zk.CachingRepository{}
)

func InitNodeTree() {
	treeView := getObject("nodesTreeView").(*gtk.TreeView)
	treeView.AppendColumn(createTextColumn("Node", core.NodeColumn))
	treeView.Connect("test-expand-row", onTestExpandRow)
	treeView.Connect("button-press-event", onButtonPressEvent)

	treeSelection, _ := treeView.GetSelection()
	treeSelection.SetMode(gtk.SELECTION_SINGLE)
	treeSelection.Connect("changed", onTreeRowSelected)

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
	ZkRepo.InvalidateAll()
	ZkPathByTreePath = make(map[string]string)
}

func ShowTreeRootNodes() {
	rootChildren, err := ZkRepo.GetRootNodeChildren(GetSelectedConn())
	if err != nil {
		log.WithError(err).Error("Failed to get read ZK root node")
	}

	// add root children to tree
	for _, rootChild := range rootChildren {
		addSubRow(nil, rootChild)
	}
}

func onTreeRowSelected(treeSelection *gtk.TreeSelection) {
	model, iter, ok := treeSelection.GetSelected()
	if ok {
		treePath, err := model.(*gtk.TreeModel).GetPath(iter)
		if err != nil {
			log.WithError(err).Errorf("Could not get path from model: %s\n", treePath)
			return
		}
		zkPath := ZkPathByTreePath[treePath.String()]
		log.Tracef("Selected tree path: %s", zkPath)

		node, _ := ZkRepo.GetValue(zkPath, GetSelectedConn())
		if node == nil {
			log.Warnf("Value nil for %s", zkPath)
		}
		setNodeValue(node)
	}
}

func onTestExpandRow(treeView *gtk.TreeView, parentIter *gtk.TreeIter, treePath *gtk.TreePath) {
	removeRowChildren(parentIter, treePath)

	//TODO use go subroutine with channel in order not to freeze UI
	//TODO add spinner in case of long running function
	parentValue, _ := NodeTreeStore.GetValue(parentIter, core.NodeColumn)
	parentGoValue, _ := parentValue.GoValue()
	log.Debugf("Add %s children", parentGoValue)

	zkPath := ZkPathByTreePath[treePath.String()]
	children, err := ZkRepo.GetChildren(zkPath, GetSelectedConn())
	if err != nil {
		//TODO show error dialog
		log.WithError(err).Errorf("Failed to get data and children for %s", zkPath)
	}

	//setNodeValue(node)
	for i := range children {
		addSubRow(parentIter, children[i])
	}
}

func removeRowChildren(parentIter *gtk.TreeIter, treePath *gtk.TreePath) {
	parentValue, _ := NodeTreeStore.GetValue(parentIter, core.NodeColumn)
	parentGoValue, _ := parentValue.GoValue()
	log.Debugf("Remove %s children", parentGoValue)

	hasChildren := NodeTreeStore.IterHasChild(parentIter)
	if hasChildren {
		childrenNum := NodeTreeStore.IterNChildren(parentIter)
		log.Debugf("Remove %v children for %s", childrenNum, parentGoValue)

		for {
			var child gtk.TreeIter
			ok := NodeTreeStore.IterChildren(parentIter, &child)
			if ok {
				childValue, _ := NodeTreeStore.GetValue(&child, core.NodeColumn)
				childGoValue, _ := childValue.GoValue()
				log.Debugf("Remove child %s at parent %s", childGoValue, parentGoValue)

				childrenRemoved := NodeTreeStore.Remove(&child)
				if !childrenRemoved {
					break
				}
			}
		}
	} else {
		log.Debugf("Row at path %s and value %s has no children", treePath, parentGoValue)
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

func setNodeValue(node *zk.Node) {
	nodeDataTextView := getObject("nodeDataTextView").(*gtk.TextView)
	textBuffer, err := nodeDataTextView.GetBuffer()
	util.CheckErrorWithMsg("Failed to get text buffer", err)

	textBuffer.SetText(node.Value)
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

func onButtonPressEvent(b *gtk.TreeView, e *gdk.Event) {
	if isMouse2ButtonClicked(e) {
		log.Tracef("Mouse button 2 pressed")

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
