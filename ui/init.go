package ui

import (
	"github.com/alivesubstance/zooverseer/core"
	"github.com/alivesubstance/zooverseer/util"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	log "github.com/sirupsen/logrus"
)

var (
	Builder    *gtk.Builder
	mainWindow *MainWindow
)

type Signals map[string]interface{}

func OnAppActivate(app *gtk.Application) func() {
	return func() {
		log.Infof("Reading glade file %v", core.Config.GladeFilePath)
		builder, err := gtk.BuilderNewFromFile(core.Config.GladeFilePath)
		util.CheckError(err)

		Builder = builder

		connectSignals(builder)

		mainWindow = NewMainWindow()
		InitConnDialog(mainWindow)
		initAppIcon(mainWindow)

		app.AddWindow(mainWindow.gtkWindow)
	}
}

func connectSignals(builder *gtk.Builder) {
	menuAboutSignals := Signals{
		"onMenuAboutBtnCloseClicked": onMenuAboutBtnCloseClicked,
		"onMenuAboutActivate":        onMenuAboutActivate,
	}

	connDlgSignals := Signals{
		"onConnCopyBtnClicked":   onConnCopyBtnClicked,
		"onConnDeleteBtnClicked": onConnDeleteBtnClicked,
		"onConnTestBtnClicked":   onConnTestBtnClicked,
		"onConnBtnClicked":       onConnBtnClicked,
		"onConnAddBtnClicked":    onConnAddBtnClicked,
		"onConnSaveBtnClicked":   onConnSaveBtnClicked,
	}

	signals := mergeSignals(Signals{}, menuAboutSignals)
	signals = mergeSignals(signals, connDlgSignals)

	builder.ConnectSignals(signals)
}

func CreateErrorDialog(parent gtk.IWindow, text string) *gtk.MessageDialog {
	return gtk.MessageDialogNew(parent, gtk.DIALOG_MODAL, gtk.MESSAGE_ERROR, gtk.BUTTONS_OK, text)
}

func createConfirmDialog(parent gtk.IWindow, text string) *gtk.MessageDialog {
	dlg := gtk.MessageDialogNew(parent, gtk.DIALOG_MODAL, gtk.MESSAGE_QUESTION, gtk.BUTTONS_YES_NO, text)
	//todo doesn't work
	//dlg.SetPosition(gtk.WIN_POS_CENTER_ON_PARENT)
	return dlg
}

func createInfoDialog(parent gtk.IWindow, text string) *gtk.MessageDialog {
	return gtk.MessageDialogNew(parent, gtk.DIALOG_MODAL, gtk.MESSAGE_INFO, gtk.BUTTONS_OK, text)
}

func createWarnDialog(parent gtk.IWindow, text string) *gtk.MessageDialog {
	return gtk.MessageDialogNew(parent, gtk.DIALOG_MODAL, gtk.MESSAGE_WARNING, gtk.BUTTONS_OK, text)
}

func GetObject(objectName string) glib.IObject {
	object, err := Builder.GetObject(objectName)
	util.CheckError(err)

	return object
}

func mergeSignals(dst Signals, src Signals) Signals {
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func initAppIcon(window *MainWindow) {
	err := window.gtkWindow.SetIconFromFile(core.Config.LogoFilePath)
	if err != nil {
		log.WithError(err).Panicf("Failed to set icon from %v", core.Config.LogoFilePath)
	}
}
