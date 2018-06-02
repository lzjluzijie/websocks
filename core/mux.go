package core

import (
	"encoding/json"
	"time"

	"github.com/xtaci/smux"
)

type MuxRequest struct {
	Host string
}

func (client *Client) OpenSession() (err error) {
	wsConn, _, err := client.Dialer.Dial(client.URL.String(), map[string][]string{
		"WebSocks-Mux": {"mux"},
	})

	if err != nil {
		logger.Errorf(err.Error())
		return
	}

	ws := &WebSocket{
		conn: wsConn,
	}

	session, err := smux.Client(ws, nil)
	if err != nil {
		return
	}

	go func() {
		for {
			if session.NumStreams() > 2 {
				time.Sleep(time.Second)
				continue
			}

			stream, err := session.OpenStream()
			if err != nil {
				session.Close()
				logger.Errorf(err.Error())
				return
			}

			client.StreamChan <- stream
		}
		return
	}()

	return
}

func (client *Client) GetStream(host string) (stream *smux.Stream, err error) {
	stream = <-client.StreamChan

	req := MuxRequest{
		Host: host,
	}

	enc := json.NewEncoder(stream)
	err = enc.Encode(req)
	return
}
