package mux

//NewMuxConn creates a new mux connection for client
func (group *Group) NewMuxConn(host string) (conn *Conn, err error) {
	conn = &Conn{
		ID:    group.NextConnID(),
		wait:  make(chan int),
		group: group,

		sendMessageNext:    1,
		receiveMessageNext: 1,
	}

	m := &Message{
		Method:    MessageMethodDial,
		MessageID: 0,
		ConnID:    conn.ID,
		Length:    uint32(len(host)),
		Data:      []byte(host),
	}

	err = group.Send(m)
	if err != nil {
		return
	}

	group.AddConn(conn)
	return
}
