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
	connDlg                      = ConnDlg{}
)

type ConnDlg struct {
	dlg                 *gtk.Dialog
	connDialogCancelBtn *gtk.Button
	connNameEntry       *gtk.Entry
	connHostEntry       *gtk.Entry
	connPortEntry       *gtk.Entry
	connUserEntry       *gtk.Entry
	connPwdEntry        *gtk.Entry
	connIdEntry         *gtk.Entry
	connListBox         *gtk.ListBox
	connCopyBtn         *gtk.Button
	connDeleteBtn       *gtk.Button
}

func InitConnDialog(mainWindow *gtk.Window) {
	connDlg.dlg = GetObject("connDialog").(*gtk.Dialog)
	connDlg.dlg.SetTransientFor(mainWindow)
	connDlg.dlg.SetPosition(gtk.WIN_POS_CENTER)

	connDlg.connPortEntry = GetObject("connPortEntry").(*gtk.Entry)
	connDlg.connPortEntry.SetWidthChars(10)

	connDlg.connIdEntry = GetObject("connIdEntry").(*gtk.Entry)
	connDlg.connNameEntry = GetObject("connNameEntry").(*gtk.Entry)
	connDlg.connHostEntry = GetObject("connHostEntry").(*gtk.Entry)
	connDlg.connPortEntry = GetObject("connPortEntry").(*gtk.Entry)
	connDlg.connUserEntry = GetObject("connUserEntry").(*gtk.Entry)
	connDlg.connPwdEntry = GetObject("connPwdEntry").(*gtk.Entry)

	connDlg.connCopyBtn = GetObject("connCopyBtn").(*gtk.Button)
	connDlg.connCopyBtn.Connect("clicked", onConnCopyBtnClicked)
	connDlg.connDeleteBtn = GetObject("connDeleteBtn").(*gtk.Button)
	connDlg.connDeleteBtn.Connect("clicked", onConnDeleteBtnClicked)
	GetObject("connDialogCancelBtn").(*gtk.Button).Connect("clicked", onConnDialogCancelBtnClicked(connDlg.dlg))
	GetObject("connAddBtn").(*gtk.Button).Connect("clicked", onConnAddBtnClicked)
	GetObject("connSaveBtn").(*gtk.Button).Connect("clicked", onConnSaveBtnClicked)
	GetObject("connTestBtn").(*gtk.Button).Connect("clicked", onConnTestBtnClicked)
	GetObject("connBtn").(*gtk.Button).Connect("clicked", onConnBtnClicked)

	connInfos := ConnRepo.FindAll()
	if len(connInfos) == 0 {
		enableConnActions(false)
	}

	connDlg.connListBox = GetObject("connList").(*gtk.ListBox)
	initConnListBox()

	connDlg.dlg.ShowAll()
}

//todo cache it for session
func getSelectedConn() *core.ConnInfo {
	//todo leave it for test
	//return &core.ConnInfo{Name: "localhost", Host: "127.0.0.1", Port: 2181}
	//return &core.ConnInfo{Name: "scotia-nightly", Host: "172.0.30.173", Port: 32090, User: "zookeeper", Password: "z00k33p3r"}
	//return &core.ConnInfo{Name: "sandbox-pleeco", Host: "10.1.1.112", Port: 2181, User: "zookeeper", Password: "z00k33p3r"}
	//return &core.ConnInfo{Name: "scotia-history", Host: "172.0.30.173", Port: 32216, User: "zookeeper", Password: "z00k33p3r"}
	connName := getSelectedConnName(connDlg.connListBox)
	return ConnRepo.FindByName(connName)
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
	connDlg.connIdEntry.SetText(fmt.Sprintf("%v", connInfo.Id))
	connDlg.connNameEntry.SetText(connInfo.Name)
	connDlg.connHostEntry.SetText(connInfo.Host)

	connDlg.connPortEntry.SetText("")
	if connInfo.Port != 0 {
		connDlg.connPortEntry.SetText(fmt.Sprintf("%v", connInfo.Port))
	}

	if withMandatory {
		context, _ := connDlg.connNameEntry.GetStyleContext()
		context.AddClass(CssClassMandatory)
		context, _ = connDlg.connHostEntry.GetStyleContext()
		context.AddClass(CssClassMandatory)
		context, _ = connDlg.connPortEntry.GetStyleContext()
		context.AddClass(CssClassMandatory)
	}

	if len(connInfo.User) != 0 && len(connInfo.Password) != 0 {
		connDlg.connUserEntry.SetText(connInfo.User)
		connDlg.connPwdEntry.SetText(util.DecryptOrPanic(connInfo.Password))
	} else {
		connDlg.connUserEntry.SetText("")
		connDlg.connPwdEntry.SetText("")
	}
}

func initConnListBox() {
	connDlg.connListBox.Connect("row-selected", onConnSelected)
	//connListBox.Connect("button-release-event", onConnListBoxBtnRelease)
	connDlg.connListBox.Connect("button-press-event", onConnListBoxBtnPress)

	// clear conn list before load saved conns from repo
	// it's done to use delete, add or copy conn buttons
	// better to redraw list box then do manual delete, add or copy ops
	children := connDlg.connListBox.GetChildren()
	children.Foreach(func(item interface{}) {
		connDlg.connListBox.Remove(item.(gtk.IWidget))
	})
	connInfos := ConnRepo.FindAll()
	sort.Slice(connInfos, func(i, j int) bool {
		return connInfos[i].Name < connInfos[j].Name
	})
	for _, connInfo := range connInfos {
		addConnListBoxItem(connInfo)
	}
	connDlg.connListBox.ShowAll()
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

	connDlg.connListBox.Add(label)
	connDlg.connListBox.ShowAll()
}

