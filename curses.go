package main

import (
	t "github.com/nsf/termbox-go"
	"log"
)

var (
	width  int
	height int

	message string
)

var chActions = map[rune]func(){}

var keyActions = map[uint16]func(){}

func unknownAction() {
	panic("unrecognized keystroke")
}

func act(a func()) {
	defer func() {
		if er := recover(); er != nil {
			e, ok := er.(error)
			if ok {
				message = e.Error()
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

func display() {

}

func UIMain() {
	t.Init()
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
	t.Shutdown()
}
