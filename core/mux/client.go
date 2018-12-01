package mux

import "math/rand"

//NewMuxConn creates a new mux connection for client
func (group *Group) NewMuxConn(host string) (conn *Conn, err error) {
	conn = &Conn{
		ID:            rand.Uint32(),
		wait:          make(chan int),
		sendMessageID: new(uint32),
		group:         group,
	}

	m := &Message{
		Method:    MessageMethodDial,
		MessageID: 4294967295,
		ConnID:    conn.ID,
		Length:    uint32(len(host)),
		Data:      []byte(host),
	}

	err = group.Send(m)
	if err != nil {
		return
	}

	group.Conns = append(group.Conns, conn)
	return
}
