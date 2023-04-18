package main

import (
	"fmt"
)

type EchoCommand struct {}

var echoInstance Command

func NewEchoInstance() Command {
	if echoInstance == nil {
		echoInstance = &EchoCommand{}
	}
	return echoInstance
}

func (this *EchoCommand)Exec(args []string) []byte {
	return []byte(fmt.Sprint("$", len(args[0]), CRLF, args[0], CRLF))
}
