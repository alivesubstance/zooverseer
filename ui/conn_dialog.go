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
	"github.com/avast/retry-go"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"sort"
	"strconv"
)

const (
	CssClassMandatory = "mandatory"
	ConnNameDraft     = ""
)

var (
	ConnRepo core.ConnRepository = &core.JsonConnRepository{}
)

func InitConnDialog(mainWindow *gtk.Window) *gtk.Dialog {
	GetObject("connPortEntry").(*gtk.Entry).SetWidthChars(10)

	connDialog := getConnDialog()
	connDialog.SetTransientFor(mainWindow)
	connDialog.SetPosition(gtk.WIN_POS_CENTER)

	GetObject("connDialogCancelBtn").(*gtk.Button).Connect("clicked", onConnDialogCancelBtnClicked(connDialog))
	GetObject("connAddBtn").(*gtk.Button).Connect("clicked", onConnAddBtnClicked)
	connSaveBtn := GetObject("connSaveBtn").(*gtk.Button)
	connSaveBtn.Connect("clicked", onConnSaveBtnClicked)

	connCopyBtn := GetObject("connCopyBtn").(*gtk.Button)
	connCopyBtn.Connect("clicked", onConnCopyBtnClicked)

	connDeleteBtn := GetObject("connDeleteBtn").(*gtk.Button)
	connDeleteBtn.Connect("clicked", onConnDeleteBtnClicked)

	connTestBtn := GetObject("connTestBtn").(*gtk.Button)
	connTestBtn.Connect("clicked", onConnTestBtnClicked)

	connInfos := ConnRepo.FindAll()
	if len(connInfos) == 0 {
		connSaveBtn.SetSensitive(false)
		connCopyBtn.SetSensitive(false)
		connDeleteBtn.SetSensitive(false)
	}

	GetObject("connBtn").(*gtk.Button).Connect("clicked", onConnBtnClicked)

	initConnListBox()

	connDialog.ShowAll()
	return connDialog
}

//todo cache it for session
func getSelectedConn() *core.ConnInfo {
	//todo leave it for test
	//return &core.ConnInfo{Name: "localhost", Host: "127.0.0.1", Port: 2181}
	//return &core.ConnInfo{Name: "scotia-nightly", Host: "172.0.30.173", Port: 32090, User: "zookeeper", Password: "z00k33p3r"}
	//return &core.ConnInfo{Name: "sandbox-pleeco", Host: "10.1.1.112", Port: 2181, User: "zookeeper", Password: "z00k33p3r"}
	//return &core.ConnInfo{Name: "scotia-history", Host: "172.0.30.173", Port: 32216, User: "zookeeper", Password: "z00k33p3r"}
	connList := getConnListBox()
	connName := getSelectedConnName(connList)
	return ConnRepo.Find(connName)
}

func getSelectedConnName(connList *gtk.ListBox) string {
	selectedRow := connList.GetSelectedRow()
	if selectedRow == nil {
		return ""
	}

	child, err := selectedRow.GetChild()
	if err != nil {
		log.WithError(err).Panicf("Failed to get list box selected child")
	}

	connName, _ := child.GetTooltipText()
	return connName
}

func drawConnInfo(connInfo *core.ConnInfo, withMandatory bool) {
	connNameEntry := GetObject("connNameEntry").(*gtk.Entry)
	connNameEntry.SetText(connInfo.Name)

	connHostEntry := GetObject("connHostEntry").(*gtk.Entry)
	connHostEntry.SetText(connInfo.Host)

	connPortEntry := GetObject("connPortEntry").(*gtk.Entry)
	if connInfo.Port != 0 {
		connPortEntry.SetText(fmt.Sprintf("%v", connInfo.Port))
	} else {
		connPortEntry.SetText("")
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
		GetObject("connUserEntry").(*gtk.Entry).SetText(connInfo.User)
		GetObject("connPwdEntry").(*gtk.Entry).SetText(connInfo.Password)
	} else {
		GetObject("connUserEntry").(*gtk.Entry).SetText("")
		GetObject("connPwdEntry").(*gtk.Entry).SetText("")
	}
}

func initConnListBox() {
	connListBox := getConnListBox()
	connListBox.Connect("row-selected", onConnSelected)
	//connListBox.Connect("button-release-event", onConnListBoxBtnRelease)
	connListBox.Connect("button-press-event", onConnListBoxBtnPress)

	// clear conn list before load saved conns from repo
	// it's done to use delete, add or copy conn buttons
	// better to redraw list box then do manual delete, add or copy ops
	children := connListBox.GetChildren()
	children.Foreach(func(item interface{}) {
		connListBox.Remove(item.(gtk.IWidget))
	})
	connInfos := ConnRepo.FindAll()
	sort.Slice(connInfos, func(i, j int) bool {
		return connInfos[i].Name < connInfos[j].Name
	})
	for _, connInfo := range connInfos {
		addConnListBoxItem(connInfo)
	}
	connListBox.ShowAll()
}

