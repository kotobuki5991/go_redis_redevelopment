package main

import (
	"fmt"
	"net"
)

type EchoCommand struct {}

var echoInstance Command

func NewEchoInstance() Command {
	if echoInstance == nil {
		echoInstance = &EchoCommand{}
	}
	return echoInstance
}

func (this *EchoCommand)Exec(conn net.Conn, args []string) []byte {
	return []byte(fmt.Sprint("$", len(args), CRLF, args[0], CRLF))
}
