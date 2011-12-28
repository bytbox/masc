// Package imap partially implements the Internet Message Access Protocol as
// defined in RFC 3501. Specifically, AUTHENTICATE, STARTLS, and SEARCH remain
// unimplemented.
//
// Untagged IMAP responses are parsed into a Mailbox struct, which tracks all
// currently known information concerning the state of the remote mailbox.
// Because no significant information is returned through tagged responses,
// interaction with Mailbox is necessary for all queries.

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

// The IMAP client.
type Client struct {
	// The underlying textproto connection may be used to extend the
	// functionality of this package; however, using Client.Cmd instead is
	// recommended, as doing so better preserves thread-safety.
	Text *tp.Conn

	Box  *Mailbox

	conn net.Conn // underlying raw connection.
	tags map[string]chan string
	tMut *sync.Mutex

	lit  chan string // channel where the literal string to be dumped is stored
}

// Represents the current known state of the remote server.
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
		conn: conn,
		tags: map[string]chan string{},
		tMut: new(sync.Mutex),
		lit: make(chan string),
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
				if l[0] == '+' {
					// server is ready for transmission of literal string
					client.Text.PrintfLine(<-client.lit)
					continue
				} else if isUntagged(l) {
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

// Sends a command and retreives the tagged response.
func (c *Client) Cmd(format string, args ...interface{}) error {
	c.tMut.Lock()
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
	c.tags[tag] = ch
	c.tMut.Unlock()

	l := <-ch
	if l[0:2] == "OK" {
		return nil
	}
	return errors.New(l)
}

// Equivalent to Cmd, but the first argument (which will be rotated to be the
// last) is sent as a literal string.
func (c *Client) CmdLit(lit, format string, args ...interface{}) error {
	c.tMut.Lock()
	t := c.Text
	id := t.Next()
	tag := fmt.Sprintf("x%d", id)
	t.StartRequest(id)
	err := t.PrintfLine("%s %s {%d}", tag, fmt.Sprintf(format, args...), len(lit))
	if err != nil {
		return err
	}
	t.EndRequest(id)

	c.lit <- lit

	t.StartResponse(id)
	defer t.EndResponse(id)

	ch := make(chan string)
	c.tags[tag] = ch
	c.tMut.Unlock()

	l := <-ch
	if l[0:2] == "OK" {
		return nil
	}
	return errors.New(l)
}

func isUntagged(l string) bool {
	return l[0] != 'x' // all tags are x00
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

// Select selects the specified IMAP mailbox, updating its information in the
// Mailbox object.
func (c *Client) Select(mb string) error {
	return c.Cmd(`SELECT "%s"`, mb)
}

// Examine is identical to select, but marks the mailbox read-only.
func (c *Client) Examine(mb string) error {
	return c.Cmd(`EXAMINE "%s"`, mb)
}

// Create creates the named mailbox.
func (c *Client) Create(mb string) error {
	return c.Cmd(`CREATE "%s"`, mb)
}

// Delete deletes the named mailbox.
func (c *Client) Delete(mb string) error {
	return c.Cmd(`DELETE "%s"`, mb)
}

// Rename renames the named mailbox to the new name.
func (c *Client) Rename(mb, name string) error {
	return c.Cmd(`RENAME "%s" "%s"`, mb, name)
}

// Subscribe adds the named mailbox to the list of "active" or "subscribed"
// mailboxes, to be used with Lsub .
func (c *Client) Subscribe(mb string) error {
	return c.Cmd(`SUBSCRIBE "%s"`, mb)
}

// Unsubscribe removes the named mailbox from the server's list of "active"
// mailboxes.
func (c *Client) Unsubscribe(mb string) error {
	return c.Cmd(`UNSUBSCRIBE "%s"`, mb)
}

// List lists all folder within basename that match the wildcard expression mb.
// The result is put into the Client's Mailbox struct.
func (c *Client) List(basename, mb string) error {
	return c.Cmd(`LIST "%s" "%s"`, basename, mb)
}

// Lsub is like List, but only operates on "active" mailboxes, as set with
// Subscribe and Unsubscribe.
func (c *Client) Lsub(basename, mb string) error {
	return c.Cmd(`LSUB "%s" "%s"`, basename, mb)
}

// Status queries the specified statuses of the indicated mailbox. This command
// should not be used on the currently selected mailbox. The legal status items
// are:
//
//	MESSAGES	The number of messages in the mailbox.
//	RECENT		The number of messages with the \Recent flag set.
//	UIDNEXT		The next unique identifier value of the mailbox.
//	UIDVALIDITY	The unique identifier validity value of the mailbox.
//	UNSEEN		The number of messages which do not have the \Seen flag set.
//
func (c *Client) Status(mb string, ss ...string) error {
	st := sliceAsString(ss)
	return c.Cmd(`STATUS "%s" %s`, mb, st)
}

// Append appends a message to the specified mailbox, which must exist.
//
// TODO handle flags and the optional date/time string.
func (c *Client) Append(mb, message string) error {
	return c.CmdLit(message, "APPEND \"%s\"", mb)
}

// Check tells the server to perform any necessary housekeeping.
func (c *Client) Check() error {
	return c.Cmd(`CHECK`)
}

// Close closes the selected mailbox, permanently deleting any marked messages
// in the process.
func (c *Client) Close() error {
	return c.Cmd(`CLOSE`)
}

// Expunge permanently removes all marked messages in the selected mailbox.
func (c *Client) Expunge() error {
	return c.Cmd(`EXPUNGE`)
}

// SEARCH remains unimplemented.

// FETCH

// STORE

// COPY

// UID

// Converts a slice of strings to a parenthesized list of space-separated
// strings.
func sliceAsString(ss []string) string {
	return "(" + strings.Join(ss, " ") + ")"
}
