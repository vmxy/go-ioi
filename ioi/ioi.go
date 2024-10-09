package ioi

import (
	"github.com/vmxy/go-ioi/ioi/stream"
	"github.com/vmxy/go-ioi/ioi/util"
)

type Network = string
type Session = stream.Session
type SessionHandle = stream.SessionHandle

var Log = util.Log

/*
var (
	Network_Tcp  Network = "tcp"
	Network_Quic         = "quic"
) */

var tcp = stream.NewTcp()

func Listen(host string, port int, handle SessionHandle) {
	tcp.Listen(host, port, func(session *Session) {
		handle(session)
	})
}

func Dail(host string, port int) (sess *Session, err error) {
	sess, err = tcp.Connect(host, port)
	if err != nil {
		return
	}
	return
}
