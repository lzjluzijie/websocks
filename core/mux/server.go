package mux

import (
	"fmt"
	"log"
	"net"

	"github.com/pkg/errors"
)

//ServerAcceptDial
func (group *Group) ServerAcceptDial(m *Message) (err error) {
	if m.Method != MessageMethodDial {
		err = errors.New(fmt.Sprintf("wrong method: %d", m.Method))
		return
	}

	host := string(m.Data)

	conn := &Conn{
		ID:    m.ConnID,
		wait:  make(chan int),
		group: group,

		sendMessageNext:    1,
		receiveMessageNext: 1,
	}

	group.AddConn(conn)

	tcpAddr, err := net.ResolveTCPAddr("tcp", host)
	if err != nil {
		conn.Close()
		log.Printf(err.Error())
		return err
	}

	tcpConn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		conn.Close()
		log.Printf(err.Error())
		return err
	}

	//debug log
	log.Printf("Accepted mux conn: %x, %s", conn.ID, host)

	conn.Run(tcpConn)
	return err
	return
}
