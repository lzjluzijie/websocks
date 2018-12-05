package mux

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

type Group struct {
	client bool

	MuxWSs []*MuxWebSocket

	sendMessageChan chan *Message

	connMap      map[uint32]*Conn
	connMapMutex sync.RWMutex

	connID      uint32
	connIDMutex sync.Mutex
}

//true: client group
//false: server group
func NewGroup(client bool) (group *Group) {
	group = &Group{
		client:          client,
		connMap:         make(map[uint32]*Conn),
		sendMessageChan: make(chan *Message, 1911),
	}
	return
}

func (group *Group) Send(m *Message) (err error) {
	group.sendMessageChan <- m
	return
}

func (group *Group) Handle(m *Message) {
	//log.Printf("group received %#v", m)

	if !group.client && m.Method != MessageMethodData {
		group.ServerHandleMessage(m)
		return
	}

	//get conn and send message
	for {
		conn := group.GetConn(m.ConnID)
		if conn == nil {
			//debug log
			err := errors.New(fmt.Sprintf("conn does not exist: %x", m.ConnID))
			log.Println(err.Error())
			log.Printf("%X %X %X %d", m.Method, m.ConnID, m.MessageID, m.Length)
			return
		}

		//this err should be nil or ErrConnClosed
		err := conn.HandleMessage(m)
		if err != nil {
			log.Println(err.Error())
			return
		}
	}
	return
}

func (group *Group) AddConn(conn *Conn) {
	group.connMapMutex.Lock()
	group.connMap[conn.ID] = conn
	group.connMapMutex.Unlock()
	return
}

func (group *Group) DeleteConn(id uint32) {
	group.connMapMutex.Lock()
	delete(group.connMap, id)
	group.connMapMutex.Unlock()
	return
}

func (group *Group) GetConn(id uint32) (conn *Conn) {
	group.connMapMutex.RLock()
	conn = group.connMap[id]
	group.connMapMutex.RUnlock()

	if conn == nil {
		t := time.Now()
		for time.Now().Before(t.Add(time.Second)) {
			group.connMapMutex.RLock()
			conn = group.connMap[id]
			group.connMapMutex.RUnlock()
			if conn != nil {
				return conn
			}
		}
	}
	return
}

func (group *Group) NextConnID() (id uint32) {
	group.connIDMutex.Lock()
	group.connID++
	id = group.connID
	group.connIDMutex.Unlock()
	return
}

func (group *Group) AddMuxWS(muxWS *MuxWebSocket) {
	muxWS.group = group
	group.MuxWSs = append(group.MuxWSs, muxWS)
	group.Listen(muxWS)
	return
}

func (group *Group) Listen(muxWS *MuxWebSocket) {
	go func() {
		for {
			m, err := muxWS.Receive()
			if err != nil {
				log.Println(err.Error())
				return
			}

			go group.Handle(m)
		}
	}()

	go func() {
		for {
			m := <-group.sendMessageChan
			err := muxWS.Send(m)
			if err != nil {
				log.Println(err.Error())
				return
			}
		}
	}()
}
