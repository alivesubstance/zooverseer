package ui

// #cgo pkg-config: gdk-3.0 glib-2.0 gobject-2.0
// #include <gdk/gdk.h>
// #include "/home/mirian/code/go/src/github.com/gotk3/gotk3/gdk/gdk.go.h"
import "C"
import (
	"fmt"
	"github.com/alivesubstance/zooverseer/core"
	"github.com/alivesubstance/zooverseer/util"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	log "github.com/sirupsen/logrus"
	"time"
	"unsafe"
)

var ConnRepo core.ConnRepository = &core.JsonConnRepository{}

var (
	lastBtnClickTime       = int64(0)
	lastSelectedConnName   string
	DoubleClickedBtnPeriod = int64(500 * 1e6)
)

func InitConnDialog(mainWindow *gtk.Window) *gtk.Dialog {
	connPortEntry := getObject("connPortEntry").(*gtk.Entry)
	connPortEntry.SetWidthChars(10)

	connDialog := getObject("connDialog").(*gtk.Dialog)
	connDialog.SetTransientFor(mainWindow)
	connDialog.SetPosition(gtk.WIN_POS_CENTER)

	connDialogCancelBtn := getObject("connDialogCancelBtn").(*gtk.Button)
	connDialogCancelBtn.Connect("clicked", onConnDialogCancelBtnClicked(connDialog))

	connAddBtn := getObject("connAddBtn").(*gtk.Button)
	connAddBtn.Connect("clicked", onConnAddBtnClicked)

	connBtn := getObject("connBtn").(*gtk.Button)
	connBtn.Connect("clicked", onConnBtnClicked(connDialog))

	initConnListBox()

	connDialog.ShowAll()

	return connDialog
}

//TODO cache it for session
func GetSelectedConn() *core.JsonConnInfo {
	//TODO leave it for test
	//return &core.JsonConnInfo{
	//	Host: "10.1.1.112",
	//	Port: 2181,
	//	User:     "zookeeper",
	//	Password: "z00k33p3r",
	//}
	connList := getObject("connList").(*gtk.ListBox)
	connName := getSelectedConnName(connList)
	connInfo := ConnRepo.Find(connName)
	if connInfo == nil {
		log.Panicf("'%v' connection setting not found. Should never happened", connName)
	}

	return connInfo
}

func getSelectedConnName(connList *gtk.ListBox) string {
	child, err := connList.GetSelectedRow().GetChild()
	if err != nil {
		log.WithError(err).Errorf("Failed to get list box selected child")
	}

	connName, _ := child.GetTooltipText()
	return connName
}

func onConnListBoxRowSelected() {
	selectedConn := GetSelectedConn()
	getObject("connNameEntry").(*gtk.Entry).SetText(selectedConn.Name)
	getObject("connHostEntry").(*gtk.Entry).SetText(selectedConn.Host)
	getObject("connPortEntry").(*gtk.Entry).SetText(fmt.Sprintf("%v", selectedConn.Port))
	if len(selectedConn.User) != 0 && len(selectedConn.Password) != 0 {
		getObject("connUserEntry").(*gtk.Entry).SetText(selectedConn.User)
		getObject("connPwdEntry").(*gtk.Entry).SetText("******")
	} else {
		getObject("connUserEntry").(*gtk.Entry).SetText("")
		getObject("connPwdEntry").(*gtk.Entry).SetText("")
	}
}

func initConnListBox() {
	connListBox := getConnListBox()
	connListBox.Connect("row-selected", onConnListBoxRowSelected)
	connListBox.Connect("button-press-event", onConnListBoxDoubleClick)

	connInfos := ConnRepo.FindAll()
	for _, connInfo := range connInfos {
		label, err := gtk.LabelNew(connInfo.Name)
		util.CheckError(err)
		// set tooltip to hold connection name and to be used further
		// to get connection settings by name.
		// looks like go gtk implementation doesn't have separate method
		// to get label text and tooltip is the only way I've found to fetch
		// connection name when connection is selected. this is looks ugly
		label.SetTooltipText(connInfo.Name)
		label.SetHAlign(gtk.ALIGN_START)

		connListBox.Add(label)
	}
	connListBox.SelectRow(connListBox.GetRowAtIndex(0))
	connListBox.ShowAll()
}

func onConnListBoxDoubleClick(listBox *gtk.ListBox, e *gdk.Event) {
	event := &gdk.EventKey{Event: e}
	mouseButton := (*C.GdkEventButton)(unsafe.Pointer(event.Event.GdkEvent)).button

	if mouseButton == 1 {
		row := listBox.GetSelectedRow()
		child, _ := row.GetChild()
		connName, _ := child.GetTooltipText()
		if time.Now().UnixNano()-lastBtnClickTime < DoubleClickedBtnPeriod && lastSelectedConnName == connName {
			log.Infof("Selected conn: %s", connName)
			onConnBtnClicked(getObject("connDialog").(*gtk.Dialog))()
		}

		lastBtnClickTime = time.Now().UnixNano()
		lastSelectedConnName = connName
	}
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
