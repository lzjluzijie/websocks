package mux

import (
	"errors"
	"io"
	"log"
	"net"
	"sync"
	"sync/atomic"
)

var ErrConnClosed = errors.New("mux conn closed")

type Conn struct {
	ID uint32

	group *Group

	mutex sync.Mutex
	buf   []byte
	wait  chan int

	closed bool

	receiveMessageID uint32
	sendMessageID    *uint32
}

func (conn *Conn) Write(p []byte) (n int, err error) {
	if conn.closed {
		return 0, ErrConnClosed
	}

	m := &Message{
		Method:    MessageMethodData,
		ConnID:    conn.ID,
		MessageID: conn.SendMessageID(),
		Length:    uint32(len(p)),
		Data:      p,
	}

	err = conn.group.Send(m)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func (conn *Conn) Read(p []byte) (n int, err error) {
	if conn.closed {
		return 0, ErrConnClosed
	}

	if len(conn.buf) == 0 {
		//log.Printf("%d buf is 0, waiting", conn.ID)
		<-conn.wait
	}

	conn.mutex.Lock()
	//log.Printf("%d buf: %v",conn.ID, conn.buf)
	n = copy(p, conn.buf)
	conn.buf = conn.buf[n:]
	conn.mutex.Unlock()
	return
}

func (conn *Conn) HandleMessage(m *Message) (err error) {
	if conn.closed {
		return ErrConnClosed
	}

	//debug log
	//log.Printf("handle message %d %d", m.ConnID, m.MessageID)

	for {
		if conn.receiveMessageID == m.MessageID {
			conn.mutex.Lock()
			conn.buf = append(conn.buf, m.Data...)
			conn.receiveMessageID++
			close(conn.wait)
			conn.wait = make(chan int)
			conn.mutex.Unlock()
			//debug log
			//log.Printf("handled message %d %d", m.ConnID, m.MessageID)
			return
		}
		<-conn.wait
	}
	return
}

func (conn *Conn) SendMessageID() (id uint32) {
	id = atomic.LoadUint32(conn.sendMessageID)
	atomic.AddUint32(conn.sendMessageID, 1)
	return
}

func (conn *Conn) Run(c *net.TCPConn) {
	go func() {
		_, err := io.Copy(c, conn)
		if err != nil {
			conn.Close()
			log.Printf(err.Error())
		}
	}()

	_, err := io.Copy(conn, c)
	if err != nil {
		conn.Close()
		log.Printf(err.Error())
	}

	return
}

func (conn *Conn) Close() (err error) {
	conn.group.DeleteConn(conn.ID)
	//close(conn.wait)
	conn.closed = true
	return
}
