package mux

import (
	"errors"
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

func (group *Group) Send(m *Message) (err error) {
	//todo
	for group.MuxWSs != nil {
		err = group.MuxWSs[0].Send(m)
		return
	}
	return
}

func (group *Group) Receive(m *Message) (err error) {
	if !group.client && m.Method != MessageMethodData {
		group.HandleMessage(m)
	}

	//get conn and send message
	//todo better way to find conn
	for _, conn := range group.Conns {
		if conn.ID == m.ConnID {
			err = conn.HandleMessage(m)
			if err != nil {
				return
			}
		}
	}

	err = errors.New("conn does not exist")
	return
}

func (group *Group) AddMuxWS(muxWS *MuxWebSocket) (err error) {
	muxWS.group = group
	group.MuxWSs = append(group.MuxWSs, muxWS)
	muxWS.Listen()
	return
}
