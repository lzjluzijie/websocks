package core

import (
	"encoding/gob"
	"io"
	"log"
	"net"

	"golang.org/x/net/websocket"
)

type Local struct {
	ListenAddr *net.TCPAddr
	LocalConn  chan *net.TCPConn
	URL        string
	Origin     string
}

func (local *Local) Listen() error {
	listener, err := net.ListenTCP("tcp", local.ListenAddr)
	if err != nil {
		return err
	}
	log.Printf("Listening at: %s", local.ListenAddr.String())

	defer listener.Close()

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Println(err)
			continue
		}

		go local.handleConn(conn)
	}
	return nil
}

func (local *Local) handleConn(conn *net.TCPConn) (err error) {
	defer conn.Close()
	conn.SetLinger(0)

	err = handShake(conn)
	if err != nil {
		log.Println(err)
		return
	}

	_, host, err := getRequest(conn)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(host)

	_, err = conn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x08, 0x43})
	if err != nil {
		log.Println(err)
		return
	}

	ws, err := websocket.Dial(local.URL, "", local.Origin)
	if err != nil {
		log.Println(err)
		return
	}

	defer ws.Close()

	enc := gob.NewEncoder(ws)
	req := Request{
		Addr: host,
	}
	err = enc.Encode(req)
	if err != nil {
		log.Println(err)
		return
	}

	go func() {
		_, err = io.Copy(ws, conn)
		if err != nil {
			log.Println(err)
			return
		}
	}()

	_, err = io.Copy(conn, ws)
	if err != nil {
		log.Println(err)
		return
	}
	return
}
