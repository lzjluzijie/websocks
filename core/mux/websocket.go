package mux

import (
	"encoding/binary"
	"log"
	"sync"

	"github.com/lzjluzijie/websocks/core"
)

type MuxWebSocket struct {
	*core.WebSocket

	group *Group

	mutex sync.Mutex
}

func NewMuxWebSocket(ws *core.WebSocket) (muxWS *MuxWebSocket) {
	muxWS = &MuxWebSocket{
		WebSocket: ws,
	}
	return
}

func (muxWS *MuxWebSocket) Send(m *MessageHead, data []byte) (err error) {
	err = binary.Write(muxWS, binary.BigEndian, m)
	if err != nil {
		return
	}

	_, err = muxWS.Write(data)
	if err != nil {
		return
	}

	log.Printf("sent %#v", m)
	return
}

func (muxWS *MuxWebSocket) Receive(m *MessageHead, data []byte) (err error) {
	err = binary.Read(muxWS, binary.BigEndian, m)
	if err != nil {
		return
	}

	_, err = muxWS.Read(data)
	if err != nil {
		return
	}

	log.Printf("received %#v", m)
	return
}

func (muxWS *MuxWebSocket) Listen() {
	go func() {
		for {
			mh := &MessageHead{}
			data := make([]byte, 0)
			err := muxWS.Receive(mh, data)
			if err != nil {
				//todo
				log.Printf(err.Error())
				continue
			}

			muxWS.group.Receive(mh, data)
		}
		return
	}()
}
