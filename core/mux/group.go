package mux

import (
	"errors"
	"log"
	"math/rand"
	"net"
)

type Group struct {
	client bool

	MuxWSs []*MuxWebSocket

	Conns []*Conn
}

//
//true: client group
//false: server group
func NewGroup(client bool) (group *Group) {
	group = &Group{
		client: client,
	}

	return
}

func (group *Group) Send(mh *MessageHead, data []byte) (err error) {
	//todo
	err = group.MuxWSs[0].Send(mh, data)
	return
}

func (group *Group) Receive(mh *MessageHead, data []byte) (err error) {
	if !group.client {
		//accept new conn
		if mh.Method == MessageMethodDial {
			host := string(data)
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
	}

	//get conn and send message
	//todo better way to find conn
	for _, conn := range group.Conns {
		if conn.ID == mh.ConnID {
			err = conn.HandleMessage(mh, data)
			if err != nil {
				return
			}
		}
	}

	err = errors.New("conn does not exist")
	return
}

func (group *Group) Start() (err error) {

	return
}

func (group *Group) AddMuxWS(muxWS *MuxWebSocket) (err error) {
	muxWS.group = group
	group.MuxWSs = append(group.MuxWSs, muxWS)
	muxWS.Listen()
	return
}
