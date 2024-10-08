package stream

import "net"

type SessionHandle func(session *Session)

type Accept interface {
	Connect(host string, port int) (*Session, error)
	Listen(host string, port int, handle SessionHandle)
}

type Listener struct {
	conn chan net.Conn
	addr net.Addr
}

func NewListener(addr net.Addr) Listener {
	l := Listener{
		conn: make(chan net.Conn, 2),
		addr: addr,
	}

	return l
}

var _ = (net.Listener)((*Listener)(nil))

func (l *Listener) Accept() (conn net.Conn, err error) {
	conn = <-l.conn
	return conn, nil
}
func (l *Listener) Close() error {
	return nil
}
func (l *Listener) Addr() net.Addr {
	return l.addr
}
