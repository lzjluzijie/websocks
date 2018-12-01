package mux

import "math/rand"

//NewMuxConn creates a new mux connection for client
func (group *Group) NewMuxConn(host string) (err error) {
	conn := &Conn{
		ID:            rand.Uint32(),
		wait:          make(chan int),
		sendMessageID: new(uint32),
	}

	mh := &MessageHead{
		Method:    MessageMethodDial,
		MessageID: 4294967295,
		ConnID:    conn.ID,
		Length:    uint32(len(host)),
	}

	err = group.Send(mh, []byte(host))
	return
}
