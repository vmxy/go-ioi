package stream

import (
	"context"
	"io"
	"net"

	"github.com/hashicorp/yamux"
	"github.com/quic-go/quic-go"
)

type Session struct {
	id            string
	clientSession any
	serverSession any
	chanClose     chan bool
}

func NewSession[T any](id string, client T, server T) Session {
	session := Session{
		id:            id,
		clientSession: client,
		serverSession: server,
		chanClose:     make(chan bool, 1),
	}
	return session
}

func (s *Session) OpenStream() (conn net.Conn, err error) {
	if sess, ok := any(s.clientSession).(*yamux.Session); ok {
		conn, err = sess.Open()
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
func (s *Session) Close() {
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
func (s *Session) Listen() (net.Listener, error) {
	lis := NewListener(&net.TCPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 0,
	})
	go func() {
		for {
			if sess, ok := any(s.serverSession).(*yamux.Session); ok {
				conn, err := sess.AcceptStream()
				if err != nil {
					if err == io.EOF {
						continue
					}
					s.Close()
					break
				}
				lis.conn <- conn
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
				lis.conn <- conn
			}
		}
	}()
	return &lis, nil
}
