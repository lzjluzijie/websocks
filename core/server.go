package core

import (
	"encoding/gob"
	"io"
	"net"
	"net/http"

	"time"

	"github.com/juju/loggo"
	"golang.org/x/net/websocket"
)

type Server struct {
	LogLevel   loggo.Level
	Pattern    string
	ListenAddr string
}

var opened = 0
var closed = 0
var t time.Time

func handler(ws *websocket.Conn) {
	opened++
	var err error
	defer func() {
		closed++
		if err != nil {
			logger.Debugf(err.Error())
		}
	}()

	defer ws.Close()

	dec := gob.NewDecoder(ws)
	req := Request{}
	err = dec.Decode(&req)
	if err != nil {
		return
	}

	logger.Debugf("Dial %s from %s", req.Addr, ws.RemoteAddr().String())
	conn, err := net.Dial("tcp", req.Addr)
	if err != nil {
		return
	}

	defer conn.Close()

	go func() {
		_, err = io.Copy(conn, ws)
		if err != nil {
			logger.Debugf(err.Error())
			return
		}
	}()

	_, err = io.Copy(ws, conn)
	if err != nil {
		return
	}
}

func (server *Server) Listen() (err error) {
	logger.SetLogLevel(server.LogLevel)

	t = time.Now()

	go func() {
		for {
			time.Sleep(time.Second)
			logger.Debugf("%s: opened%d, closed%d", time.Since(t), opened, closed)
		}
	}()

	http.Handle(server.Pattern, websocket.Handler(handler))
	err = http.ListenAndServe(server.ListenAddr, nil)
	if err != nil {
		return err
	}
	return
}
