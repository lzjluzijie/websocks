package core

import (
	"encoding/gob"
	"golang.org/x/net/websocket"
	"io"
	"log"
	"net"
	"net/http"
)

type Server struct {
	ListenAddr string
}

func handler(ws *websocket.Conn) {
	dec := gob.NewDecoder(ws)
	s := Socks5{}
	err := dec.Decode(&s)
	if err != nil {
		log.Println(err.Error())
	}

	log.Printf("Dial %s from %s",s.Addr, ws.RemoteAddr().String())
	conn, err := net.Dial("tcp", s.Addr)
	if err != nil {
		log.Println(err.Error())
		return
	}

	go func() {
		_, err = io.Copy(conn, ws)
		if err != nil {
			log.Println(err)
		}
	}()

	_, err = io.Copy(ws, conn)
	if err != nil {
		log.Println(err)
	}
}

func (server *Server)Listen()(err error) {
	http.Handle("/ws", websocket.Handler(handler))
	err = http.ListenAndServe(server.ListenAddr, nil)
	if err != nil {
		return err
	}
	return
}
