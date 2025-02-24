package client

import (
	"context"
	"fmt"
	"log"
	"testing"
)

func TestNewCient(t *testing.T) {

	client, err := New("localhost:8080")
	if err != nil {
		log.Fatal(err)
	}

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
