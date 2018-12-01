package mux

import (
	"io"
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

func (muxWS *MuxWebSocket) Send(m *Message) (err error) {
	_, err = io.Copy(muxWS, m)
	if err != nil {
		return
	}

	log.Printf("sent %#v", m)
	return
}

func (muxWS *MuxWebSocket) Receive(m *Message) (err error) {
	_, err = io.Copy(m, muxWS)
	if err != nil {
		return
	}

	log.Printf("received %#v", m)
	return
}

func (muxWS *MuxWebSocket) Listen() {
	go func() {
		for {
			m := &Message{}
			err := muxWS.Receive(m)
			if err != nil {
				//todo
				log.Printf(err.Error())
				continue
			}

			muxWS.group.Receive(m)
		}
		return
	}()
}
