package cmd

import "net"

type Command interface {
	Exec(conn net.Conn, args []string) []byte
}
