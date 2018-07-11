package client

import (
	"net"

	"github.com/lzjluzijie/websocks/core"
)

func (client *WebSocksClient) OpenMux() (err error) {
	wsConn, _, err := client.dialer.Dial(client.ServerURL.String(), map[string][]string{
		"WebSocks-Mux": {"mux"},
	})

	if err != nil {
		return
	}

	ws := core.NewWebSocket(wsConn)

	muxWS := core.NewMuxWebSocket(ws)
	client.muxWS = muxWS
	return
}

func (client *WebSocksClient) DialMuxConn(host string, conn *net.TCPConn) {
	muxConn := core.NewMuxConn(client.muxWS)

	err := muxConn.DialMessage(host)
	if err != nil {
		log.Errorf(err.Error())
		err = client.OpenMux()
		if err != nil {
			log.Errorf(err.Error())
		}
		return
	}

	muxConn.MuxWS.PutMuxConn(muxConn)

	log.Debugf("dialed mux for %s", host)

	muxConn.Run(conn)
	return
}

func (client *WebSocksClient) ListenMuxWS(muxWS *core.MuxWebSocket) {
	for {
		m, err := muxWS.ReceiveMessage()
		if err != nil {
			log.Debugf(err.Error())
			return
		}

		//get conn and send message
		conn := muxWS.GetMuxConn(m.ConnID)
		err = conn.HandleMessage(m)
		if err != nil {
			log.Debugf(err.Error())
			continue
		}
	}
}
