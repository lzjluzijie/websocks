package mux

import (
	"log"
	"math/rand"
	"net"
)

//HandleMessage is a server group function
func (group *Group) HandleMessage(m *Message) (err error) {
	//accept new conn
	if m.Method == MessageMethodDial {
		host := string(m.Data)
		log.Printf("start to dial %s", host)
		conn := &Conn{
			ID:            rand.Uint32(),
			wait:          make(chan int),
			sendMessageID: new(uint32),
		}

		tcpAddr, err := net.ResolveTCPAddr("tcp", host)
		if err != nil {
			log.Printf(err.Error())
			return err
		}

		tcpConn, err := net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			log.Printf(err.Error())
			return err
		}

		log.Printf("Accepted mux conn %s", host)

		conn.Run(tcpConn)
		return err
	}
	return
}
