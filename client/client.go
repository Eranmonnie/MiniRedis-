package client

import (
	"bytes"
	"context"
	"net"

	"github.com/tidwall/resp"
)

type Client struct {
	addr string
	conn net.Conn
}

func New(address string) (*Client, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	return &Client{
		addr: address,
		conn: conn,
	}, nil
}

func (c *Client) Set(ctx context.Context, key string, val string) error {

	var buf bytes.Buffer
	wr := resp.NewWriter(&buf)
	wr.WriteArray([]resp.Value{resp.StringValue("set"), resp.StringValue(key), resp.StringValue(val)})

	_, err := c.conn.Write(buf.Bytes())
	return err
}

func (c *Client) Get(ctx context.Context, key string) (string, error) {

	var buf bytes.Buffer
	wr := resp.NewWriter(&buf)
	wr.WriteArray([]resp.Value{resp.StringValue("get"), resp.StringValue(key)})

	_, err := c.conn.Write(buf.Bytes())
	if err != nil {
		return "", err
	}

	b := make([]byte, 1024)
	n, err := c.conn.Read(b)
	if err != nil {
		return "", err
	}

	return string(b[:n]), nil
}
