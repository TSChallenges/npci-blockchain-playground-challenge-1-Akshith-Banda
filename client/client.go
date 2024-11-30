package client

import "net"

// Client represents a chat user
type Client struct {
	Conn     net.Conn
	Username string
}
