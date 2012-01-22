package main

import (
	"fmt"
	t "github.com/bytbox/termbox-go"
	. "github.com/bytbox/go-mail"
	"log"
)

var (
	width  int
	height int
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
			if e, ok := er.(string); ok {
				message <- e
			} else if e, ok := er.(error); ok {
				message <- e.Error()
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
	escape = make(chan interface{})  // exits safely
	updates = make(chan interface{}) // updates the screen
	message = make(chan string)      // updates the message displayed
)

type Display struct {
	message string
}

func lim(s string, i int) string {
	if len(s) > i {
		return s[:i]
	}
	return s
}

func updateMessages() {
	// TODO do this in a way that doesn't risk race conditions
	message <- "updating..."
	mc := UpdateAll()
	i := 0
	for m := range mc {
		store.Add(m)
		store.messageList = append(store.messageList, m)
		i++
	}
	message <- fmt.Sprintf("read %d messages", i)
}

func lookup(m Message, k string) string {
	// TODO this is a hack
	for _, h := range m.RawHeaders {
		if k == h.Key {
			return h.Value
		}
	}
	return "---"
}

func display(d Display) {
	t.Clear()

	// Calculate tab-stops
	tabs := []int{
		1,
		2 * width / 8,
		4 * width / 8,
		7 * width / 8,
		width,
	}

	headers := []string{
		"From",
		"To",
		"Subject",
	}

	// Headers
	for i, h := range headers {
		t.WriteAt(tabs[i], 0, h, t.GREEN, t.BLACK)
	}

	// Message
	for i, m := range store.messageList {
		for j, h := range headers {
			c := lim(lookup(m, h), tabs[j+1] - tabs[j] - 1)
			t.WriteAt(tabs[j], i+1, c, t.WHITE, t.BLACK)
		}
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
		defer func() {
			escape <- nil
		}()
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
		case <-escape:
			goto Exit
		}
	}

Exit:
}
