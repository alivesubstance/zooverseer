package main

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"log"
)

func getObject(builder *gtk.Builder, objectName string) glib.IObject {
	object, err := builder.GetObject(objectName)
	checkError(err)

	return object
}

func checkError(e error) {
	if e != nil {
		log.Panic(e)
	}
}
