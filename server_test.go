package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/eranmonnie/go-redis/client"
)

func TestServerWithClients(t *testing.T) {
	server := NewServer(Config{
		ListenAddress: ":8080",
	})
	go func() {
		log.Fatal(server.Start())
	}()

	time.Sleep(2 * time.Second)

	nClients := 10
	wg := sync.WaitGroup{}
	wg.Add(nClients)
	for i := 0; i < nClients; i++ {
		go func(it int) {
			c, err := client.New("localhost:8080")
			if err != nil {
				log.Fatal(err)
			}
			defer c.Close()
			key := fmt.Sprintf("client_foo_%d", i)
			val := fmt.Sprintf("client_bar_%d", i)
			if err := c.Set(context.TODO(), key, val); err != nil {
				log.Fatal(err)
			}
			val, err = c.Get(context.TODO(), key)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("client_%d, got this value back  => %s \n", i, val)
			wg.Done()
		}(i)

	}
	wg.Wait()

	time.Sleep(2 * time.Second)

	if len(server.peers) != 0 {
		t.Fatalf("expected 0 peers, got %d", len(server.peers))
	}

}
