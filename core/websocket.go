package core

import (
	"github.com/gorilla/websocket"
	"github.com/juju/loggo"
)

var logger = loggo.GetLogger("core")

type WebSocket struct {
	conn *websocket.Conn
	buf  []byte
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
	err = ws.conn.WriteMessage(websocket.BinaryMessage, p)
	if err != nil {
		return
	}

	return len(p), nil
}

func (ws *WebSocket) Close() (err error) {
	ws.conn.Close()
	return
}
