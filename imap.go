package main

import (
	tp "net/textproto"
)

type Client struct {
	Text *tp.Conn
}
