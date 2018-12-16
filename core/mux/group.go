package mux

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

type Group struct {
	client bool

	MuxWSs     []*MuxWebSocket
	muxWSMutex sync.Mutex

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
		sendMessageChan: make(chan *Message, 8),
	}
	return
}

func (group *Group) Send(m *Message) (err error) {
	group.sendMessageChan <- m
	return
}

func (group *Group) Handle(m *Message) {
	//log.Printf("group received %#v", m)

	//Server handle new mux conn request
	if !group.client && m.Method == MessageMethodDial {
		err := group.ServerAcceptDial(m)
		if err != nil {
			log.Println(err.Error())
		}
		return
	}

	if m.Method == MessageMethodClose {
		conn := group.GetConn(m.ConnID)
		if conn == nil {
			//err := errors.New(fmt.Sprintf("conn does not exist: %d", m.ConnID))
			//log.Println(err.Error())
			return
		}

		if conn.closed {
			//err := errors.New(fmt.Sprintf("conn already closed: %d", m.ConnID))
			//log.Println(err.Error())
			return
		}

		err := conn.Close()
		if err != nil {
			log.Println(err.Error())
		}
		return
	}

	//get conn and send message
	conn := group.GetConn(m.ConnID)
	if conn == nil {
		//debug log
		err := errors.New(fmt.Sprintf("conn does not exist: %x", m.ConnID))
		log.Println(err.Error())
		//log.Printf("%X %X %X %d", m.Method, m.ConnID, m.MessageID, m.Length)
		return
	}

	conn.HandleMessage(m)
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
	id = rand.Uint32()
	//group.connIDMutex.Lock()
	//group.connID++
	//id = group.connID
	//group.connIDMutex.Unlock()
	return
}

func (group *Group) AddMuxWS(muxWS *MuxWebSocket) {
	muxWS.group = group
	group.MuxWSs = append(group.MuxWSs, muxWS)
	group.Listen(muxWS)
	return
}

func (group *Group) DeleteMuxWS(id uint32) {
	group.muxWSMutex.Lock()
	for i, muxWS := range group.MuxWSs {
		if muxWS.ID == id {
			group.MuxWSs = append(group.MuxWSs[:i], group.MuxWSs[i+1:]...)
			group.muxWSMutex.Unlock()
			return
		}
	}

	group.muxWSMutex.Unlock()
	log.Printf("Cannot find muxWS: %d", id)
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
