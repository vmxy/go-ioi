package stream

import "net"

type SessionHandle func(conn net.Conn)

type Accept interface {
	Listen(host string, port int, handle SessionHandle)
}

type Connect interface {
	Connect(host string, port int) (net.Conn, error)
}
