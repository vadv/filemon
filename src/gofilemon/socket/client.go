package socket

import (
	"net"
	"time"
)

const (
	MaxSizeRequest = 1024
)

type Client struct {
	conn net.Conn
}

func NewClient() (*Client, error) {
	conn, err := net.DialTimeout("unix", FilePath(), 200*time.Millisecond)
	if err != nil {
		return nil, err
	}
	result := &Client{conn: conn}
	return result, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) Write(data []byte) (err error) {
	c.conn.SetWriteDeadline(time.Now().Add(100 * time.Millisecond))
	_, err = c.conn.Write(data)
	return
}

func (c *Client) Read() ([]byte, error) {
	buf := make([]byte, MaxSizeRequest)
	c.conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
	n, err := c.conn.Read(buf[:])
	if err != nil {
		return nil, err
	}
	return buf[0:n], nil
}
