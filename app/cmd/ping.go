package cmd

import (
	"fmt"
	"net"
	"app/consts"
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
	return []byte(fmt.Sprint("+PONG", consts.CRLF))
}
