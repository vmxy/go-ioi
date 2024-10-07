package stream

import (
	"context"
	"net"

	"github.com/hashicorp/yamux"
	"github.com/quic-go/quic-go"
)

type Session[T any] struct {
	clientSession T
	serverSession T
	chanClose     chan bool
}

func NewSession[T any](client T, server T) Session[T] {
	session := Session[T]{
		clientSession: client,
		serverSession: server,
		chanClose:     make(chan bool, 1),
	}
	return session
}

func (s *Session[T]) OpenStream() (conn net.Conn, err error) {
	if sess, ok := any(s.clientSession).(*yamux.Session); ok {
		conn, err = sess.OpenStream()
		return
	}
	if sess, ok := any(s.clientSession).(*quic.Connection); ok {
		quicCon, e1 := (*sess).OpenStream()
		if e1 != nil {
			return nil, e1
		}
		conn = &TStream{
			stream:     &quicCon,
			localAddr:  (*sess).LocalAddr(),
			remoteAddr: (*sess).RemoteAddr(),
		}
		return
	}
	return
}
func (s *Session[T]) Close() {
	defer func() {
		s.chanClose <- true
	}()
	if sess, ok := any(s.serverSession).(*yamux.Session); ok {
		if !sess.IsClosed() {
			sess.Close()
		}
	}
	if sess, ok := any(s.clientSession).(*yamux.Session); ok {
		if !sess.IsClosed() {
			sess.Close()
		}
	}
}
func (s *Session[T]) Accept(handle SessionHandle) {
	go func() {
		for {
			if sess, ok := any(s.serverSession).(*yamux.Session); ok {
				conn, err := sess.AcceptStream()
				if err != nil {
					s.Close()
					return
				}
				handle(conn)
			} else if sess, ok := any(s.serverSession).(*quic.Connection); ok {
				quicCon, err := (*sess).AcceptStream(context.Background())
				if err != nil {
					s.Close()
				}
				conn := &TStream{
					stream:     &quicCon,
					localAddr:  (*sess).LocalAddr(),
					remoteAddr: (*sess).RemoteAddr(),
				}
				handle(conn)
			}
		}
	}()
}
