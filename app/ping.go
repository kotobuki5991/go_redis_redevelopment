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
	return []byte(fmt.Sprint("+PONG", CRLF))
}
