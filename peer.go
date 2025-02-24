package main

import (
	"net"
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
	buf := make([]byte, 1024)
	for {
		n, err := p.conn.Read(buf)
		if err != nil {
			return err
		}
		msgBuf := make([]byte, n)
		copy(msgBuf, buf[:n])
		p.msgch <- Message{
			data: msgBuf,
			peer: p,
		}
	}
}

func (p *Peer) Send(b []byte) (int, error) {
	return p.conn.Write(b)
}
