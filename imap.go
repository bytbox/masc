// Package imap partially implements the Internet Message Access Protocol as
// defined in RFC 3501.
//
// Untagged IMAP responses are parsed into a Mailbox struct, which tracks all
// currently known information concerning the state of the remote mailbox.

package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	tp "net/textproto"
	"strings"
	"sync"
)

type Client struct {
	// The underlying textproto connection may be used to extend the
	// functionality of this package; however, using Client.Cmd instead is
	// recommended, as doing so better preserves thread-safety.
	Text *tp.Conn

	Box  *Mailbox

	tags map[string]chan string
	tMut *sync.Mutex
}

// Represents the current known state of the remote mailbox.
type Mailbox struct {
	capabilities []string
	mut          *sync.RWMutex
}

func (m *Mailbox) Capable(c string) bool {
	m.mut.RLock()
	for _, ca := range m.capabilities {
		if c == ca {
			return true
		}
	}
	m.mut.RUnlock()
	return false
}

// Dial creates an unsecured connection to the IMAP server at the given address
// and returns the corresponding Client.
func Dial(addr string) (*Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return NewClient(conn)
}

// DialTLS creates a TLS_secured connection to the IMAP server at the given
// address and returns the corresponding Client.
func DialTLS(addr string) (*Client, error) {
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		return nil, err
	}
	return NewClient(conn)
}

// NewClient returns a new Client using an existing connection.
func NewClient(conn net.Conn) (*Client, error) {
	text := tp.NewConn(conn)
	client := &Client{
		Text: text,
		Box:  &Mailbox{
			capabilities: []string{},
			mut: new(sync.RWMutex),
		},
		tags: map[string]chan string{},
		tMut: new(sync.Mutex),
	}

	input := make(chan string)

	// Read all input from conn
	go func() {
		l, err := text.ReadLine()
		for err == nil {
			input <- l
			l, err = text.ReadLine()
		}
		if err == io.EOF {
			close(input)
		} else {
			panic(err)
		}
	}()

	// Start the serving goroutine
	go func() {
		for {
			select {
			case l := <-input:
				if isUntagged(l) {
					client.handleUntagged(l[2:])
					continue
				}
				// handle tagged response
				ps := strings.SplitN(l, " ", 2)
				tag := ps[0]
				l = ps[1]
				client.tMut.Lock()
				client.tags[tag] <- l
				client.tMut.Unlock()
			}
		}
	}()
	return client, nil
}

func (c *Client) handleUntagged(l string) {
	c.Box.mut.Lock()
	switch l[0:strings.Index(l, " ")] {
	case "CAPABILITY":
		c.Box.capabilities = strings.Split(l, " ")[1:]
	default:
		println(l)
	}
	c.Box.mut.Unlock()
}

func (c *Client) Cmd(format string, args ...interface{}) error {
	t := c.Text
	id := t.Next()
	tag := fmt.Sprintf("x%d", id)
	t.StartRequest(id)
	err := t.PrintfLine("%s %s", tag, fmt.Sprintf(format, args...))
	if err != nil {
		return err
	}
	t.EndRequest(id)

	t.StartResponse(id)
	defer t.EndResponse(id)

	ch := make(chan string)
	c.tMut.Lock()
	c.tags[tag] = ch
	c.tMut.Unlock()

	l := <-ch
	if l[0:2] == "OK" {
		return nil
	}
	return errors.New(l)
}

func isUntagged(l string) bool {
	return l[0:2] == "* "
}

// Noop sends a NOOP command to the server, which may be abused to test that
// the connection is still working, or keep it active.
func (c *Client) Noop() error {
	return c.Cmd("NOOP")
}

// Capability determines the server's capabilities.
func (c *Client) Capability() error {
	return c.Cmd("CAPABILITY")
}

// Login authenticates a client using the provided username and password. This
// method is only secure if TLS is being used. AUTHENTICATE and STARTTLS are
// not supported.
func (c *Client) Login(username, password string) error {
	return c.Cmd("LOGIN %s %s", username, password)
}

// Logout closes the connection, after instructing the server to do the same.
func (c *Client) Logout() error {
	return c.Cmd("LOGOUT")
}

// SELECT

// EXAMINE

// CREATE

// DELETE

// RENAME

// SUBSCRIBE

// UNSUBSCRIBE

// List lists all folder within basename that match the wildcard expression mb.
// The result is put into the Client's Mailbox struct.
func (c *Client) List(basename, mb string) error {
	return c.Cmd(`LIST "%s" "%s"`, basename, mb)
}

// LSUB

// STATUS

// APPEND

// CHECK

// CLOSE

// EXPUNGE

// SEARCH

// FETCH

// STORE

// COPY

// UID

