package core

import (
	"encoding/gob"
	"sync"
)

type MuxWebSocket struct {
	*WebSocket
	Decoder *gob.Decoder
	Encoder *gob.Encoder

	connMap sync.Map

	//mutex sync.RWMutex
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
	logger.Debugf("sent %#v", m)
	return
}

func (muxWS *MuxWebSocket) ReceiveMessage() (m *Message, err error) {
	m = &Message{}
	err = muxWS.Decoder.Decode(m)
	logger.Debugf("received %#v", m)
	return
}

func (muxWS *MuxWebSocket) PutMuxConn(conn *MuxConn) {
	muxWS.connMap.Store(conn.ID, conn)
	return
}

func (muxWS *MuxWebSocket) GetMuxConn(id uint64) (conn *MuxConn) {
	c, ok := muxWS.connMap.Load(id)
	if !ok {
		panic("not ok!")
	}

	return c.(*MuxConn)
}
