package core

import "math/rand"

//CreateMuxConn create new mux connection for client
func CreateMuxConn(muxWS *MuxWebSocket) (conn *MuxConn) {
	return &MuxConn{
		ID:            rand.Uint64(),
		MuxWS:         muxWS,
		wait:          make(chan int),
		sendMessageID: new(uint64),
	}
}

//client dial remote
func (conn *MuxConn) DialMessage(host string) (err error) {
	m := &Message{
		Method:    MessageMethodDial,
		MessageID: 18446744073709551615,
		ConnID:    conn.ID,
		Data:      []byte(host),
	}

	//log.Printf("dial for %s", host)

	err = conn.MuxWS.SendMessage(m)
	if err != nil {
		return
	}

	//log.Printf("%d %s", conn.ID, host)
	return
}
