package cmd

import "net"

type Command interface {
	exec(conn net.Conn)
}
