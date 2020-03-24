package main

import (
	"github.com/gotk3/gotk3/glib"
	"log"
	"os"

	"github.com/gotk3/gotk3/gtk"
)

const appId = "com.github.alivesubstance.zooverseer"

var signals = map[string]interface{}{
	"onConnCancelBtnClicked": onConnCancelBtnClicked,
	"openConnDialog":         openConnDialog,
}

func main() {
	log.Print("Starting zooverseer")

	app, err := gtk.ApplicationNew(appId, glib.APPLICATION_FLAGS_NONE)
	checkError(err)

	app.Connect("activate", func() {
		// Создаём билдер
		builder, err := gtk.BuilderNew()
		checkError(err)

		// Загружаем в билдер окно из файла Glade
		// todo change to relative path
		err = builder.AddFromFile("/home/mirian/code/go/src/github.com/alivesubstance/zooverseer/assets/main.glade")
		checkError(err)

		builder.ConnectSignals(signals)

		connectDialog := initConnDialog(builder)
		connectDialog.ShowAll()

		// Преобразуем из объекта именно окно типа gtk.Window
		// и соединяем с сигналом "destroy" чтобы можно было закрыть
		// приложение при закрытии окна
		mainWindow := getObject(builder, "mainWindow").(*gtk.Window)
		//mainWindow.Connect("destroy", func() {
		//	gtk.MainQuit()
		//})

		// Отображаем все виджеты в окне
		mainWindow.ShowAll()
		app.AddWindow(mainWindow)
	})

	os.Exit(app.Run(os.Args))
}

func initConnDialog(builder *gtk.Builder) *gtk.Dialog {
	portEntry := getObject(builder, "connPortEntry").(*gtk.Entry)
	portEntry.SetWidthChars(10)

	connectDialog := getObject(builder, "connDialog").(*gtk.Dialog)

	return connectDialog
}

func onConnCancelBtnClicked() {
	log.Print("onConnCancelBtnClicked")
}

func openConnDialog() {
	log.Print("Open conn dialog")
}
