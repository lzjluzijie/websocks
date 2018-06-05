package core

import (
	"net"
)

func (muxWS *MuxWebSocket) ClientListen() {
	for {
		m, err := muxWS.ReceiveMessage()
		if err != nil {
			logger.Debugf(err.Error())
			return
		}

		//get conn and send message
		conn := muxWS.GetMuxConn(m.ConnID)
		err = conn.HandleMessage(m)
		if err != nil {
			logger.Debugf(err.Error())
			continue
		}
	}
}

func (client *Client) OpenMux() (err error) {
	wsConn, _, err := client.Dialer.Dial(client.URL.String(), map[string][]string{
		"WebSocks-Mux": {"mux"},
	})

	if err != nil {
		return
	}

	ws := &WebSocket{
		conn: wsConn,
	}

	muxWS := NewMuxWebSocket(ws)
	client.MuxWS = muxWS
	return
}
func (client *Client) DialMuxConn(host string, conn *net.TCPConn) {
	muxConn := NewMuxConn(client.MuxWS)

	err := muxConn.DialMessage(host)
	if err != nil {
		logger.Errorf(err.Error())
		err = client.OpenMux()
		if err != nil {
			logger.Errorf(err.Error())
		}
		return
	}

	muxConn.muxWS.PutMuxConn(muxConn)

	logger.Debugf("dialed mux for %s", host)

	muxConn.Run(conn)
	return
}

//client dial remote
func (conn *MuxConn) DialMessage(host string) (err error) {
	m := &Message{
		Method:    MessageMethodDial,
		MessageID: 18446744073709551615,
		ConnID:    conn.ID,
		Data:      []byte(host),
	}

	logger.Debugf("dial for %s", host)

	err = conn.muxWS.SendMessage(m)
	if err != nil {
		return
	}

	logger.Debugf("%d %s", conn.ID, host)
	return
}
