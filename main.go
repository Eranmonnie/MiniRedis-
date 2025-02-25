package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net"
)

const defaultListener = ":5001"

type Config struct {
	ListenAddress string
}

type Message struct {
	cmd  Command
	peer *Peer
}

type Server struct {
	Config
	peers     map[*Peer]bool
	ln        net.Listener
	addPeerch chan *Peer
	quitch    chan struct{}
	msgch     chan Message

	kv *KV
}

func NewServer(cfg Config) *Server {
	if len(cfg.ListenAddress) == 0 {
		cfg.ListenAddress = defaultListener
	}
	return &Server{
		Config:    cfg,
		peers:     make(map[*Peer]bool),
		addPeerch: make(chan *Peer),
		quitch:    make(chan struct{}),
		msgch:     make(chan Message),
		kv:        NewKV(),
	}
}

func (s *Server) Start() error {
	slog.Info("go-redis server running", "listenAddress", s.ListenAddress)
	ln, err := net.Listen("tcp", s.ListenAddress)
	if err != nil {
		return err
	}
	s.ln = ln

	go s.loop()

	return s.acceptLoop()
}

func (s *Server) handleMessage(msg Message) error {
	switch v := msg.cmd.(type) {
	case SetCommand:
		// slog.Info("someone wants to set a key into the hash table", "key", v.key, "value", v.val)
		return s.kv.Set(v.key, v.val)
	case GetCommand:
		// slog.Info("someone wants to get a key from the hash table", "key", v.key)
		val, ok := s.kv.Get(v.key)

		if !ok {
			return fmt.Errorf("key not found")
		}
		_, err := msg.peer.Send(val)
		if err != nil {
			slog.Error("error sending value to peer", "error", err)
		}
	}

	return nil
}

func (s *Server) loop() {
	for {
		select {
		case msg := <-s.msgch:
			if err := s.handleMessage(msg); err != nil {
				slog.Error("handle raw message error", "error", err)
			}
		case <-s.quitch:
			return

		case p := <-s.addPeerch:
			slog.Info("new peer added", "remoteAddress", p.conn.RemoteAddr())
			s.peers[p] = true
		}
	}
}

func (s *Server) acceptLoop() error {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			select {
			case <-s.quitch:
				return nil
			default:
				fmt.Println(err)
				continue
			}
		}
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	p := NewPeer(conn, s.msgch)
	s.addPeerch <- p
	slog.Info("new peer connected", "remoteAddress", conn.RemoteAddr())
	if err := p.readLoop(); err != nil {
		slog.Error("peer read error", "error", err, "remoteAddress", conn.RemoteAddr())
	}

}

func main() {

	listenAddr := flag.String("listenAddr", defaultListener, "listen address of the go-redis server")
	flag.Parse()
	server := NewServer(Config{
		ListenAddress: *listenAddr,
	})
	log.Fatal(server.Start())
}
