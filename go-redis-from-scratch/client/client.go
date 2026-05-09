package client

import (
	"bytes"
	"context"
	"log"
	"net"

	"github.com/tidwall/resp"
)

// Client Struct
type Client struct {
	addr string   // client address
	conn net.Conn // client connection
}

// Initialize New Client
func New(addr string) *Client {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	return &Client{
		addr: addr,
		conn: conn,
	}
}

// Set Key and Value
func (c *Client) Set(ctx context.Context, key, val string) error {
	conn, err := net.Dial("tcp", c.addr)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	wr := resp.NewWriter(&buf)
	wr.WriteArray(
		[]resp.Value{
			resp.StringValue("SET"),
			resp.StringValue(key),
			resp.StringValue(val),
		},
	)

	_, err = conn.Write(buf.Bytes())
	return err
}

// Get Value with Key
func (c *Client) Get(ctx context.Context, key string) (string, error) {
	conn, err := net.Dial("tcp", c.addr)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	wr := resp.NewWriter(&buf)
	wr.WriteArray(
		[]resp.Value{
			resp.StringValue("GET"),
			resp.StringValue(key),
		},
	)
	_, err = conn.Write(buf.Bytes())
	if err != nil {
		return "", err
	}

	b := make([]byte, 1024)
	n, err := conn.Read(b)

	return string(b[:n]), err
}
