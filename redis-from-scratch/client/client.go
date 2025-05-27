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
	buf := &bytes.Buffer{}
	wr := resp.NewWriter(buf)

	wr.WriteArray(
		[]resp.Value{
			resp.StringValue("SET"),
			resp.StringValue(key),
			resp.StringValue(val),
		},
	)

	_, err := c.conn.Write(buf.Bytes())
	buf.Reset()
	return err
}
