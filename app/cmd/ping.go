package cmd

import (
	"fmt"
	"github.com/kotobuki5991/go_redis_redevelopment/app/consts"
)

type PingCommand struct {}

var pingInstance Command

func NewPingInstance() Command {
	if pingInstance == nil {
		pingInstance = &PingCommand{}
	}
	return pingInstance
}

func (cmd *PingCommand)Exec(args []string) []byte {
	return []byte(fmt.Sprint("+PONG", consts.CRLF))
}
