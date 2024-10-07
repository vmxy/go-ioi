package stream

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"net"
	"os"
	"time"

	"github.com/quic-go/quic-go"
)

type Quic struct {
	Manager *Manager[*quic.Connection]
}

func NewQuic() *Quic {
	m := NewManager[*quic.Connection]()
	q := &Quic{
		Manager: &m,
	}
	return q
}

type TStream struct {
	stream     *quic.Stream
	localAddr  net.Addr
	remoteAddr net.Addr
}

// Read reads data from the connection.
// Read can be made to time out and return an error after a fixed
// time limit; see SetDeadline and SetReadDeadline.
func (s *TStream) Read(b []byte) (n int, err error) {
	return (*s.stream).Read(b)
}
func (s *TStream) Write(b []byte) (n int, err error) {
	return (*s.stream).Write(b)
}
func (s *TStream) Close() error {
	return (*s.stream).Close()
}
func (s *TStream) LocalAddr() net.Addr {
	return s.localAddr
}
func (s *TStream) RemoteAddr() net.Addr {
	return s.remoteAddr
}
func (s *TStream) SetDeadline(t time.Time) error {
	return (*s.stream).SetDeadline(t)
}
func (s *TStream) SetReadDeadline(t time.Time) error {
	return (*s.stream).SetReadDeadline(t)
}
func (s *TStream) SetWriteDeadline(t time.Time) error {
	return (*s.stream).SetWriteDeadline(t)
}

var _ = (net.Conn)((*TStream)(nil))
var _ = (Accept)((*Quic)(nil))
var _ = (Connect)((*Quic)(nil))

//var _ = (Connect)((*AcceptQuic)(nil))

func (accept *Quic) Listen(host string, port int, handle SessionHandle) {
	hp := fmt.Sprintf("%s:%d", host, port)
	listener, err := quic.ListenAddr(hp, generateTLSConfig(), nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("quic listening on ", hp)

	for {
		sess, err := listener.Accept(context.Background())
		if err != nil {
			log.Println(err)
			continue
		}
		go accept.handleSession(sess, handle)
	}
}

func (accept *Quic) Connect(host string, port int) (net.Conn, error) {
	hp := fmt.Sprintf("%s:%d", host, port)
	var session *quic.Connection
	if sess, ok := accept.Manager.Get(hp); ok {
		session = sess
	} else {
		sess, err := quic.DialAddr(context.Background(), hp, generateTLSConfig(), nil)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		session = &sess
		accept.Manager.Set(hp, session)
	}
	stream, err := (*session).OpenStream()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	conn := &TStream{
		stream:     &stream,
		localAddr:  (*session).LocalAddr(),
		remoteAddr: (*session).RemoteAddr(),
	}
	err = connect(conn, "")
	if err != nil {
		conn.Close()
		return nil, err
	}
	return conn, nil
}

func (accept *Quic) handleSession(sess quic.Connection, handle SessionHandle) {
	defer sess.Context().Done()
	var bs []byte = make([]byte, 0, 2048)
	for {
		stream, err := sess.AcceptStream(context.Background())
		if err != nil {
			log.Println(err)
			return
		}
		var conn net.Conn = &TStream{
			stream:     &stream,
			localAddr:  sess.LocalAddr(),
			remoteAddr: sess.RemoteAddr(),
		}
		size, err := conn.Read(bs)
		if err != nil {
			log.Println("read err", err)
			continue
		}
		handshake := bs[0:size]
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
		conn.Write([]byte("vmfs/1 200 ok\r\n\r\n"))
		go handle(conn)
	}
}

func generateTLSConfig() *tls.Config {
	// 生成和返回 TLS 配置
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal(err)
	}

	serialNumber, err := rand.Int(rand.Reader, big.NewInt(1<<62))
	if err != nil {
		log.Fatal(err)
	}

	cert := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"My Organization"},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(365 * 24 * time.Hour), // Valid for 1 year

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &priv.PublicKey, priv)
	if err != nil {
		log.Fatal(err)
	}

	// Save the certificate and private key to files
	certOut, err := os.Create("cert.pem")
	if err != nil {
		log.Fatal(err)
	}
	defer certOut.Close()
	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		log.Fatal(err)
	}

	keyOut, err := os.Create("key.pem")
	if err != nil {
		log.Fatal(err)
	}
	defer keyOut.Close()
	if err := pem.Encode(keyOut, &pem.Block{Type: "PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)}); err != nil {
		log.Fatal(err)
	}

	return &tls.Config{
		Certificates:       []tls.Certificate{{Certificate: [][]byte{derBytes}, PrivateKey: priv}},
		NextProtos:         []string{"quic-v1"}, // Specify QUIC protocol version
		InsecureSkipVerify: true,                // 开发阶段可以使用，生产环境请使用有效证书
	}
}
