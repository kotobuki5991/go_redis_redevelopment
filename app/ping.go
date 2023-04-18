package main

import (
	"fmt"
	"net"
)

type PingCommand struct {}

var pingInstance Command

func NewPingInstance() Command {
	if pingInstance == nil {
		pingInstance = &EchoCommand{}
	}
	return pingInstance
}

func (this *PingCommand)Exec(conn net.Conn, args []string) []byte {
	fmt.Println("call ping")
	return []byte(fmt.Sprint("+PONG", CRLF))
}