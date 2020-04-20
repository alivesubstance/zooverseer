package ui

import (
	"github.com/atotto/clipboard"
	"github.com/gotk3/gotk3/gtk"
	log "github.com/sirupsen/logrus"
)

func InitContextMenu() {
	getObject("popupMenuCopyValue").(*gtk.MenuItem).Connect("activate", onPopupMenuCopyValue)
}

func onPopupMenuCopyValue() {
	treeView := getObject("nodesTreeView").(*gtk.TreeView)
	treeSelection, err := treeView.GetSelection()
	if err != nil {
		log.WithError(err).Errorf("Failed to get nodes tree selection")
	}
	node := getTreeSelectedValue(treeSelection)
	if node != nil {
		log.Tracef("Add text value to clipboard")
		err := clipboard.WriteAll(node.Value)
		if err != nil {
			log.WithError(err).Errorf("Failed to copy value to clipboard")
		}

	}

}
