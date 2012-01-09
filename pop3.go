// Package pop3 implements the Post Office Protocol, Version 3 as defined in
// RFC 1939.

package main

import (
	"bufio"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
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
// output.
func (c *Client) cmd(format string, args ...interface{}) (string, error) {
	fmt.Fprintf(c.conn, format, args...)
	line, _, err := c.bin.ReadLine()
	l := string(line)
	if l[0:3] != "+OK" {
		err = errors.New(l[5:])
	}
	return l, err
}

// Quit sends the QUIT message to the POP3 server and closes the connection.
func (c *Client) Quit() error {
	_, err := c.cmd("QUIT\n")
	if err != nil { return err }
	c.conn.Close()
	return nil
}
