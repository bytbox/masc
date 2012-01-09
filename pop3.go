// Package pop3 provides an implementation of the Post Office Protocol, Version
// 3 as defined in RFC 1939.

package main

import (
	"bufio"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
)

// The POP3 client.
type Client struct {
	conn net.Conn
	bin  *bufio.Reader
}

// Dial creates an unsecured connection to the POP3 server at the given address
// and returns the corresponding Client.
func Dial(addr string) (*Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return NewClient(conn)
}

// DialTLS creates a TLS-secured connection to the POP3 server at the given
// address and returns the corresponding Client.
func DialTLS(addr string) (*Client, error) {
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		return nil, err
	}
	return NewClient(conn)
}

// NewClient returns a new Client object using an existing connection.
func NewClient(conn net.Conn) (*Client, error) {
	client := &Client{
		bin: bufio.NewReader(conn),
		conn: conn,
	}
	// send dud command, to read a line
	_, err := client.cmd("")
	if err != nil { return nil, err }
	return client, nil
}

// Convenience function to synchronously run an arbitrary command and wait for
// output. The terminating CRLF must be included in the format string.
func (c *Client) cmd(format string, args ...interface{}) (string, error) {
	fmt.Fprintf(c.conn, format, args...)
	line, _, err := c.bin.ReadLine()
	l := string(line)
	if l[0:3] != "+OK" {
		err = errors.New(l[5:])
	}
	return l[4:], err
}

// User sends the given username to the server. Generally, there is no reason
// not to use the Auth convenience method.
func (c *Client) User(username string) (err error) {
	_, err = c.cmd("USER %s\r\n", username)
	return
}

// Pass sends the given password to the server. The password is sent
// unencrypted unless the connection is already secured by TLS (via DialTLS or
// some other mechanism). Generally, there is no reason not to use the Auth
// convenience method.
func (c *Client) Pass(password string) (err error) {
	_, err = c.cmd("PASS %s\r\n", password)
	return
}

// Auth sends the given username and password to the server, calling the User
// and Pass methods as appropriate.
//
// Technically speaking, the server may opt not to support either
// authentication mechanism; however, in practice, all implement both.
func (c *Client) Auth(username, password string) (err error) {
	err = c.User(username)
	if err != nil { return }
	err = c.Pass(password)
	return
}

// Apop sends the given username and password hash (MD5 digest) to the server.
// This method does not offer any more real security over Auth.
func (c *Client) Apop(username, digest string) (err error) {
	_, err = c.cmd("APOP %s %s\r\n", username, digest)
	return
}

// Stat retrieves a drop listing for the current maildrop, consisting of the
// number of messages and the total size (in octets) of the maildrop.
// Information provided besides the number of messages and the size of the
// maildrop is ignored. In the event of an error, all returned numeric values
// will be 0.
func (c *Client) Stat() (count, size int, err error) {
	l, err := c.cmd("STAT\r\n")
	if err != nil { return 0, 0, err }
	parts := strings.Fields(l)
	count, err = strconv.Atoi(parts[0])
	if err != nil { return 0, 0, errors.New("Invalid server response") }
	size, err = strconv.Atoi(parts[1])
	if err != nil { return 0, 0, errors.New("Invalid server response") }
	return
}

// List returns the size of the given message, if it exists. If the message
// does not exist, or another error is encountered, the returned size will be
// 0.
func (c *Client) List(msg int) (size int, err error) {
	l, err := c.cmd("LIST %d\r\n", msg)
	if err != nil { return 0, err }
	size, err = strconv.Atoi(strings.Fields(l)[1])
	if err != nil { return 0, errors.New("Invalid server response") }
	return size, nil
}

// Quit sends the QUIT message to the POP3 server and closes the connection.
func (c *Client) Quit() error {
	_, err := c.cmd("QUIT\r\n")
	if err != nil { return err }
	c.conn.Close()
	return nil
}