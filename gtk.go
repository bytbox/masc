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

	container := gtk.VBox(false, 5)

	// toolbar
	toolbar := gtk.HBox(false, 1)
	sendButton := gtk.ButtonWithLabel("Send")
	toolbar.PackStart(sendButton, false, false, 0)
	//container.Add(toolbar)
	container.PackStart(toolbar, false, false, 0)

	// message content
	swin := gtk.ScrolledWindow(nil, nil)
	swin.SetPolicy(gtk.GTK_POLICY_AUTOMATIC, gtk.GTK_POLICY_ALWAYS)
	contentField := gtk.TextView()
	swin.Add(contentField)
	container.Add(swin)

	window.Add(container)

	window.SetSizeRequest(800, 600)
	window.ShowAll()
}

func UIMain() {
	window := gtk.Window(gtk.GTK_WINDOW_TOPLEVEL)
	window.SetPosition(gtk.GTK_WIN_POS_CENTER)
	window.SetTitle("Masc")

	window.Connect("destroy", func(ctx interface{}) {
		gtk.MainQuit()
	}, nil)

	outer := gtk.VBox(false, 4)

	// toolbar
	toolbar := gtk.HBox(false, 1)

	container := gtk.HBox(true, 4)

	// message list
	lScroll := gtk.ScrolledWindow(nil, nil)
	lScroll.SetPolicy(gtk.GTK_POLICY_AUTOMATIC, gtk.GTK_POLICY_ALWAYS)
	container.Add(lScroll)

	// message body
	bScroll := gtk.ScrolledWindow(nil, nil)
	bScroll.SetPolicy(gtk.GTK_POLICY_AUTOMATIC, gtk.GTK_POLICY_ALWAYS)
	bodyView := gtk.TextView()
	bodyView.SetEditable(false)
	bScroll.Add(bodyView)
	container.Add(bScroll)

	outer.PackStart(toolbar, false, false, 0)
	outer.Add(container)
	window.Add(outer)

	window.SetSizeRequest(800, 600)
	window.ShowAll()

	gtk.Main()
}
