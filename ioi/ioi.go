package ioi

import (
	"github.com/vmxy/go-ioi/ioi/stream"
)

type Network = string

var (
	Network_Tcp  Network = "tcp"
	Network_Quic         = "quic"
)

var tcp = stream.NewTcp()

func Listen(host string, port int, handle stream.SessionHandle) {
	tcp.Listen(host, port, func(session *stream.Session) {})
}

func Dail(host string, port int) (sess *stream.Session, err error) {
	sess, err = tcp.Connect(host, port)
	if err != nil {
		return
	}
	return
}
