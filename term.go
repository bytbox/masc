package main

import (
	"os"
	"os/exec"
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
	'r': readMessage,
	'R': replyMessage,
	'u': updateMessages,
	'm': writeMessage,
}

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
	execProg    = make(chan []string)    // execProgutes program with arguments
)

// display information
var (
	d_message   string
	d_selected  int // this is the row of the currently selected message
	d_screenpos int
)

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

func readMessage() {
	tmpFile := mkTemp(store.messageList[d_selected].Body)
	execProg <- []string{"less", tmpFile}
}

func replyMessage() {
	execProg <- []string{"vim"}
}

func writeMessage() {
	execProg <- []string{"vim"}
}

func display() {
	t.Clear()

	// Calculate tab-stops
	tabs := []int{
		1,
		3 * width / 16,
		6 * width / 16,
		13 * width / 16,
		width,
	}

	headers := []string{
		"From",
		"To",
		"Subject",
		"Date",
	}

	// Headers
	for i, h := range headers {
		t.WriteAt(tabs[i], 0, h, t.GREEN, t.BLACK)
	}

	// Message
	for
		i := d_screenpos;
		i < min(height-2+d_screenpos, len(store.messageList));
		i++ {
		m := store.messageList[i]
		var fg = uint16(t.WHITE)
		var bg = uint16(t.BLACK)
		hp := i-d_screenpos+1
		if i == d_selected {
			bg = t.RED
		}
		for x := 0; x < width; x++ {
			t.ChangeCell(x, hp, ' ', fg, bg)
		}
		for j, h := range headers {
			c := lim(lookup(m, h), tabs[j+1] - tabs[j] - 1)
			t.WriteAt(tabs[j], hp, c, fg, bg)
		}
	}

	// Write the error message or otherwise
	messageWriter := t.Writer(0, height-1, t.WHITE, t.BLACK)
	fmt.Fprint(messageWriter, d_message)

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

	for {
		display()
		select {
		case e := <-events:
			switch e.Type {
			case t.EVENT_KEY:
				if e.Ch == 'q' {
					goto Exit
				}
				if e.Key == t.KEY_ARROW_UP || e.Ch == 'k'{
					if d_selected > 0 {
						d_selected--
					}
				} else if e.Key == t.KEY_ARROW_DOWN || e.Ch == 'j' {
					if d_selected < len(store.messageList)-1 {
						d_selected++
					}
				} else {
					a, ok := chActions[e.Ch]
					if ok {
						go act(a)
					} else {
						go act(unknownAction)
					}
				}
				if d_selected < d_screenpos+1 {
					d_screenpos = max(d_selected-1, 0)
				}
				if d_selected > d_screenpos + height-4 {
					d_screenpos++
				}
			case t.EVENT_RESIZE:
				updateSize()
			default:
				log.Print("warning: unknown event type")
			}
		case <-updates:
		case p := <-execProg:
			t.Shutdown()
			cmd := p[0]
			args := p[1:]
			c := exec.Command(cmd, args...)
			c.Stdin = os.Stdin
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			e := c.Start()
			if e != nil { panic(e) }
			e = c.Wait()
			if e != nil { panic(e) }
			t.Init()
			updateSize()
		case m := <-message:
			d_message = m
		case <-escape:
			goto Exit
		}
	}

Exit:
}
