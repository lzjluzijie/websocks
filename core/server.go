package core

import (
	"encoding/gob"
	"io"
	"log"
	"net"
	"net/http"

	"golang.org/x/net/websocket"
)

type Server struct {
	ListenAddr string
}

func handler(ws *websocket.Conn) {
	defer ws.Close()

	dec := gob.NewDecoder(ws)
	req := Request{}
	err := dec.Decode(&req)
	if err != nil {
		log.Println(err.Error())
		return
	}

	log.Printf("Dial %s from %s", req.Addr, ws.RemoteAddr().String())
	conn, err := net.Dial("tcp", req.Addr)
	if err != nil {
		log.Println(err.Error())
		return
	}

	defer conn.Close()

	go func() {
		_, err = io.Copy(conn, ws)
		if err != nil {
			log.Println(err)
			return
		}
	}()

	_, err = io.Copy(ws, conn)
	if err != nil {
		log.Println(err)
		return
	}
}

func (server *Server) Listen() (err error) {
	http.Handle("/ws", websocket.Handler(handler))
	err = http.ListenAndServe(server.ListenAddr, nil)
	if err != nil {
		return err
	}
	return
}
