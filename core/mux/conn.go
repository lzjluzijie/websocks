package mux

import (
	"bytes"
	"crypto/sha256"
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

	buf         bytes.Buffer
	wait        chan int
	handleMutex sync.Mutex

	receiveMessageNext uint32
	sendMessageNext    uint32
}

func (conn *Conn) Write(p []byte) (n int, err error) {
	if conn.closed {
		return 0, io.EOF
		//return 0, ErrConnClosed
	}

	m := &Message{
		Method:    MessageMethodData,
		ConnID:    conn.ID,
		MessageID: conn.sendMessageNext,
		Length:    uint32(len(p)),
		Data:      p,
	}
	//log.Printf("%d: %d", conn.ID, conn.sendMessageNext)
	conn.sendMessageNext++

	err = conn.group.Send(m)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func (conn *Conn) Read(p []byte) (n int, err error) {
	if conn.closed {
		return 0, io.EOF
		//return 0, ErrConnClosed
	}

	if conn.buf.Len() == 0 {
		//log.Printf("%d buf is 0, waiting", conn.ID)
		<-conn.wait
	}

	n, err = conn.buf.Read(p)
	return
}

func (conn *Conn) HandleMessage(m *Message) {
	//debug log
	//log.Printf("handle message %d %d", m.ConnID, m.MessageID)

	//log.Printf("%d: %d %d", conn.ID, m.MessageID, conn.receiveMessageNext)

	for {
		if conn.closed {
			return
		}

		conn.handleMutex.Lock()
		if conn.receiveMessageNext == m.MessageID {
			conn.buf.Write(m.Data)
			conn.receiveMessageNext = m.MessageID + 1
			close(conn.wait)
			conn.wait = make(chan int)
			conn.handleMutex.Unlock()
			//debug log
			//log.Printf("handled message %d %d", m.ConnID, m.MessageID)
			return
		}
		conn.handleMutex.Unlock()
		<-conn.wait
	}
	return
}

func (conn *Conn) Run(c *net.TCPConn) {
	go func() {
		h := sha256.New()
		r := io.TeeReader(conn, h)
		_, err := io.Copy(c, r)
		conn.Close()
		c.Close()
		log.Printf("%x write: %x", conn.ID, h.Sum(nil))
		if err != nil {
			log.Printf(err.Error())
		}
	}()

	h := sha256.New()
	r := io.TeeReader(c, h)
	_, err := io.Copy(conn, r)
	conn.Close()
	c.Close()
	log.Printf("%x read: %x", conn.ID, h.Sum(nil))
	if err != nil {
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
