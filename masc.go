package main

import (
	"flag"

	"github.com/mattn/go-gtk/gtk"
)

func main() {
	flag.Parse()

	gtk.Init(nil)

	window := gtk.Window(gtk.GTK_WINDOW_TOPLEVEL)
	window.SetPosition(gtk.GTK_WIN_POS_CENTER)
	window.SetTitle("GTK Go!")

	window.SetSizeRequest(600, 600)
	window.ShowAll()

	gtk.Main()
}
