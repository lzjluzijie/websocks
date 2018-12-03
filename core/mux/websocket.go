package mux

import (
	"io"
	"sync"

	"github.com/lzjluzijie/websocks/core"
)

type MuxWebSocket struct {
	*core.WebSocket

	group *Group

	sMutex sync.Mutex
	rMutex sync.Mutex
}

func NewMuxWebSocket(ws *core.WebSocket) (muxWS *MuxWebSocket) {
	muxWS = &MuxWebSocket{
		WebSocket: ws,
	}
	return
}

func (muxWS *MuxWebSocket) Send(m *Message) (err error) {
	muxWS.sMutex.Lock()
	_, err = io.Copy(muxWS, m)
	if err != nil {
		//muxWS.Close()
		return
	}
	muxWS.sMutex.Unlock()

	//debug log
	//log.Printf("sent %#v", m)
	return
}

func (muxWS *MuxWebSocket) Receive() (m *Message, err error) {
	muxWS.rMutex.Lock()
	h := make([]byte, 13)
	_, err = muxWS.Read(h)
	if err != nil {
		//muxWS.Close()
		return
	}

	m = LoadMessage(h)
	data := make([]byte, m.Length)

	_, err = muxWS.Read(data)
	if err != nil {
		//muxWS.Close()
		return
	}
	muxWS.rMutex.Unlock()

	m.Data = data

	////debug log
	//log.Printf("received %#v", m)
	return
}
