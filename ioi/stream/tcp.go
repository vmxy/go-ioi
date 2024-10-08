package stream

import (
	"fmt"
	"log"
	"net"

	"github.com/google/uuid"
	"github.com/hashicorp/yamux"
	"github.com/vmxy/go-ioi/ioi/util"
)

type Tcp struct {
	chanConn chan net.Conn
}

var _ = (Accept)((*Tcp)(nil))

//var _ = (Connect)((*Tcp)(nil))

//var _ = (Connect)((*AcceptTcp)(nil))

func NewTcp() *Tcp {
	q := &Tcp{
		chanConn: make(chan net.Conn, 1),
	}
	return q
}
func (accept *Tcp) Listen(host string, port int, handle SessionHandle) {
	maps := util.NewMap[string, *Session]()
	defer maps.Clear()
	hp := fmt.Sprintf("%s:%d", host, port)
	listener, err := net.Listen("tcp", hp)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("tcp listening on ", hp)
	var bs []byte = make([]byte, 2048)
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		size, err := conn.Read(bs)
		if err != nil {
			log.Println("read err", err)
			continue
		}
		chunk := bs[0:size]
		fmt.Println("read====", size, string(chunk))
		if IsWebSocket(chunk) {
			conn, err = UpdateWebSocket(chunk, &conn)
			if err != nil {
				log.Println("websocket accept error", err)
				continue
			}
			size, err = conn.Read(bs)
			if err != nil {
				log.Println("websocket read err", err)
				continue
			}
			chunk = bs[0:size]
		}
		if !IsVMFS(chunk) {
			conn.Close()
			continue
		}
		_, connectType, sid := parseVMFSRequest(chunk)
		fmt.Println("id==", sid)
		conn.Write([]byte("vmfs/1 200 ok\r\n\r\n"))
		sess, find := maps.Get(sid)
		if !find {
			sess1 := NewSession[*yamux.Session](nil, nil)
			sess = &sess1
			maps.Set(sid, &sess1)
		}
		if connectType == "server" {
			// Setup server side of yamux
			client, err := yamux.Client(conn, nil)
			if err != nil {
				log.Println(err)
				continue
			}
			sess.clientSession = client
			if sess.serverSession == nil {
				maps.Delete(sid)
				sess.Close()
				continue
			}
			handle(sess)
		} else {
			// Setup server side of yamux
			server, err := yamux.Server(conn, nil)
			if err != nil {
				log.Println(err)
				continue
			}
			sess.serverSession = server
		}
	}
}

func (accept *Tcp) Connect(host string, port int) (*Session, error) {
	sid := uuid.New().String()

	hostport := fmt.Sprintf("%s:%d", host, port)
	client, err := connectClient(hostport, sid)
	if err != nil {
		return nil, err
	}
	server, err := connectServer(hostport, sid)
	if err != nil {
		return nil, err
	}
	session := NewSession(client, server)
	return &session, nil
}
