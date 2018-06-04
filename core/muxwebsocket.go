package core

import (
	"encoding/gob"
	"fmt"
	"net"
	"sync"

	"github.com/pkg/errors"
)

type MuxWebSocket struct {
	*WebSocket
	Decoder *gob.Decoder
	Encoder *gob.Encoder

	connMap sync.Map

	mutex sync.RWMutex
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

func (muxWS *MuxWebSocket) AcceptMuxConn(m *Message) (conn *MuxConn, host string, err error) {
	if m.Method != MessageMethodDial {
		err = errors.New(fmt.Sprintf("wrong message method %d", m.Method))
		return
	}

	host = string(m.Data)

	conn = &MuxConn{
		ID:    m.ConnID,
		muxWS: muxWS,
	}
	muxWS.PutMuxConn(conn)
	return
}

func (muxWS *MuxWebSocket) SendMessage(m *Message) (err error) {
	muxWS.mutex.Lock()
	err = muxWS.Encoder.Encode(m)
	muxWS.mutex.Unlock()
	return
}

func (muxWS *MuxWebSocket) ReceiveMessage() (m *Message, err error) {
	m = &Message{}
	muxWS.mutex.RLock()
	err = muxWS.Decoder.Decode(m)
	muxWS.mutex.RUnlock()
	return
}

func (muxWS *MuxWebSocket) Listen() (err error) {
	//block and listen
	for {
		m, err := muxWS.ReceiveMessage()
		if err != nil {
			return err
		}

		fmt.Println(m)

		//accept new conn
		if m.Method == MessageMethodDial {
			conn, host, err := muxWS.AcceptMuxConn(m)
			if err != nil {
				logger.Debugf(err.Error())
				continue
			}

			logger.Debugf("Accepted mux conn %s", host)

			tcpAddr, err := net.ResolveTCPAddr("tcp", host)
			if err != nil {
				logger.Debugf(err.Error())
				continue
			}

			tcpConn, err := net.DialTCP("tcp", nil, tcpAddr)
			if err != nil {
				logger.Debugf(err.Error())
				continue
			}

			conn.Run(tcpConn)
			continue
		}

		//get conn and send message
		conn := muxWS.GetMuxConn(m.ConnID)
		err = conn.ReceiveMessage(m)
		if err != nil {
			logger.Debugf(err.Error())
			continue
		}
	}
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
