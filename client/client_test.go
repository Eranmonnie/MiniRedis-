package client

import (
	"context"
	"fmt"
	"log"
	"testing"
)

func TestNewClient1(t *testing.T) {
	client, err := New("localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	err = client.Set(context.Background(), "foo", 1)
	if err != nil {
		log.Fatal(err)
	}

	val, err := client.Get(context.Background(), "foo")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("value is %s \n", string(val))

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
