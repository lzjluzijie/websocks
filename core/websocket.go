package core

import (
	"errors"
	"time"

	"github.com/gorilla/websocket"
	"github.com/juju/loggo"
)

var logger = loggo.GetLogger("core")

type WebSocket struct {
	conn *websocket.Conn
	buf  []byte

	//stats
	createdAt     time.Time
	closed        bool
	AddDownloaded func(downloaded uint64)
	AddUploaded   func(uploaded uint64)
}

func NewWebSocket(conn *websocket.Conn) (ws *WebSocket) {
	ws = &WebSocket{
		conn:      conn,
		createdAt: time.Now(),
	}
	return
}

func (ws *WebSocket) Read(p []byte) (n int, err error) {
	if ws.closed == true {
		return 0, errors.New("websocket closed")
	}

	if len(ws.buf) == 0 {
		_, ws.buf, err = ws.conn.ReadMessage()
		if err != nil {
			return
		}
	}

	n = copy(p, ws.buf)
	ws.buf = ws.buf[n:]

	if ws.AddDownloaded != nil {
		go ws.AddDownloaded(uint64(n))
	}
	return
}

func (ws *WebSocket) Write(p []byte) (n int, err error) {
	if ws.closed == true {
		return 0, errors.New("websocket closed")
	}

	err = ws.conn.WriteMessage(websocket.BinaryMessage, p)
	if err != nil {
		return
	}

	n = len(p)

	if ws.AddUploaded != nil {
		go ws.AddUploaded(uint64(n))
	}
	return
}

func (ws *WebSocket) Close() (err error) {
	ws.conn.Close()
	ws.closed = true
	return
}
