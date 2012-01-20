package main

import (
	"fmt"
	t "github.com/bytbox/termbox-go"
	"log"
)

var (
	width  int
	height int

	message string
)

var messageList = []Message{
	Message{
		To: []string{"hey"},
		Title: "TITLE HERE",
		From: "there",
		Content: "hoho",
	},
}

var chActions = map[rune]func(){
	'u': updateMessages,
}

var keyActions = map[uint16]func(){}

func unknownAction() {
	panic("unrecognized keystroke")
}

func act(a func()) {
	defer func() {
		if er := recover(); er != nil {
			e, ok := er.(string)
			if ok {
				message = e
			} else {
				message = "unknown error"
			}
		}
	}()
	a()
}

func updateSize() {
	width = t.Width()
	height = t.Height()
}

func updateMessages() {
	panic("not yet implemented")
}

func display() {
	t.Clear()

	// List all messages
	for i, m := range messageList {
		w := t.Writer(2, i, t.WHITE, t.BLACK)
		fmt.Fprintf(w, "%s", m.Title)
	}

	// Write the error message or otherwise
	messageWriter := t.Writer(0, height-1, t.WHITE, t.BLACK)
	fmt.Fprint(messageWriter, message)

	t.Present()
}

func UIMain() {
	t.Init()
	defer func() {
		t.Shutdown()
	}()
	updateSize()
	e := t.Event{}
	for {
		display()
		e.Poll()
		switch e.Type {
		case t.EVENT_KEY:
			if e.Ch == 'q' {
				goto Exit
			}
			a, ok := chActions[e.Ch]
			if !ok {
				a, ok = keyActions[e.Key]
			}
			if ok {
				act(a)
			} else {
				act(unknownAction)
			}
		case t.EVENT_RESIZE:
			updateSize()
		default:
			log.Print("warning: unknown event type")
		}
	}

Exit:
}
