package main

import "net"

type Command interface {
	Exec(conn net.Conn, args []string) []byte
}
