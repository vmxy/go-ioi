package stream

import (
	"fmt"
	"log"
	"net"

	"github.com/hashicorp/yamux"
)

type Tcp struct {
	Manager  *Manager[*yamux.Session]
	chanConn chan net.Conn
}

var _ = (Accept)((*Tcp)(nil))
var _ = (Connect)((*Tcp)(nil))

//var _ = (Connect)((*AcceptTcp)(nil))

func NewTcp() *Tcp {
	m := NewManager[*yamux.Session]()
	q := &Tcp{
		Manager:  &m,
		chanConn: make(chan net.Conn, 1),
	}
	return q
}
func (accept *Tcp) Listen(host string, port int, handle SessionHandle) {
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
			log.Println(err)
			continue
		}
		size, err := conn.Read(bs)
		if err != nil {
			log.Println("read err", err)
			continue
		}
		handshake := bs[0:size]
		fmt.Println("read====", size, string(handshake))
		if IsWebSocket(handshake) {
			conn, err = UpdateWebSocket(handshake, &conn)
			if err != nil {
				log.Println("websocket accept error", err)
				continue
			}
			size, err = conn.Read(bs)
			if err != nil {
				log.Println("websocket read err", err)
				continue
			}
			handshake = bs[0:size]
		}
		if !IsVMFS(handshake) {
			conn.Close()
			continue
		}
		_, id := parseVMFSRequest(handshake)
		fmt.Println("id==", id)
		conn.Write([]byte("vmfs/1 200 ok\r\n\r\n"))
		// Setup server side of yamux
		sess, err := yamux.Server(conn, nil)
		if err != nil {
			log.Println(err)
			continue
		}
		fmt.Println("xxxxx add mmmm ", sess.RemoteAddr().String())
		go accept.handleSession(sess, handle)
	}
}

/*
	func (accept *Tcp) ConnectVirtual(id string) (net.Conn, error) {

}
*/
func (accept *Tcp) handleSession(sess *yamux.Session, handle SessionHandle) {
	key := sess.RemoteAddr().String()
	fmt.Println("------------>key", key)
	defer func() {
		sess.Close()
		accept.Manager.Delete(key)
	}()
	for {
		stream, err := sess.AcceptStream()
		if err != nil {
			log.Println(err)
			return
		}
		go handle(stream)
	}
}

func (accept *Tcp) Connect(host string, port int) (net.Conn, error) {
	hp := fmt.Sprintf("%s:%d", host, port)
	var session *yamux.Session
	if sess, ok := accept.Manager.Get(hp); ok {
		if sess.IsClosed() {
			accept.Manager.Delete(hp)
		} else {
			session = sess
		}
	}
	if session == nil {
		conn, err := net.Dial("tcp", hp)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		err = connect(conn, "1")
		fmt.Println("ddd", "tcp", hp, err)

		if err != nil {
			conn.Close()
			return nil, err
		}

		sess, err := yamux.Client(conn, nil)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		go accept.handleSession(sess, func(conn net.Conn) {
			fmt.Println("==========>>>>>>>>>>>>>accept----")
			accept.chanConn <- conn
		})
		session = sess
		accept.Manager.Set(hp, session)
	}
	conn, err := session.Open() //OpenStream
	if err != nil {
		log.Println(err)
		accept.Manager.Delete(hp)
		return nil, err
	}
	return conn, nil
}
