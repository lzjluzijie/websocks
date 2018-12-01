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
	muxWS.mutex.Lock()
	_, err = io.Copy(muxWS, m)
	if err != nil {
		return
	}

	log.Printf("sent %#v", m)
	muxWS.mutex.Unlock()
	return
}

func (muxWS *MuxWebSocket) Receive(m *Message) (err error) {
	h := make([]byte, 13)
	_, err = muxWS.Read(h)
	if err != nil {
		return
	}

	m.SetHead(h)
	data := make([]byte, m.Length)

	_, err = muxWS.Read(data)
	if err != nil {
		return
	}

	m.Data = data
	log.Printf("received %#v", m)
	return
}

func (muxWS *MuxWebSocket) Listen() {
	go func() {
		for {
			m := &Message{}
			err := muxWS.Receive(m)
			if err != nil {
				log.Printf(err.Error())
				return
			}

			go muxWS.group.Receive(m)
		}
		return
	}()
}
