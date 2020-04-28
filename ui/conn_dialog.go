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
	"strconv"
)

const (
	CssClassMandatory = "mandatory"
	ConnNameDraft     = ""
)

var ConnRepo core.ConnRepository = &core.JsonConnRepository{}

var (
	lastBtnClickTime       = int64(0)
	lastSelectedConnName   string
	DoubleClickedBtnPeriod = int64(500 * 1e6)
)

func InitConnDialog(mainWindow *gtk.Window) *gtk.Dialog {
	getObject("connPortEntry").(*gtk.Entry).SetWidthChars(10)

	connDialog := getObject("connDialog").(*gtk.Dialog)
	connDialog.SetTransientFor(mainWindow)
	connDialog.SetPosition(gtk.WIN_POS_CENTER_ON_PARENT)

	getObject("connDialogCancelBtn").(*gtk.Button).Connect("clicked", onConnDialogCancelBtnClicked(connDialog))
	getObject("connAddBtn").(*gtk.Button).Connect("clicked", onConnAddBtnClicked)
	connSaveBtn := getObject("connSaveBtn").(*gtk.Button)
	connSaveBtn.Connect("clicked", onConnSaveBtnClicked)

	connCopyBtn := getObject("connCopyBtn").(*gtk.Button)
	connCopyBtn.Connect("clicked", onConnCopyBtnClicked)

	connDeleteBtn := getObject("connDeleteBtn").(*gtk.Button)
	connDeleteBtn.Connect("clicked", onConnDeleteBtnClicked)

	connInfos := ConnRepo.FindAll()
	if len(connInfos) == 0 {
		connSaveBtn.SetSensitive(false)
		connCopyBtn.SetSensitive(false)
		connDeleteBtn.SetSensitive(false)
	}

	getObject("connBtn").(*gtk.Button).Connect("clicked", onConnBtnClicked(connDialog))

	initConnListBox()
	initCssProvider()

	connDialog.ShowAll()
	return connDialog
}

