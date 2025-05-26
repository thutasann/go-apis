package client

import (
	"bytes"
	"context"
	"io"
	"net"

	"github.com/tidwall/resp"
)

// Client Struct
type Client struct {
	addr string // client address
}

// Initialize New Client
func New(addr string) *Client {
	return &Client{
		addr: addr,
	}
}

// Set Key and Value
func (c *Client) Set(ctx context.Context, key, val string) error {
	conn, err := net.Dial("tcp", c.addr)
	if err != nil {
		return err
	}

	buf := &bytes.Buffer{}
	wr := resp.NewWriter(buf)
	wr.WriteArray(
		[]resp.Value{
			resp.StringValue("SET"),
			resp.StringValue(key),
			resp.StringValue(val),
		},
	)

	_, err = io.Copy(conn, buf)
	return err
}
