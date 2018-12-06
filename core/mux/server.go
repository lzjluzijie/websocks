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

		//debug log
		//log.Printf("start to dial %s", host)

		conn := &Conn{
			ID:    m.ConnID,
			wait:  make(chan int),
			group: group,

			sendMessageNext:    1,
			receiveMessageNext: 1,
		}

		//add to group before receive data
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
	}
	return
}
