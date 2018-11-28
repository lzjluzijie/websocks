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
	MuxWS *MuxWebSocket

	mutex sync.Mutex
	buf   []byte
	wait  chan int

	receiveMessageID uint64
	sendMessageID    *uint64
}

//NewMuxConn create new mux connection for client
func NewMuxConn(muxWS *MuxWebSocket) (conn *MuxConn) {
	return &MuxConn{
		ID:            rand.Uint64(),
		MuxWS:         muxWS,
		wait:          make(chan int),
		sendMessageID: new(uint64),
	}
}

func (conn *MuxConn) Write(p []byte) (n int, err error) {
	m := &Message{
		Method:    MessageMethodData,
		ConnID:    conn.ID,
		MessageID: conn.SendMessageID(),
		Data:      p,
	}

	err = conn.MuxWS.SendMessage(m)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func (conn *MuxConn) Read(p []byte) (n int, err error) {
	if len(conn.buf) == 0 {
		//log.Printf("%d buf is 0, waiting", conn.ID)
		<-conn.wait
	}

	conn.mutex.Lock()
	//log.Printf("%d buf: %v", conn.buf)
	n = copy(p, conn.buf)
	conn.buf = conn.buf[n:]
	conn.mutex.Unlock()
	return
}

func (conn *MuxConn) HandleMessage(m *Message) (err error) {
	//log.Printf("handle message %d %d", m.ConnID, m.MessageID)
	for {
		if conn.receiveMessageID == m.MessageID {
			conn.mutex.Lock()
			conn.buf = append(conn.buf, m.Data...)
			conn.receiveMessageID++
			close(conn.wait)
			conn.wait = make(chan int)
			conn.mutex.Unlock()
			//log.Printf("handled message %d %d", m.ConnID, m.MessageID)
			return
		}
		<-conn.wait
	}
	return
}

func (conn *MuxConn) SendMessageID() (id uint64) {
	id = atomic.LoadUint64(conn.sendMessageID)
	atomic.AddUint64(conn.sendMessageID, 1)
	return
}

func (conn *MuxConn) Run(c *net.TCPConn) {
	go func() {
		_, err := io.Copy(c, conn)
		if err != nil {
			//log.Printf(err.Error())
		}
	}()

	_, err := io.Copy(conn, c)
	if err != nil {
		//log.Printf(err.Error())
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

	//log.Printf("dial for %s", host)

	err = conn.MuxWS.SendMessage(m)
	if err != nil {
		return
	}

	//log.Printf("%d %s", conn.ID, host)
	return
}