func initCssProvider() {
	providerNew, err := gtk.CssProviderNew()
	if err != nil {
		log.WithError(err).Errorf("Failed to create CSS provider")
	}

	err = providerNew.LoadFromPath(core.CssStyleFilePath)
	if err != nil {
		log.WithError(err).Errorf("Failed to load CSS styles")
	}

	defaultScreen, err := gdk.ScreenGetDefault()
	if err != nil {
		log.WithError(err).Errorf("Failed to get default screen")
	}

	gtk.AddProviderForScreen(defaultScreen, providerNew, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
}

//TODO cache it for session
func GetSelectedConn() *core.ConnInfo {
	//TODO leave it for test
	//return &core.ConnInfo{
	//	Name:     "kelp-nightly",
	//	Host:     "172.0.30.173",
	//	Port:     32050,
	//	User:     "zookeeper",
	//	Password: "z00k33p3r",
	//}
	connList := getConnListBox()
	connName := getSelectedConnName(connList)
	if len(connName) == 0 {
		return nil
	}

	return ConnRepo.Find(connName)
}

func getSelectedConnName(connList *gtk.ListBox) string {
	selectedRow := connList.GetSelectedRow()
	if selectedRow == nil {
		return ""
	}

	child, err := selectedRow.GetChild()
	if err != nil {
		log.WithError(err).Errorf("Failed to get list box selected child")
	}

	connName, _ := child.GetTooltipText()
	return connName
}

func drawConnInfo(connInfo *core.ConnInfo, withMandatory bool) {
	connNameEntry := getObject("connNameEntry").(*gtk.Entry)
	connNameEntry.SetText(connInfo.Name)

	connHostEntry := getObject("connHostEntry").(*gtk.Entry)
	connHostEntry.SetText(connInfo.Host)

	connPortEntry := getObject("connPortEntry").(*gtk.Entry)
	if connInfo.Port != 0 {
		connPortEntry.SetText(fmt.Sprintf("%v", connInfo.Port))
	}

	if withMandatory {
		context, _ := connNameEntry.GetStyleContext()
		context.AddClass(CssClassMandatory)
		context, _ = connHostEntry.GetStyleContext()
		context.AddClass(CssClassMandatory)
		context, _ = connPortEntry.GetStyleContext()
		context.AddClass(CssClassMandatory)
	}

	if len(connInfo.User) != 0 && len(connInfo.Password) != 0 {
		getObject("connUserEntry").(*gtk.Entry).SetText(connInfo.User)
		getObject("connPwdEntry").(*gtk.Entry).SetText("******")
	} else {
		getObject("connUserEntry").(*gtk.Entry).SetText("")
		getObject("connPwdEntry").(*gtk.Entry).SetText("")
	}
}

func initConnListBox() {
	connListBox := getConnListBox()
	connListBox.Connect("row-selected", onConnListBoxRowSelected)
	//connListBox.Connect("button-release-event", onConnListBoxBtnPress)

	// clear conn list before load saved conns from repo
	// it's done to use delete, add or copy conn buttons
	// better to redraw list box then do manual delete, add or copy ops
	children := connListBox.GetChildren()
	children.Foreach(func(item interface{}) {
		connListBox.Remove(item.(gtk.IWidget))
	})
	for _, connInfo := range ConnRepo.FindAll() {
		addConnListBoxItem(connInfo)
	}
	connListBox.ShowAll()
}

func addConnListBoxItem(conn *core.ConnInfo) {
	label, err := gtk.LabelNew(conn.Name)
	util.CheckError(err)
	// set tooltip to hold connection name and to be used further
	// to get connection settings by name.
	// looks like go gtk implementation doesn't have separate method
	// to get label text and tooltip is the only way I've found to fetch
	// connection name when connection is selected. this is looks ugly
	label.SetTooltipText(conn.Name)
	label.SetHAlign(gtk.ALIGN_START)

	connListBox := getConnListBox()
	connListBox.Add(label)
	connListBox.ShowAll()
}

func onConnListBoxRowSelected() {
	selectedConn := GetSelectedConn()
	if selectedConn == nil {
		return
	}

	drawConnInfo(selectedConn, false)
}

// TODO Double click handle has a bug. To reproduce:
// - create new connection
// - remove it
// - try to double click to another(not currently selected)
// - error occurs
//func onConnListBoxBtnPress(listBox *gtk.ListBox, e *gdk.Event) {
//	event := &gdk.EventKey{Event: e}
//	mouseButton := (*C.GdkEventButton)(unsafe.Pointer(event.Event.GdkEvent)).button
//	row := listBox.GetSelectedRow()
//
//	// row is nil when conn list is empty
//	if row != nil && mouseButton == 1 {
//		child, _ := row.GetChild()
//		connName, _ := child.GetTooltipText()
//		// didn't find a way how to register double button click
//		// looks like GO GTK adapter doesn't support such event
//		// implement my own bicycle
//		if time.Now().UnixNano()-lastBtnClickTime < DoubleClickedBtnPeriod && lastSelectedConnName == connName {
//			log.Infof("Selected conn: %s", connName)
//			onConnBtnClicked(getConnDialog())()
//		}
//
//		lastBtnClickTime = time.Now().UnixNano()
//		lastSelectedConnName = connName
//	}
//}

func onConnAddBtnClicked() {
	getConnListBox().UnselectAll()
	drawConnInfo(&core.ConnInfo{Name: ConnNameDraft}, true)
}

func onConnSaveBtnClicked() {
	connNameEntry := getObject("connNameEntry").(*gtk.Entry)
	connName, _ := connNameEntry.GetText()
	if len(connName) > 0 && isConnExists(connName) {
		dialog := createInfoDialog(
			getConnDialog(),
			"Connection name '"+connName+"' is not unique. Please choose another name",
		)
		dialog.Run()
		dialog.Hide()

		return
	}
	connHostEntry := getObject("connHostEntry").(*gtk.Entry)
	connHost, _ := connHostEntry.GetText()

	connPortEntry := getObject("connPortEntry").(*gtk.Entry)
	connPort, _ := connPortEntry.GetText()

	if len(connName) == 0 || len(connHost) == 0 || len(connPort) == 0 {
		dialog := createInfoDialog(
			getConnDialog(),
			"Mandatory fields should not be empty",
		)
		dialog.Run()
		dialog.Hide()

		return
	}

	connUserEntry := getObject("connUserEntry").(*gtk.Entry)
	connUser, _ := connUserEntry.GetText()

	connPwdEntry := getObject("connPwdEntry").(*gtk.Entry)
	connPwd, _ := connPwdEntry.GetText()
	if len(connUser) != 0 && len(connPwd) == 0 {
		dialog := createInfoDialog(
			getConnDialog(),
			"Password should be provided together with user",
		)
		dialog.Run()
		dialog.Hide()

		context, _ := connUserEntry.GetStyleContext()
		context.AddClass(CssClassMandatory)

		context, _ = connPwdEntry.GetStyleContext()
		context.AddClass(CssClassMandatory)

		return
	}

	context, _ := connNameEntry.GetStyleContext()
	context.RemoveClass(CssClassMandatory)
	context, _ = connHostEntry.GetStyleContext()
	context.RemoveClass(CssClassMandatory)
	context, _ = connUserEntry.GetStyleContext()
	context.RemoveClass(CssClassMandatory)
	context, _ = connPwdEntry.GetStyleContext()
	context.RemoveClass(CssClassMandatory)

	connPortInt, err := strconv.Atoi(connPort)
	if err != nil || connPortInt < 0 {
		dialog := createInfoDialog(
			getConnDialog(),
			"Connection port should be positive number",
		)
		dialog.Run()
		dialog.Hide()
		return
	}

	context, _ = connPortEntry.GetStyleContext()
	context.RemoveClass(CssClassMandatory)

	connInfo := &core.ConnInfo{
		Name:     connName,
		Host:     connHost,
		Port:     connPortInt,
		User:     connUser,
		Password: connPwd,
	}
	ConnRepo.Upsert(connInfo)
	initConnListBox()
}

func isConnExists(connName string) bool {
	for _, connInfo := range ConnRepo.FindAll() {
		if connName == connInfo.Name {
			return true
		}
	}

	return false
}

func onConnCopyBtnClicked() {
	log.Print("Conn add btn clicked")
}

func onConnDeleteBtnClicked() {
	selectedConn := GetSelectedConn()
	dialog := createConfirmDialog(getConnDialog(), "Are you sure you want to delete "+selectedConn.Name+"?")
	resp := dialog.Run()
	if resp == gtk.RESPONSE_YES {
		ConnRepo.Delete(selectedConn)
		connListBox := getConnListBox()
		connListBox.Remove(connListBox.GetSelectedRow())
	}

	dialog.Hide()
}

func onConnBtnClicked(connDialog *gtk.Dialog) func() {
	return func() {
		ClearNodeTree()
		err := ShowTreeRootNodes()
		if err != nil {
			connInfo := GetSelectedConn()
			dialog := createErrorDialog(getConnDialog(), "Unable to connect to "+connInfo.Name)
			dialog.Run()
			dialog.Hide()
			return
		}
		connDialog.Hide()
	}
}

func onConnDialogCancelBtnClicked(connDialog *gtk.Dialog) func() {
	return func() {
		connDialog.Hide()
	}
}

func getConnDialog() *gtk.Dialog {
	return getObject("connDialog").(*gtk.Dialog)
}

func getConnListBox() *gtk.ListBox {
	return getObject("connList").(*gtk.ListBox)
}
