package core

import (
	"io"
	"math/rand"
	"net"
	"sync"
	"sync/atomic"
)

const (
	MessageMethodData = iota
	MessageMethodDial
)

type Message struct {
	Method    byte
	ConnID    uint64
	MessageID uint64
	Data      []byte
}

type MuxConn struct {
	ID    uint64
	muxWS *MuxWebSocket

	messages []*Message
	mutex    sync.Mutex
	buf      []byte

	receiveMessageID uint64
	sendMessageID    *uint64
}

//client use
func NewMuxConn(muxWS *MuxWebSocket) (conn *MuxConn) {
	conn = new(MuxConn)
	conn.muxWS = muxWS
	conn.ID = rand.Uint64()
	return
}

func (conn *MuxConn) Write(p []byte) (n int, err error) {
	m := &Message{
		Method:    MessageMethodData,
		ConnID:    conn.ID,
		MessageID: conn.SendMessageID(),
		Data:      p,
	}

	err = conn.muxWS.SendMessage(m)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func (conn *MuxConn) Read(p []byte) (n int, err error) {
	for {
		if len(conn.buf) != 0 {
			break
		}
	}
	println("readed")

	conn.mutex.Lock()
	n = copy(p, conn.buf)
	conn.buf = conn.buf[n:]
	conn.mutex.Unlock()
	return
}

func (conn *MuxConn) ReceiveMessage(m *Message) (err error) {
	for {
		if conn.receiveMessageID == m.MessageID {
			conn.mutex.Lock()
			conn.buf = append(conn.buf, m.Data...)
			conn.receiveMessageID++
			conn.mutex.Unlock()
			return
		}
	}
	return
}

//client dial remote
func (conn *MuxConn) DialMessage(host string) (err error) {
	m := &Message{
		Method:    MessageMethodDial,
		MessageID: 18446744073709551615,
		ConnID:    conn.ID,
		Data:      []byte(host),
	}

	err = conn.muxWS.SendMessage(m)
	return
}

func (conn *MuxConn) SendMessageID() (id uint64) {
	id = atomic.LoadUint64(conn.sendMessageID)
	atomic.AddUint64(conn.sendMessageID, 1)
	return
}

func (conn *MuxConn) Run(c *net.TCPConn) {
	go func() {
		_, err := io.Copy(conn, c)
		if err != nil {
			logger.Debugf(err.Error())
		}
	}()

	go func() {
		_, err := io.Copy(c, conn)
		if err != nil {
			logger.Debugf(err.Error())
		}
	}()
	return
}
