package core

import (
	"encoding/gob"
	"io"
	"net"
	"net/http"

	"github.com/juju/loggo"
	"golang.org/x/net/websocket"
)

type Server struct {
	LogLevel   loggo.Level
	ListenAddr string
}

func handler(ws *websocket.Conn) {
	var err error
	defer logger.Debugf(err.Error())
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

	http.Handle("/ws", websocket.Handler(handler))
	err = http.ListenAndServe(server.ListenAddr, nil)
	if err != nil {
		return err
	}
	return
}
