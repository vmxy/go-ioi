package stream

import "net"

type SessionHandle func(session *Session[any])

type Accept interface {
	Listen(host string, port int, handle SessionHandle)
}

type Connect interface {
	Connect(host string, port int) (net.Conn, error)
}
