package stream

import (
	"bytes"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	// 允许跨域访问
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func UpdateWebSocket(header []byte, conn *net.Conn) (net.Conn, error) {
	headers := ParseHttpHeaders(string(header))
	var response http.ResponseWriter = NewHttpResponseWriter(conn, headers)
	resHeaders := make(http.Header)
	reader := bytes.NewReader(header)
	request, _ := http.NewRequest("GET", headers.Get("url"), reader)
	wsConn, err := upgrader.Upgrade(response, request, resHeaders)
	//wsocket := stream.NewWSocket(wsConn)
	if err != nil {
		return nil, err
	}
	nconn := WSConn{stream: wsConn}
	return &nconn, nil
}

type HttpResponseWriter struct {
	Headers    http.Header
	conn       *net.Conn
	StatusCode int
}

func NewHttpResponseWriter(conn *net.Conn, headers http.Header) *HttpResponseWriter {
	return &HttpResponseWriter{
		Headers:    headers,
		StatusCode: http.StatusOK, // 默认状态码为200 OK
		conn:       conn,
	}
}

func (m *HttpResponseWriter) Header() http.Header {
	return m.Headers
}

func (m *HttpResponseWriter) Write(b []byte) (int, error) {
	return (*m.conn).Write(b)
}

func (m *HttpResponseWriter) WriteHeader(statusCode int) {
	m.StatusCode = statusCode
}

type WSConn struct {
	stream *websocket.Conn
}

var _ = (net.Conn)((*WSConn)(nil))

// Read reads data from the connection.
// Read can be made to time out and return an error after a fixed
// time limit; see SetDeadline and SetReadDeadline.
func (s *WSConn) Read(b []byte) (n int, err error) {
	dataType, bs, err := (*s.stream).ReadMessage()
	clear(b)
	if err != nil {
		return 0, err
	}
	if dataType == websocket.TextMessage || dataType == websocket.BinaryMessage {
		n = copy(b, bs)
		return n, err
	}
	return 0, err
}
func (s *WSConn) Write(b []byte) (n int, err error) {
	err = (*s.stream).WriteMessage(websocket.BinaryMessage, b)
	if err != nil {
		return 0, err
	}
	return len(b), nil
}
func (s *WSConn) Close() error {
	return (*s.stream).Close()
}
func (s *WSConn) LocalAddr() net.Addr {
	return (*s.stream).LocalAddr()
}
func (s *WSConn) RemoteAddr() net.Addr {
	return (*s.stream).RemoteAddr()
}
func (s *WSConn) SetDeadline(t time.Time) error {
	return nil
}
func (s *WSConn) SetReadDeadline(t time.Time) error {
	return (*s.stream).SetReadDeadline(t)
}
func (s *WSConn) SetWriteDeadline(t time.Time) error {
	return (*s.stream).SetWriteDeadline(t)
}