func onConnSelected() {
	selectedConn := getSelectedConn()
	if selectedConn == nil {
		return
	}

	enableConnActions(true)
	drawConnInfo(selectedConn, false)
}

func enableConnActions(value bool) {
	connDlg.connCopyBtn.SetSensitive(value)
	connDlg.connDeleteBtn.SetSensitive(value)
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
	connDlg.connListBox.UnselectAll()
	enableConnActions(false)
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
	connName, _ := connDlg.connNameEntry.GetText()
	connHost, _ := connDlg.connHostEntry.GetText()
	connPort, _ := connDlg.connPortEntry.GetText()

	if len(connName) == 0 || len(connHost) == 0 || len(connPort) == 0 {
		dialog := createInfoDialog(
			connDlg.dlg,
			"Mandatory fields should not be empty",
		)
		dialog.Run()
		dialog.Hide()

		return nil
	}

	connUser, _ := connDlg.connUserEntry.GetText()
	connPwd, _ := connDlg.connPwdEntry.GetText()
	if len(connUser) != 0 && len(connPwd) == 0 {
		dialog := createInfoDialog(connDlg.dlg, "Password should be provided together with user")
		dialog.Run()
		dialog.Hide()

		context, _ := connDlg.connUserEntry.GetStyleContext()
		context.AddClass(CssClassMandatory)

		context, _ = connDlg.connPwdEntry.GetStyleContext()
		context.AddClass(CssClassMandatory)

		return nil
	}

	context, _ := connDlg.connNameEntry.GetStyleContext()
	context.RemoveClass(CssClassMandatory)
	context, _ = connDlg.connHostEntry.GetStyleContext()
	context.RemoveClass(CssClassMandatory)
	context, _ = connDlg.connUserEntry.GetStyleContext()
	context.RemoveClass(CssClassMandatory)
	context, _ = connDlg.connPwdEntry.GetStyleContext()
	context.RemoveClass(CssClassMandatory)

	connPortInt, err := strconv.Atoi(connPort)
	if err != nil || connPortInt < 0 {
		dialog := createInfoDialog(connDlg.dlg, "Connection port should be positive number")
		dialog.Run()
		dialog.Hide()
		return nil
	}

	context, _ = connDlg.connPortEntry.GetStyleContext()
	context.RemoveClass(CssClassMandatory)

	id := int64(0)
	idTxt, _ := connDlg.connIdEntry.GetText()
	if len(idTxt) != 0 {
		id, _ = strconv.ParseInt(idTxt, 10, 64)
	}

	encryptedPwd, err := util.Encrypt(connPwd)
	if err != nil {
		dialog := createInfoDialog(connDlg.dlg, err.Error())
		dialog.Run()
		dialog.Hide()
		return nil
	}

	connInfo := &core.ConnInfo{
		Id:       id,
		Name:     connName,
		Host:     connHost,
		Port:     connPortInt,
		User:     connUser,
		Password: encryptedPwd,
		Type:     core.ConnManual,
	}
	return connInfo
}

func onConnCopyBtnClicked() {
	connInfo := getSelectedConn()
	connInfoCopy := connInfo.Copy()
	connInfoCopy.Name += " - copy"

	connDlg.connListBox.UnselectAll()
	enableConnActions(false)
	drawConnInfo(connInfoCopy, false)
}

func onConnDeleteBtnClicked() {
	selectedConn := getSelectedConn()
	dialog := createConfirmDialog(connDlg.dlg, "Are you sure you want to delete "+selectedConn.Name+"?")
	resp := dialog.Run()
	if resp == gtk.RESPONSE_YES {
		ConnRepo.Delete(selectedConn)
		connDlg.connListBox.Remove(connDlg.connListBox.GetSelectedRow())
	}

	dialog.Hide()
	conns := ConnRepo.FindAll()
	if len(conns) == 0 {
		enableConnActions(false)
	}
}

func onConnBtnClicked() {
	zk.Reset()
	ClearNodeTree()

	connInfo := getSelectedConn()
	if connInfo == nil {
		return
	}
	zk.CachingRepo.SetConnInfo(connInfo)

	err := ShowTreeRootNodes()
	if err != nil {
		dialog := CreateErrorDialog(connDlg.dlg, "Unable to connect to "+connInfo.Name)
		dialog.Run()
		dialog.Hide()
		return
	}
	connDlg.dlg.Hide()
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
		dialog = createInfoDialog(connDlg.dlg, "Successfully connected to "+connInfo.Name)
	} else {
		cause := errors.Cause(err).(retry.Error)
		wrappedErrors := cause.WrappedErrors()
		errMsg := wrappedErrors[len(wrappedErrors)-1].Error()
		dialog = createWarnDialog(connDlg.dlg, errMsg)
	}
	dialog.Run()
	dialog.Hide()
}

func getConnForm() *core.ConnInfo {
	connName, _ := connDlg.connNameEntry.GetText()
	connHost, _ := connDlg.connHostEntry.GetText()
	connPort, _ := connDlg.connPortEntry.GetText()
	connUser, _ := connDlg.connUserEntry.GetText()
	connPwd, _ := connDlg.connPwdEntry.GetText()

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
