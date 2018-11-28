package core

import (
	"encoding/gob"
	"sync"
)

type MuxWebSocket struct {
	*WebSocket
	Decoder *gob.Decoder
	Encoder *gob.Encoder

	muxConns  []*MuxConn
	muxConnID []uint64
	mutex     sync.Mutex
}

func NewMuxWebSocket(ws *WebSocket) (muxWS *MuxWebSocket) {
	dec := gob.NewDecoder(ws)
	enc := gob.NewEncoder(ws)

	muxWS = &MuxWebSocket{
		WebSocket: ws,
		Decoder:   dec,
		Encoder:   enc,
	}
	return
}

func (muxWS *MuxWebSocket) SendMessage(m *Message) (err error) {
	err = muxWS.Encoder.Encode(m)
	//log.Printf("sent %#v", m)
	return
}

func (muxWS *MuxWebSocket) ReceiveMessage() (m *Message, err error) {
	m = &Message{}
	err = muxWS.Decoder.Decode(m)
	//log.Printf("received %#v", m)
	return
}

func (muxWS *MuxWebSocket) PutMuxConn(conn *MuxConn) {
	muxWS.mutex.Lock()
	muxWS.muxConns = append(muxWS.muxConns, conn)
	muxWS.muxConnID = append(muxWS.muxConnID, conn.ID)
	muxWS.mutex.Unlock()
	return
}

func (muxWS *MuxWebSocket) GetMuxConn(connID uint64) (conn *MuxConn) {
	muxWS.mutex.Lock()
	for n, id := range muxWS.muxConnID {
		if id == connID {
			conn = muxWS.muxConns[n]
			break
		}
	}
	muxWS.mutex.Unlock()
	return
}
