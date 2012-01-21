package main

import (
	"fmt"
	t "github.com/bytbox/termbox-go"
	"log"
)

var (
	width  int
	height int

	messageList []Message
)

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
				message <- e
			} else {
				message <- "unknown error"
			}
		}
	}()
	a()
}

func updateSize() {
	width = t.Width()
	height = t.Height()
}

// Channels for signaling the main routine
var (
	updates = make(chan interface{}) // updates the screen
	message = make(chan string) // updates the message displayed
)

type Display struct {
	message string
}

func updateMessages() {
	message <- "updating..."
	messageList = UpdateAllList()
	message <- ""
}

func display(d Display) {
	t.Clear()

	// List all messages
	for i, m := range messageList {
		w := t.Writer(2, i, t.WHITE, t.BLACK)
		fmt.Fprintf(w, "%s", m.Title)
	}

	// Write the error message or otherwise
	messageWriter := t.Writer(0, height-1, t.WHITE, t.BLACK)
	fmt.Fprint(messageWriter, d.message)

	t.Present()
}

func UIMain() {
	t.Init()
	defer func() {
		t.Shutdown()
	}()
	updateSize()

	events := make(chan t.Event)

	go func() {
		for {
			e := t.Event{}
			e.Poll()
			events <- e
		}
	}()

	d := Display{}

	for {
		display(d)
		select {
		case e := <-events:
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
					go act(a)
				} else {
					go act(unknownAction)
				}
			case t.EVENT_RESIZE:
				updateSize()
			default:
				log.Print("warning: unknown event type")
			}
		case <-updates:
		case m := <-message:
			d.message = m
		}
	}

Exit:
}
