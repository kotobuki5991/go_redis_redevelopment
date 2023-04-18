package main

import (
	"fmt"
)

type PingCommand struct {}

var pingInstance Command

func NewPingInstance() Command {
	if pingInstance == nil {
		pingInstance = &PingCommand{}
	}
	return pingInstance
}

func (this *PingCommand)Exec(args []string) []byte {
	fmt.Println("call ping")
	return []byte(fmt.Sprint("+PONG", CRLF))
}
