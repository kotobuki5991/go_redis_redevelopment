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
	fmt.Println("args[0]")
	fmt.Println(args[0])
	fmt.Println("args[0]")
	fmt.Println("args[1]")
	fmt.Println(args[1])
	fmt.Println("args[1]")
	return []byte(fmt.Sprint("$", len(args), CRLF, args[0], CRLF))
}
