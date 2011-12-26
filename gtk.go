package main

import (
	"github.com/mattn/go-gtk/gtk"
)

func init() {
	gtk.Init(nil)
}

// Open the composition window.
func Compose() {
	window := gtk.Window(gtk.GTK_WINDOW_TOPLEVEL)
	window.SetPosition(gtk.GTK_WIN_POS_CENTER)
	window.SetTitle("Compose Message")

	window.ShowAll()
}

func GUIMain() {
	gtk.Main()
}

/*
	gtk.Init(nil)

	window := gtk.Window(gtk.GTK_WINDOW_TOPLEVEL)
	window.SetPosition(gtk.GTK_WIN_POS_CENTER)
	window.SetTitle("GTK Go!")

	window.Connect("destroy", func(ctx interface{}) {
		gtk.MainQuit()
	}, nil)

	window.SetSizeRequest(600, 600)
	window.ShowAll()

	gtk.Main()
*/