func addConnListBoxItem(conn *core.ConnInfo) {
	label, err := gtk.LabelNew(conn.Name)
	util.CheckError(err)
	// set tooltip to hold connection name for further using
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

func onConnSelected() {
	selectedConn := getSelectedConn()
	if selectedConn == nil {
		return
	}

	setConnListBoxBtnsSensitivity(true)
	drawConnInfo(selectedConn, false)
}

func setConnListBoxBtnsSensitivity(value bool) {
	GetObject("connCopyBtn").(*gtk.Button).SetSensitive(value)
	GetObject("connDeleteBtn").(*gtk.Button).SetSensitive(value)
}

func onConnListBoxBtnPress(listBox *gtk.ListBox, e *gdk.Event) {
	btnEvent := gdk.EventButtonNewFromEvent(e)
	row := listBox.GetSelectedRow()

	// row is nil when conn list is empty
	if row != nil && btnEvent.Button() == 1 && btnEvent.Type() == gdk.EVENT_DOUBLE_BUTTON_PRESS {
		child, _ := row.GetChild()
		connName, _ := child.GetTooltipText()

		log.Infof("Selected conn: %s", connName)
		onConnBtnClicked()
	}
}

func onConnAddBtnClicked() {
	getConnListBox().UnselectAll()
	setConnListBoxBtnsSensitivity(false)
	drawConnInfo(&core.ConnInfo{Name: ConnNameDraft}, true)
}

func onConnSaveBtnClicked() {
	connInfo := validateAndGetConn()
	if connInfo == nil {
		return
	}
	ConnRepo.Upsert(connInfo)
	initConnListBox()
}

func validateAndGetConn() *core.ConnInfo {
	connNameEntry := GetObject("connNameEntry").(*gtk.Entry)
	connName, _ := connNameEntry.GetText()

	connHostEntry := GetObject("connHostEntry").(*gtk.Entry)
	connHost, _ := connHostEntry.GetText()

	connPortEntry := GetObject("connPortEntry").(*gtk.Entry)
	connPort, _ := connPortEntry.GetText()

	if len(connName) == 0 || len(connHost) == 0 || len(connPort) == 0 {
		dialog := createInfoDialog(
			getConnDialog(),
			"Mandatory fields should not be empty",
		)
		dialog.Run()
		dialog.Hide()

		return nil
	}

	connUserEntry := GetObject("connUserEntry").(*gtk.Entry)
	connUser, _ := connUserEntry.GetText()

	connPwdEntry := GetObject("connPwdEntry").(*gtk.Entry)
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

		return nil
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
		return nil
	}

	context, _ = connPortEntry.GetStyleContext()
	context.RemoveClass(CssClassMandatory)

	connInfo := &core.ConnInfo{
		Name:     connName,
		Host:     connHost,
		Port:     connPortInt,
		User:     connUser,
		Password: connPwd,
		Type:     core.ConnManual,
	}
	return connInfo
}

func onConnCopyBtnClicked() {
	connInfo := getSelectedConn()
	connInfoCopy := connInfo.Copy()
	connInfoCopy.Name += " - copy"

	getConnListBox().UnselectAll()
	setConnListBoxBtnsSensitivity(false)
	drawConnInfo(connInfoCopy, false)
}

func onConnDeleteBtnClicked() {
	selectedConn := getSelectedConn()
	dialog := createConfirmDialog(getConnDialog(), "Are you sure you want to delete "+selectedConn.Name+"?")
	resp := dialog.Run()
	if resp == gtk.RESPONSE_YES {
		ConnRepo.Delete(selectedConn)
		connListBox := getConnListBox()
		connListBox.Remove(connListBox.GetSelectedRow())
	}

	dialog.Hide()
}

func onConnBtnClicked() {
	ClearNodeTree()

	connInfo := getSelectedConn()
	if connInfo == nil {
		return
	}
	zk.CachingRepo.SetConnInfo(connInfo)

	err := ShowTreeRootNodes()
	if err != nil {
		dialog := CreateErrorDialog(getConnDialog(), "Unable to connect to "+connInfo.Name)
		dialog.Run()
		dialog.Hide()
		return
	}
	getConnDialog().Hide()
}

func onConnDialogCancelBtnClicked(connDialog *gtk.Dialog) func() {
	return func() {
		connDialog.Hide()
	}
}

func onConnTestBtnClicked() {
	connInfo := getConnForm()
	if connInfo.Host == "" || connInfo.Port == 0 {
		return
	}

	repo := zk.Repository{}
	repo.SetConnInfo(connInfo)
	_, err := repo.GetMeta(core.NodeRootName)

	var dialog *gtk.MessageDialog
	if err == nil {
		dialog = createInfoDialog(getConnDialog(), "Successfully connected to "+connInfo.Name)
	} else {
		cause := errors.Cause(err).(retry.Error)
		wrappedErrors := cause.WrappedErrors()
		errMsg := wrappedErrors[len(wrappedErrors)-1].Error()
		dialog = createWarnDialog(getConnDialog(), errMsg)
	}
	dialog.Run()
	dialog.Hide()
}

func getConnForm() *core.ConnInfo {
	connName, _ := GetObject("connNameEntry").(*gtk.Entry).GetText()
	connHost, _ := GetObject("connHostEntry").(*gtk.Entry).GetText()
	connPort, _ := GetObject("connPortEntry").(*gtk.Entry).GetText()
	connUser, _ := GetObject("connUserEntry").(*gtk.Entry).GetText()
	connPwd, _ := GetObject("connPwdEntry").(*gtk.Entry).GetText()

	connPortInt, _ := strconv.Atoi(connPort)
	connInfo := &core.ConnInfo{
		Name:     connName,
		Host:     connHost,
		Port:     connPortInt,
		User:     connUser,
		Password: connPwd,
	}
	return connInfo
}

func getConnDialog() *gtk.Dialog {
	return GetObject("connDialog").(*gtk.Dialog)
}

func getConnListBox() *gtk.ListBox {
	return GetObject("connList").(*gtk.ListBox)
}
