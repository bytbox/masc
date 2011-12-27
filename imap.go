package main

import (
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

func NewClient(conn net.Conn, host string) (*Client, error) {
	text := tp.NewConn(conn)
	client := &Client{
		Text: text,
		conn: conn,
		host: host,
	}
	return client, nil
}
