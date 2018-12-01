package mux

import (
	"errors"
	"log"
	"time"
)

type Group struct {
	client bool

	MuxWSs []*MuxWebSocket

	Conns []*Conn
}

//true: client group
//false: server group
func NewGroup(client bool) (group *Group) {
	group = &Group{
		client: client,
	}
	return
}

func (group *Group) Send(m *Message) (err error) {
	//todo
	for group.MuxWSs != nil {
		err = group.MuxWSs[0].Send(m)
		return
	}
	return
}

func (group *Group) Handle(m *Message) {
	//log.Printf("group received %#v", m)

	if !group.client && m.Method != MessageMethodData {
		group.ServerHandleMessage(m)
		return
	}

	//get conn and send message
	//todo better way to find conn
	for {
		t := time.Now()
		for _, conn := range group.Conns {
			if conn.ID == m.ConnID {
				log.Printf("find conn id %x", conn.ID)
				err := conn.HandleMessage(m)
				if err != nil {
					log.Println(err.Error())
					return
				}
				return
			}
		}
		if time.Now().After(t.Add(time.Second * 3)) {
			err := errors.New("conn does not exist")
			log.Println(err.Error())
			return
		}
	}
	return
}

func (group *Group) AddMuxWS(muxWS *MuxWebSocket) (err error) {
	muxWS.group = group
	group.MuxWSs = append(group.MuxWSs, muxWS)
	group.Listen(muxWS)
	return
}

func (group *Group) Listen(muxWS *MuxWebSocket) {
	go func() {
		for {
			log.Println("ready to receive")
			m, err := muxWS.Receive()
			if err != nil {
				log.Printf(err.Error())
				return
			}

			go group.Handle(m)
		}
		return
	}()
}
