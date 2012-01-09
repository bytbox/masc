package main

import (
	"gobject/gtk-3.0"
)

func init() {
	gtk.Init(nil)
}

// Open the composition window.
func Compose() {
	window := gtk.NewWindow(gtk.WindowTypeToplevel)
	window.SetPosition(gtk.WindowPositionCenter)
	window.SetTitle("Compose Message")

	container := gtk.NewVBox(false, 5)

	// toolbar
	toolbar := gtk.NewHBox(false, 1)
	sendButton := gtk.NewButtonWithLabel("Send")
	toolbar.PackStart(sendButton, false, false, 0)
	//container.Add(toolbar)
	container.PackStart(toolbar, false, false, 0)

	// message content
	swin := gtk.NewScrolledWindow(nil, nil)
	swin.SetPolicy(gtk.PolicyTypeAutomatic, gtk.PolicyTypeAlways)
	contentField := gtk.NewTextView()
	swin.Add(contentField)
	container.Add(swin)

	window.Add(container)

	window.SetSizeRequest(800, 600)
	window.ShowAll()
}

func UIMain() {
	window := gtk.NewWindow(gtk.WindowTypeToplevel)
	window.SetPosition(gtk.WindowPositionCenter)
	window.SetTitle("Masc")

	window.Connect("destroy", func(ctx interface{}) {
		gtk.MainQuit()
	})

	outer := gtk.NewVBox(false, 4)

	// toolbar
	toolbar := gtk.NewHBox(false, 1)

	container := gtk.NewHBox(true, 4)

	// message list
	lScroll := gtk.NewScrolledWindow(nil, nil)
	lScroll.SetPolicy(gtk.PolicyTypeAutomatic, gtk.PolicyTypeAlways)
	container.Add(lScroll)

	// message body
	bScroll := gtk.NewScrolledWindow(nil, nil)
	bScroll.SetPolicy(gtk.PolicyTypeAutomatic, gtk.PolicyTypeAlways)
	bodyView := gtk.NewTextView()
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
