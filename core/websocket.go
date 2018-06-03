package core

import (
	"sync"

	"github.com/gorilla/websocket"
)

type WebSocket struct {
	conn  *websocket.Conn
	buf   []byte
	mutex sync.Mutex
}

func (ws *WebSocket) Read(p []byte) (n int, err error) {
	if len(ws.buf) == 0 {
		_, ws.buf, err = ws.conn.ReadMessage()
		if err != nil {
			return
		}
	}

	n = copy(p, ws.buf)
	ws.buf = ws.buf[n:]

	return
}

func (ws *WebSocket) Write(p []byte) (n int, err error) {
	ws.mutex.Lock()
	err = ws.conn.WriteMessage(websocket.BinaryMessage, p)
	if err != nil {
		return
	}

	ws.mutex.Unlock()
	return len(p), nil
}

func (ws *WebSocket) Close() (err error) {
	ws.conn.Close()
	return
}
