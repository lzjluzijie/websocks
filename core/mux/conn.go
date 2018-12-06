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

	buf           []byte
	bufMutex      sync.Mutex
	wait          chan int
	messages      []*Message
	messagesMutex sync.Mutex

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
	log.Println(conn.sendMessageNext)
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
	log.Printf("%d %d", m.MessageID, conn.receiveMessageNext)
	if conn.closed {
		return ErrConnClosed
	}

	//debug log
	//log.Printf("handle message %d %d", m.ConnID, m.MessageID)

	if m.MessageID < conn.receiveMessageNext {
		err = errors.New("invalid message id")
		return
	}

	if m.MessageID == conn.receiveMessageNext {
		conn.bufMutex.Lock()
		conn.messagesMutex.Lock()
		conn.buf = append(conn.buf, m.Data...)
		conn.receiveMessageNext++
		for i, m := range conn.messages {
			if m == nil {
				conn.messages = conn.messages[i:]
				continue
			}

			conn.buf = append(conn.buf, m.Data...)
			conn.receiveMessageNext++
		}
		conn.messagesMutex.Unlock()
		conn.bufMutex.Unlock()

		close(conn.wait)
		conn.wait = make(chan int)
		return
	}

	conn.messagesMutex.Lock()
	i := m.MessageID - conn.receiveMessageNext
	if i < uint32(len(conn.messages)) {
		conn.messages[i] = m
		conn.messagesMutex.Unlock()
		return
	}

	if i == uint32(len(conn.messages)) {
		conn.messages = append(conn.messages, m)
		conn.messagesMutex.Unlock()
		return
	}

	d := i - uint32(len(conn.messages)) + 1
	log.Printf("%d %d %d %d %d", i, m.MessageID, conn.receiveMessageNext, uint32(len(conn.messages)), d)
	s := make([]*Message, d)
	s[d-1] = m
	conn.messages = append(conn.messages, s...)
	conn.messagesMutex.Unlock()
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
	conn.group.DeleteConn(conn.ID)
	//close(conn.wait)
	conn.closed = true
	return
}
