package client

import (
	"github.com/lzjluzijie/websocks/core"
	"github.com/lzjluzijie/websocks/core/mux"
)

func (client *WebSocksClient) OpenMux() (err error) {
	wsConn, _, err := client.dialer.Dial(client.ServerURL.String(), map[string][]string{
		"WebSocks-Mux": {"v0.16"},
	})

	if err != nil {
		return
	}

	ws := core.NewWebSocket(wsConn, client.Stats)

	muxWS := mux.NewMuxWebSocket(ws)
	client.muxGroup.AddMuxWS(muxWS)

	//log.Printf("new mux websocket")
	return
}
