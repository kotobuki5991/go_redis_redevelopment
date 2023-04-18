package cmd

import (
	"fmt"
	"github.com/kotobuki5991/go_redis_redevelopment/app/consts"
)

type EchoCommand struct {}

var echoInstance Command

func NewEchoInstance() Command {
	if echoInstance == nil {
		echoInstance = &EchoCommand{}
	}
	return echoInstance
}

func (cmd *EchoCommand)Exec(args []string) []byte {
	return []byte(fmt.Sprint("$", len(args[0]), consts.CRLF, args[0], consts.CRLF))
}
