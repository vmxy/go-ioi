package stream

import (
	"net"
	"strings"

	"github.com/hashicorp/yamux"
	"github.com/vmxy/go-ioi/ioi/util"
)

var mtcp = NewTcp()
var mquic = NewQuic()

//var ManagerAccept = NewManager[net.Conn]()

func Listen(host string, port int, handle SessionHandle) {
	go mquic.Listen(host, port, func(conn net.Conn) {
		go handle(conn)

	})

	go mtcp.Listen(host, port, func(conn net.Conn) {
		go handle(conn)

	})
	//return &l
}

func Dial(host string, port int) (net.Conn, error) {
	return mtcp.Connect(host, port)
	/*
		 	if network == "quic" {
				return mquic.Connect(host, port)
			} else {
				return mtcp.Connect(host, port)
			}
	*/
}

type Stream struct {
	server        *string
	signConn      *net.Conn
	clientSession *yamux.Session
	serverSession *yamux.Session
	chanConn      chan net.Conn
}

func NewStream(server *string) *Stream {
	s := &Stream{
		server:   server,
		chanConn: make(chan net.Conn, 1),
	}
	return s
}
func (s *Stream) Listen(host string, port int, handle SessionHandle) {
	go mquic.Listen(host, port, func(conn net.Conn) {
		//go s.handleReceive(conn, handle)
	})
	go mtcp.Listen(host, port, func(conn net.Conn) {
		//go s.handleReceive(conn, handle)
	})
}
func (s *Stream) Server() (host string, port int) {
	if s.server == nil || len(*s.server) < 1 {
		return "", 0
	}
	sp := strings.Split(*s.server, ":")
	if len(sp) < 2 {
		return "", 0
	}
	host = sp[0]
	port = util.ParseInt(sp[1])
	return
}
func (s *Stream) Connect(host string, port int) (conn net.Conn, err error) {

	return
}
