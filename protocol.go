package main

import (
	"bytes"
	"fmt"

	"github.com/tidwall/resp"
)

const (
	commandSet    = "set"
	commandGet    = "get"
	commandHello  = "hello"
	commandClient = "client"
)

type Command interface {
}

type SetCommand struct {
	key, val []byte
}

type ClientCommand struct {
	val string
}

type HelloCommand struct {
	val string
}
type GetCommand struct {
	key []byte
}

func respWriteMap(m map[string]string) []byte {
	buf := &bytes.Buffer{}
	buf.WriteString("%" + fmt.Sprintf("%d\r\n", len(m)))
	rw := resp.NewWriter(buf)
	for k, v := range m {
		rw.WriteString(k)
		rw.WriteString(":" + v)
	}
	return buf.Bytes()
}

// func parseCommand(raw string) (Command, error) {
// 	rd := resp.NewReader(bytes.NewBufferString(raw))
// 	for {
// 		v, _, err := rd.ReadValue()
// 		if err == io.EOF {
// 			break
// 		}
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		if v.Type() == resp.Array {
// 			for _, value := range v.Array() {
// 				switch value.String() {
// 				case commandSet:
// 					if len(v.Array()) != 3 {
// 						return nil, fmt.Errorf("invalid number of variables for %s command", commandSet)
// 					}
// 					cmd := SetCommand{
// 						key: v.Array()[1].Bytes(),
// 						val: v.Array()[2].Bytes(),
// 					}
// 					return cmd, nil
// 				case commandGet:
// 					if len(v.Array()) != 2 {
// 						return nil, fmt.Errorf("invalid number of variables for %s command", commandGet)
// 					}
// 					cmd := GetCommand{
// 						key: v.Array()[1].Bytes(),
// 						// val: v.Array()[2].Bytes(),
// 					}
// 					return cmd, nil

// 				case commandHello:
// 					cmd := HelloCommand{
// 						val: v.Array()[1].String(),
// 					}
// 					return cmd, nil
// 				}
// 			}
// 		}
// 	}
// 	return nil, fmt.Errorf("invalid  or unknown command recieved, %s", raw)
// }
