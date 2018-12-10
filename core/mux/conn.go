package mux

import (
	"errors"
	"io"
	"log"
	"net"
	"sync"
)

var ErrConnClosed = errors.New("mux conn closed")

type Conn struct {
	closed bool

	ID uint32

	group *Group

	buf      []byte
	bufMutex sync.Mutex
	wait     chan int

	receiveMessageNext uint32
	sendMessageNext    uint32
}

func (conn *Conn) Write(p []byte) (n int, err error) {
	if conn.closed {
		return 0, ErrConnClosed
	}

	m := &Message{
		Method:    MessageMethodData,
		ConnID:    conn.ID,
		MessageID: conn.sendMessageNext,
		Length:    uint32(len(p)),
		Data:      p,
	}
	log.Printf("%d: %d", conn.ID, conn.sendMessageNext)
	conn.sendMessageNext++

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

	conn.bufMutex.Lock()
	//log.Printf("%d buf: %v",conn.ID, conn.buf)
	n = copy(p, conn.buf)
	conn.buf = conn.buf[n:]
	conn.bufMutex.Unlock()
	return
}

func (conn *Conn) HandleMessage(m *Message) (err error) {
	//debug log
	//log.Printf("handle message %d %d", m.ConnID, m.MessageID)

	log.Printf("%d: %d %d", conn.ID, m.MessageID, conn.receiveMessageNext)

	for {
		if conn.closed {
			return ErrConnClosed
		}

		if conn.receiveMessageNext == m.MessageID {
			conn.bufMutex.Lock()
			conn.buf = append(conn.buf, m.Data...)
			conn.receiveMessageNext++
			conn.bufMutex.Unlock()
			close(conn.wait)
			conn.wait = make(chan int)
			//debug log
			//log.Printf("handled message %d %d", m.ConnID, m.MessageID)
			return
		}
		<-conn.wait
	}
	return
}

func (conn *Conn) Run(c *net.TCPConn) {
	go func() {
		_, err := io.Copy(c, conn)
		if err != nil {
			conn.Close()
			c.Close()
			log.Printf(err.Error())
		}
	}()

	_, err := io.Copy(conn, c)
	if err != nil {
		conn.Close()
		c.Close()
		log.Printf(err.Error())
	}

	return
}

func (conn *Conn) Close() (err error) {
	conn.closed = true
	//close(conn.wait)

	go func() {
		m := &Message{
			Method: MessageMethodClose,
			ConnID: conn.ID,
		}

		err = conn.group.Send(m)
		if err != nil {
			log.Println(err.Error())
		}
	}()

	conn.group.DeleteConn(conn.ID)
	return
}
