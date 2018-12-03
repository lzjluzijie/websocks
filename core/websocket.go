package core

import (
	"errors"
	"time"

	"github.com/gorilla/websocket"
)

type WebSocket struct {
	conn *websocket.Conn
	buf  []byte

	//stats
	createdAt time.Time
	closed    bool
	stats     *Stats
}

func NewWebSocket(conn *websocket.Conn, stats *Stats) (ws *WebSocket) {
	ws = &WebSocket{
		conn:      conn,
		createdAt: time.Now(),
		stats:     stats,
	}
	return
}

func (ws *WebSocket) Read(p []byte) (n int, err error) {
	if ws.closed == true {
		return 0, errors.New("websocket closed")
	}

	if len(ws.buf) == 0 {
		//debug log
		//log.Println("empty buf, waiting")
		_, ws.buf, err = ws.conn.ReadMessage()
		if err != nil {
			return
		}
	}

	n = copy(p, ws.buf)
	ws.buf = ws.buf[n:]

	if ws.stats != nil {
		go ws.stats.AddDownloaded(uint64(n))
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

	if ws.stats != nil {
		go ws.stats.AddUploaded(uint64(n))
	}
	return
}

func (ws *WebSocket) Close() (err error) {
	ws.conn.Close()
	ws.closed = true
	return
}
