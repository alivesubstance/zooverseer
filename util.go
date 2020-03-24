package main

import (
	"github.com/gotk3/gotk3/glib"
	"log"
)

func getObject(objectName string) glib.IObject {
	object, err := Builder.GetObject(objectName)
	checkError(err)

	return object
}

func checkError(e error) {
	if e != nil {
		log.Panic(e)
	}
}
