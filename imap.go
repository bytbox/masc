// Package imap implements the Internet Message Access Protocol as defined in
// RFC 3501.
//
// This implementation is thread-safe, but does not take advantage of all the
// parallelism provided for in the standard. Commands will be performed
// sequentially.
//
// Untagged IMAP responses are parsed into a Mailbox struct, which tracks all
// currently known information concerning the state of the remote mailbox.

package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	tp "net/textproto"
	"strings"
)

type Client struct {
	Text *tp.Conn
	Box  *Mailbox
	conn net.Conn
	host string
}

// Represents the current known state of the remote mailbox.
type Mailbox struct {
	Capabilities []string
}

func Dial(addr string) (*Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	host := addr[:strings.Index(addr, ":")]
	return NewClient(conn, host)
}

func DialTLS(addr string) (*Client, error) {
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		return nil, err
	}
	host := addr[:strings.Index(addr, ":")]
	return NewClient(conn, host)
}

// NewClient returns a new Client using an existing connection.
func NewClient(conn net.Conn, host string) (*Client, error) {
	text := tp.NewConn(conn)
	client := &Client{
		Text: text,
		Box:  &Mailbox{
			Capabilities: []string{},
		},
		conn: conn,
		host: host,
	}
	return client, nil
}

func (c *Client) handleUntagged(l string) {
	switch l[0:strings.Index(l, " ")] {
	case "CAPABILITY":
		c.Box.Capabilities = strings.Split(l, " ")[1:]
	default:
		println(l)
	}
}

func (c *Client) cmd(format string, args ...interface{}) error {
	t := c.Text
	id := t.Next()
	t.StartRequest(id)
	err := t.PrintfLine("x%d %s", id, fmt.Sprintf(format, args...))
	if err != nil {
		return err
	}
	t.EndRequest(id)

	t.StartResponse(id)
	defer t.EndResponse(id)

	l, err := t.ReadLine()
	if err != nil {
		return err
	}
	for isUntagged(l) {
		c.handleUntagged(l[2:])
		l, err = t.ReadLine()
		if err != nil {
			return err
		}
	}

	l = strings.SplitN(l, " ", 2)[1]
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
	return c.cmd("NOOP")
}

// Capability determines the server's capabilities.
func (c *Client) Capability() error {
	return c.cmd("CAPABILITY")
}

// Login authenticates a client using the provided username and password. This
// method is only secure if TLS is being used.
func (c *Client) Login(username, password string) error {
	return c.cmd("LOGIN %s %s", username, password)
}

// Logout closes the connection, after instructing the server to do the same.
func (c *Client) Logout() error {
	return c.cmd("LOGOUT")
}
