package client

import (
	"context"
	"fmt"
	"log"
	"sync"
	"testing"
)

func TestNewClients(t *testing.T) {
	nClients := 10
	wg := sync.WaitGroup{}
	wg.Add(nClients)
	for i := 0; i < nClients; i++ {
		go func(it int) {
			c, err := New("localhost:8080")
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
			fmt.Printf("client_%d, got this value back  => %s", i, val)
			wg.Done()
		}(i)

	}
	wg.Wait()
}

func TestNewCient(t *testing.T) {

	client, err := New("localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	for i := 0; i < 10; i++ {
		err := client.Set(context.Background(), fmt.Sprintf("leader_%d", i), fmt.Sprintf("Charlie_%d", i))
		if err != nil {
			log.Fatal(err)
		}

		val, err := client.Get(context.Background(), fmt.Sprintf("leader_%d", i))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("value is %s \n", string(val))
	}
}
