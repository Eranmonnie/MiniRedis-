package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"time"

	"github.com/eranmonnie/go-redis/client"
)

const defaultListener = ":5001"

type Config struct {
	ListenAddress string
}

type Server struct {
	Config
	peers     map[*Peer]bool
	ln        net.Listener
	addPeerch chan *Peer
	quitch    chan struct{}
	msgch     chan []byte

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
		msgch:     make(chan []byte),
		kv:        NewKV(),
	}
}

func (s *Server) Start() error {
	slog.Info("server running", "listenAddress", s.ListenAddress)
	ln, err := net.Listen("tcp", s.ListenAddress)
	if err != nil {
		return err
	}
	s.ln = ln

	go s.loop()

	return s.acceptLoop()
}

func (s *Server) handleRawMessage(rawMsg []byte) error {
	cmd, err := parseCommand(string(rawMsg))
	if err != nil {
		return err
	}

	switch v := cmd.(type) {
	case SetCommand:
		slog.Info("someone wants to set a key into the hash table", "key", v.key, "value", v.val)
		return s.kv.Set(v.key, v.val)
	}

	return nil
}

func (s *Server) loop() {
	for {
		select {
		case rawMsg := <-s.msgch:
			if err := s.handleRawMessage(rawMsg); err != nil {
				slog.Error("handle raw message error", "error", err)
			}
		case <-s.quitch:
			return

		case p := <-s.addPeerch:
			s.peers[p] = true
		}
	}
}

func (s *Server) acceptLoop() error {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go s.handleConn(conn)
		//  return nil
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
	server := NewServer(Config{
		ListenAddress: ":3000",
	})

	go func() {
		log.Fatal(server.Start())
	}()

	time.Sleep(2 * time.Second)

	client := client.New("localhost:3000")

	for i := 0; i < 10; i++ {
		err := client.Set(context.Background(), fmt.Sprintf("leader_%d", i), fmt.Sprintf("Charlie_%d", i))
		if err != nil {
			log.Fatal(err)
		}

		time.Sleep(2 * time.Second)

		fmt.Println(server.kv.data)
	}

	select {}
}
