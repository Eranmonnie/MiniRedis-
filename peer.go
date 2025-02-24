package main

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/tidwall/resp"
)

type Peer struct {
	conn  net.Conn
	msgch chan Message
}

func NewPeer(conn net.Conn, msgch chan Message) *Peer {
	return &Peer{
		conn:  conn,
		msgch: msgch,
	}
}

func (p *Peer) readLoop() error {
	rd := resp.NewReader(p.conn)

	for {
		v, _, err := rd.ReadValue()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		if v.Type() == resp.Array {
			for _, value := range v.Array() {
				switch value.String() {
				case commandSet:
					if len(v.Array()) != 3 {
						return fmt.Errorf("invalid number of variables for %s command", commandSet)
					}
					cmd := SetCommand{
						key: v.Array()[1].Bytes(),
						val: v.Array()[2].Bytes(),
					}

					p.msgch <- Message{
						peer: p,
						cmd:  cmd,
					}

				case commandGet:
					if len(v.Array()) != 2 {
						return fmt.Errorf("invalid number of variables for %s command", commandGet)
					}
					cmd := GetCommand{
						key: v.Array()[1].Bytes(),
					}
					p.msgch <- Message{
						peer: p,
						cmd:  cmd,
					}

				}
			}
		}
	}
	return nil
}

func (p *Peer) Send(b []byte) (int, error) {
	return p.conn.Write(b)
}
