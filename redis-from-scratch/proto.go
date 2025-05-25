package main

import (
	"bytes"
	"fmt"
	"io"
	"log"

	"github.com/tidwall/resp"
)

const (
	CommandSET = "SET" // SET Command
)

// Command Interface
type Command interface{}

// Set Command Struct
type SetCommand struct {
	key, val string
}

// Parse Command
func parseCommand(raw string) (Command, error) {
	rd := resp.NewReader(bytes.NewBufferString(raw))

	for {
		v, _, err := rd.ReadValue()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Read %s\n", v.Type())

		if v.Type() == resp.Array {
			for _, value := range v.Array() {
				switch value.String() {
				case CommandSET:
					fmt.Printf("%+v\n", v.Array())
					if len(v.Array()) != 3 {
						return nil, fmt.Errorf("invalid number of variables for SET command")
					}
					cmd := SetCommand{
						key: v.Array()[1].String(),
						val: v.Array()[2].String(),
					}
					return cmd, nil
				default:
					fmt.Println("no value type matched")
				}
			}
		}
	}

	return "foo", nil
}
