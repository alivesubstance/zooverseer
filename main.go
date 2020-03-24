package main

import (
	"log"

	"github.com/gotk3/gotk3/gtk"
)

func main() {
	// Инициализируем GTK.
	gtk.Init(nil)

	// Создаём билдер
	builder, err := gtk.BuilderNew()
	if err != nil {
		log.Fatal("Ошибка:", err)
	}

	// Загружаем в билдер окно из файла Glade
	err = builder.AddFromFile("/home/mirian/code/go/src/github.com/alivesubstance/zooverseer/assets/main.glade")
	if err != nil {
		log.Fatal("Ошибка:", err)
	}

	dlg := initConnDialog(builder)
	dlg.showAll()

	portEntry, err := builder.GetObject("portEntry")
	if err != nil {
		log.Fatal("Failed to get port entry")
	}
	port := portEntry.(*gtk.Entry)
	port.SetWidthChars(10)

	dlg.ShowAll()

	// Получаем объект главного окна по ID
	//obj, err := builder.GetObject("mainWindow")
	//if err != nil {
	//	log.Fatal("Ошибка:", err)
	//}

	// Преобразуем из объекта именно окно типа gtk.Window
	// и соединяем с сигналом "destroy" чтобы можно было закрыть
	// приложение при закрытии окна
	//win := obj.(*gtk.Window)
	//win.Connect("destroy", func() {
	//	gtk.MainQuit()
	//})

	// Отображаем все виджеты в окне
	//win.ShowAll()

	// Выполняем главный цикл GTK (для отрисовки). Он остановится когда
	// выполнится gtk.MainQuit()
	gtk.Main()
}

func initConnDialog(builder *gtk.Builder) *gtk.Dialog {
	connectDialog, err := builder.GetObject("connectDialog")

	dlg := connectDialog.(*gtk.Dialog)
	dlg.Connect("destroy", func() {
		gtk.MainQuit()
	})
}
