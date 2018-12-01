package mux

import (
	"log"
	"net"
)

//ServerHandleMessage is a server group function
func (group *Group) ServerHandleMessage(m *Message) (err error) {
	//accept new conn
	if m.Method == MessageMethodDial {
		host := string(m.Data)
		log.Printf("start to dial %s", host)
		conn := &Conn{
			ID:            m.ConnID,
			wait:          make(chan int),
			sendMessageID: new(uint32),
			group:         group,
		}

		//add to group before receive data
		group.Conns = append(group.Conns, conn)

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
