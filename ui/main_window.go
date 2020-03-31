package ui

// #cgo pkg-config: gdk-3.0 glib-2.0 gobject-2.0
// #include <gdk/gdk.h>
// #include "/home/mirian/code/go/src/github.com/gotk3/gotk3/gdk/gdk.go.h"
import "C"
import (
	"github.com/alivesubstance/zooverseer/core/zk"
	"github.com/alivesubstance/zooverseer/util"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	log "github.com/sirupsen/logrus"
	"unsafe"
)

// ID to access the tree view columns by
const (
	NodeColumn    = 0
	NodeRootPath  = "0"
	NodeRootValue = "/"
)

func InitMainWindow(mainWindow *gtk.Window) {
	treeView := getObject("nodesTreeView").(*gtk.TreeView)
	treeView.Connect("row-expanded", onTreeRowExpanded)

	mainWindow.ShowAll()
}

func onTreeRowExpanded(treeView *gtk.TreeView, treeIter *gtk.TreeIter, treePath *gtk.TreePath) {
	//TODO use go subroutine with channel in order not to freeze UI
	//TODO add spinner in case of long running function

	nodeChannel := make(chan zk.Node)

	//TODO FUCK! tree path ash for 0:1 but zk need /env/sandbox-pleeco
	treePathStr := treePath.String()
	err := zk.Get("/", GetSelectedConn(), nodeChannel)
	if err != nil {
		//TODO show error dialog
		util.CheckErrorWithMsg("Failed to get data and children for ["+treePathStr+"]", err)
	}
	node := <-nodeChannel

	log.Debugln("Node " + node.Name + " has " + string(len(node.Children)) + " children")

	nodesStore := getObject("nodesStore").(*gtk.TreeStore)
	setNodeValue(nodesStore, treeIter, node.Value)

	childIter := nodesStore.Append(treeIter)
	for _, child := range node.Children {
		log.Println("Add child " + child.Name)
		addSubRow(nodesStore, childIter, child.Name)
	}
}

func on_button_press_event(b *gtk.Window, e *gdk.Event) {
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
	event := &gdk.EventKey{e}
	mouseButton := (*C.GdkEventButton)(unsafe.Pointer(event.Event.GdkEvent)).button
	return uint(mouseButton) == 3
}

func ShowTreeRootData() {
	// for tree path see
	//https://developer.gnome.org/gtk3/stable/GtkTreeModel.html#gtk-tree-path-new-from-string
	rootTreePath, err := gtk.TreePathNewFromString(NodeRootPath)
	util.CheckErrorWithMsg("Failed to get root tree path", err)

	nodesTreeView := getObject("nodesTreeView").(*gtk.TreeView)
	nodesStore := getObject("nodesStore").(*gtk.TreeStore)

	// append root node
	rootIter := nodesStore.Append(nil)
	addSubRow(nodesStore, rootIter, NodeRootValue)

	onTreeRowExpanded(nodesTreeView, rootIter, rootTreePath)
}

func addSubRow(nodesStore *gtk.TreeStore, treeIter *gtk.TreeIter, value string) {
	childIter := nodesStore.Append(treeIter)
	setNodeValue(nodesStore, childIter, value)
}

func setNodeValue(nodesStore *gtk.TreeStore, treeIter *gtk.TreeIter, value string) {
	err := nodesStore.SetValue(treeIter, NodeColumn, value)
	if err != nil {
		path, err := nodesStore.GetPath(treeIter)
		util.CheckError(err)

		log.Panic("Unable set value ["+value+"] for ["+path.String()+"]", err)
	}
}

func ClearNodesTree() {
	nodesStore := getObject("nodesStore").(*gtk.TreeStore)
	nodesStore.Clear()
}
