package main

// #cgo pkg-config: gdk-3.0 glib-2.0 gobject-2.0
// #include <gdk/gdk.h>
// #include "/home/mirian/code/go/src/github.com/gotk3/gotk3/gdk/gdk.go.h"
import "C"
import (
	"fmt"
	"github.com/alivesubstance/zooverseer/core"
	"github.com/alivesubstance/zooverseer/util"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"log"
	"os"
	"unsafe"
)

// there are rumors that global variable is evil. why?
var (
	Builder        *gtk.Builder
	ConnRepository core.ConnRepository = core.JsonConnInfo{}
)

func main() {
	log.Print("Starting zooverseer")

	app, err := gtk.ApplicationNew(core.AppId, glib.APPLICATION_FLAGS_NONE)
	util.CheckError(err)

	app.Connect("activate", onAppActivate(app))

	os.Exit(app.Run(os.Args))
}

func onAppActivate(app *gtk.Application) func() {
	return func() {
		log.Print("Reading glade file")
		builder, err := gtk.BuilderNewFromFile(core.GladeFilePath)
		util.CheckError(err)

		Builder = builder

		mainWindow := getObject("mainWindow").(*gtk.Window)
		mainWindow.Connect("button-release-event", on_button_press_event)
		//connDialog := initConnDialog(mainWindow)
		//connDialog.ShowAll()

		mainWindow.ShowAll()
		app.AddWindow(mainWindow)

	}
}

func on_button_press_event(b *gtk.Window, e *gdk.Event) {
	mouseButton := getMouseButton(e)
	if mouseButton == 3 {
		// if mouse2 button pressed
		menu := getObject("popupMenu").(*gtk.Menu)

		menu.ShowAll()
		menu.PopupAtPointer(e)
	}
}

func getMouseButton(e *gdk.Event) uint {
	event := &gdk.EventKey{e}
	mouseButton := (*C.GdkEventButton)(unsafe.Pointer(event.Event.GdkEvent)).button
	return uint(mouseButton)
}

func initConnDialog(mainWindow *gtk.Window) *gtk.Dialog {
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

	return connDialog
}

func initConnsListBox() {
	connListBox := getConnListBox()
	connListBox.Connect("row-selected", onConnListBoxRowSelected())

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

func onConnListBoxRowSelected() func(listBox *gtk.ListBox, row *gtk.ListBoxRow) {
	return func(listBox *gtk.ListBox, row *gtk.ListBoxRow) {
		selectedConn := getSelectedConn(row)
		getObject("connNameEntry").(*gtk.Entry).SetText(selectedConn.Name)
		getObject("connHostEntry").(*gtk.Entry).SetText(selectedConn.Host)
		getObject("connPortEntry").(*gtk.Entry).SetText(fmt.Sprintf("%v", selectedConn.Port))
		getObject("connUserEntry").(*gtk.Entry).SetText(selectedConn.User)
		getObject("connPwdEntry").(*gtk.Entry).SetText("***")
	}
}

func getSelectedConn(row *gtk.ListBoxRow) *core.JsonConnInfo {
	child, err := row.GetChild()
	util.CheckError(err)

	connName, _ := child.GetTooltipText()
	connInfo, ok := ConnRepository.Find(connName)
	if !ok {
		log.Panicf("'%s' connection setting not found. Should never happened", connName)
	}

	return connInfo
}

func onConnAddBtnClicked() {
	log.Print("Conn add btn clicked")
}

func onConnBtnClicked(connDialog *gtk.Dialog) func() {
	return func() {
		connDialog.Hide()

	}
}

func onConnDialogCancelBtnClicked(connDialog *gtk.Dialog) func() {
	return func() {
		connDialog.Hide()
	}
}

func getConnListBox() *gtk.ListBox {
	return getObject("connList").(*gtk.ListBox)
}

func getObject(objectName string) glib.IObject {
	object, err := Builder.GetObject(objectName)
	util.CheckError(err)

	return object
}
