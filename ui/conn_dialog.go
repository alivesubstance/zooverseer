package ui

import (
	"fmt"
	"github.com/alivesubstance/zooverseer/core"
	"github.com/alivesubstance/zooverseer/util"
	"github.com/gotk3/gotk3/gtk"
	"log"
)

func InitConnDialog(mainWindow *gtk.Window) *gtk.Dialog {
	connPortEntry := getObject("connPortEntry").(*gtk.Entry)
	connPortEntry.SetWidthChars(10)

	connDialog := getObject("connDialog").(*gtk.Dialog)
	connDialog.SetTransientFor(mainWindow)

	connDialogCancelBtn := getObject("connDialogCancelBtn").(*gtk.Button)
	connDialogCancelBtn.Connect("clicked", onConnDialogCancelBtnClicked(connDialog))

	connAddBtn := getObject("connAddBtn").(*gtk.Button)
	connAddBtn.Connect("clicked", onConnAddBtnClicked)

	connBtn := getObject("connBtn").(*gtk.Button)
	connBtn.Connect("clicked", onConnBtnClicked(connDialog))

	initConnsListBox()

	connDialog.ShowAll()

	return connDialog
}

//TODO cache it for session
func GetSelectedConn() *core.JsonConnInfo {
	//TODO leave it for test
	return &core.JsonConnInfo{
		Host: "10.1.1.1",
		Port: 2181,
		//User:     "zookeeper",
		//Password: "z00k33p3r",
	}
	//connList := getObject("connList").(*gtk.ListBox)
	//child, err := connList.GetSelectedRow().GetChild()
	//util.CheckError(err)
	//
	//connName, _ := child.GetTooltipText()
	//connInfo, ok := ConnRepository.Find(connName)
	//if !ok {
	//	log.Panicf("'%s' connection setting not found. Should never happened", connName)
	//}
	//
	//return connInfo
}

func onConnListBoxRowSelected() {
	selectedConn := GetSelectedConn()
	getObject("connNameEntry").(*gtk.Entry).SetText(selectedConn.Name)
	getObject("connHostEntry").(*gtk.Entry).SetText(selectedConn.Host)
	getObject("connPortEntry").(*gtk.Entry).SetText(fmt.Sprintf("%v", selectedConn.Port))
	if len(selectedConn.User) != 0 && len(selectedConn.Password) != 0 {
		getObject("connUserEntry").(*gtk.Entry).SetText(selectedConn.User)
		getObject("connPwdEntry").(*gtk.Entry).SetText("***")
	}
}

func initConnsListBox() {
	connListBox := getConnListBox()
	connListBox.Connect("row-selected", onConnListBoxRowSelected)

	connInfos := ConnRepository.FindAll()
	for _, connInfo := range connInfos {
		label, err := gtk.LabelNew(connInfo.Name)
		util.CheckError(err)
		// set tooltip to hold connection name and to be used further
		// to get connection settings by name.
		// looks like go gtk implementation doesn't have separate method
		// to get label text and tooltip is the only way I've found to fetch
		// connection name when connection is selected. this is looks ugly
		label.SetTooltipText(connInfo.Name)

		connListBox.Add(label)
	}
	connListBox.SelectRow(connListBox.GetRowAtIndex(0))
	connListBox.ShowAll()
}

func getConnListBox() *gtk.ListBox {
	return getObject("connList").(*gtk.ListBox)
}

func onConnAddBtnClicked() {
	log.Print("Conn add btn clicked")
}

func onConnBtnClicked(connDialog *gtk.Dialog) func() {
	return func() {
		connDialog.Hide()
		ClearNodeTree()
		ShowTreeRootNodes()
	}
}

func onConnDialogCancelBtnClicked(connDialog *gtk.Dialog) func() {
	return func() {
		connDialog.Hide()
	}
}
