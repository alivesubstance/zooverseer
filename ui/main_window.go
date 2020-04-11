package ui

// #cgo pkg-config: gdk-3.0 glib-2.0 gobject-2.0
// #include <gdk/gdk.h>
// #include "/home/mirian/code/go/src/github.com/gotk3/gotk3/gdk/gdk.go.h"
import "C"
import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"unsafe"
)

func InitMainWindow(mainWindow *gtk.Window) {
	InitNodeTree()

	mainWindow.SetTitle("Zooverseer")
	mainWindow.ShowAll()
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
