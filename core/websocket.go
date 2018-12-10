package core

import (
	"errors"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

var ErrWebSocketClosed = errors.New("websocket closed")

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
		return 0, ErrWebSocketClosed
	}

	if len(ws.buf) == 0 {
		//debug log
		//log.Println("empty buf, waiting")
		_, ws.buf, err = ws.conn.ReadMessage()
		if err != nil {
			e := ws.Close()
			if e != nil {
				log.Println(err.Error())
			}
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
		return 0, ErrWebSocketClosed
	}

	err = ws.conn.WriteMessage(websocket.BinaryMessage, p)
	if err != nil {
		e := ws.Close()
		if e != nil {
			log.Println(err.Error())
		}
		return
	}

	n = len(p)

	if ws.stats != nil {
		go ws.stats.AddUploaded(uint64(n))
	}
	return
}

func (ws *WebSocket) Close() (err error) {
	ws.closed = true
	err = ws.conn.Close()
	return
}
