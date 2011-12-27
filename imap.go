package main

import (
	"crypto/tls"
	"net"
	tp "net/textproto"
	"strings"
)

type Client struct {
	Text *tp.Conn
	conn net.Conn
	host string
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
		conn: conn,
		host: host,
	}
	return client, nil
}

func (c *Client) cmd() {

}

// Noop sends a NOOP command to the server.
func (c *Client) Noop() {

}

// Login authenticates a client using the provided username and password. This
// method is only secure if TLS is being used.
func (c *Client) Login(username, password string) {

}

// Logout un-authenticates a client.
func (c *Client) Logout() {

}
