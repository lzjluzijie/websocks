package mux

import (
	"bytes"
	"io"
	"log"
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
		e := muxWS.Close()
		if e != nil {
			log.Println(e.Error())
		}
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
		e := muxWS.Close()
		if e != nil {
			log.Println(e.Error())
		}
		return
	}

	//debug log
	//log.Printf("%d %x",n, h)

	m = LoadMessage(h)
	buf := &bytes.Buffer{}
	r := io.LimitReader(muxWS, int64(m.Length))

	_, err = io.Copy(buf, r)
	if err != nil {
		e := muxWS.Close()
		if e != nil {
			log.Println(e.Error())
		}
		return
	}
	muxWS.rMutex.Unlock()

	m.Data = buf.Bytes()

	////debug log
	//log.Printf("received %#v", m)
	return
}

func (muxWS *MuxWebSocket) Close() (err error) {
	muxWS.group.MuxWSs = nil
	err = muxWS.WebSocket.Close()
	return
}
