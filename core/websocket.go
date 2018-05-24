package core

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:   4 * 1024,
	WriteBufferSize:  4 * 1024,
	HandshakeTimeout: 10 * time.Second,
}

type WebSocket struct {
	conn   *websocket.Conn
	buf []byte
}

func NewWebSocket(w http.ResponseWriter, r *http.Request) (ws *WebSocket, err error) {
	c, err := upgrader.Upgrade(w, r, nil)
	ws = &WebSocket{
		conn: c,
	}
	return
}

func (ws *WebSocket) Read(p []byte) (n int, err error) {
	if len(ws.buf) == 0 {
		_, ws.buf, err = ws.conn.ReadMessage()
		if err != nil{
			return
		}
	}

	n = copy(p, ws.buf)
	ws.buf = ws.buf[n:]
	return
}

func (ws *WebSocket) Write(p []byte) (n int, err error) {
	err = ws.conn.WriteMessage(websocket.BinaryMessage, p)
	if err != nil {
		return
	}

	return len(p), nil
}
