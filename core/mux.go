package core

import (
	"encoding/gob"
	"math/rand"
	"net"
	"time"
)

type Mux struct {
	WS      *WebSocket
	Decoder *gob.Decoder
	Encoder *gob.Encoder
}

type MuxRequest struct {
	ID     uint64
	Method string
	Data   []byte
}

type MuxResponse struct {
	ID   uint64
	Data []byte
}

//DialMux dial new mux conn, listen and write to local conn
func (client *Client) DialMux() (mux *Mux, err error) {
	conn, _, err := client.Dialer.Dial(client.URL.String(), map[string][]string{
		"WebSocks-Mux": {"mux"},
	})

	if err != nil {
		return
	}

	ws := &WebSocket{
		conn: conn,
	}

	mux = &Mux{
		WS:      ws,
		Decoder: gob.NewDecoder(ws),
		Encoder: gob.NewEncoder(ws),
	}

	go func() {
		for {
			resp := &MuxResponse{}
			err := mux.Decoder.Decode(resp)
			if err != nil {
				logger.Errorf(err.Error())
				return
			}

			conn := client.MuxTCPConn[resp.ID]
			_, err = conn.Write(resp.Data)
			if err != nil {
				logger.Errorf(err.Error())
				continue
			}
		}
	}()

	return
}

func (server *Server) ServerHandleMux(ws *WebSocket) {
	dec := gob.NewDecoder(ws)
	enc := gob.NewEncoder(ws)

	for {
		req := &MuxRequest{}
		err := dec.Decode(req)
		if err != nil {
			logger.Errorf(err.Error())
			continue
		}

		go server.ServerHandleMuxRequest(req, enc)
	}
	return
}

func (server *Server) ServerHandleMuxRequest(req *MuxRequest, enc *gob.Encoder) {
	id := req.ID
	method := req.Method
	var err error

	//If method is "new", dial new and listen
	if method == "new" {
		host := string(req.Data)
		conn, err := net.Dial("tcp", host)
		if err != nil {
			logger.Errorf(err.Error())
			return
		}

		server.MuxConn[id] = conn
		logger.Debugf("dialed %s, id %d", host, id)

		go func() {
			data := make([]byte, 32*1024)
			for {
				n, err := conn.Read(data)
				if err != nil {
					logger.Errorf(err.Error())
					return
				}

				resp := &MuxResponse{
					ID:   id,
					Data: data[:n],
				}

				err = enc.Encode(resp)
				if err != nil {
					logger.Errorf(err.Error())
					continue
				}
			}
		}()

		return
	}

	//If method is not "new", write data
	conn := server.MuxConn[id]
	if conn == nil {
		time.Sleep(3 * time.Second)
		conn = server.MuxConn[id]
		if conn == nil {
			logger.Errorf("conn %d does not exist", id)
		}
		return
	}
	_, err = conn.Write(req.Data)
	if err != nil {
		logger.Errorf(err.Error())
		conn.Close()
		delete(server.MuxConn, id)
		logger.Debugf("conn %d closed", id)
		return
	}
}

//ClientHandleMux read from local conn and write to remote server
func (client *Client) ClientHandleMux(conn *net.TCPConn, host string) {
	mux := client.MuxS[0]
	id, err := mux.NewConn(host)
	if err != nil {
		logger.Errorf(err.Error())
		return
	}

	client.MuxTCPConn[id] = conn

	data := make([]byte, 32*1024)

	for {
		n, err := conn.Read(data)
		if err != nil {
			logger.Errorf(err.Error())
			return
		}

		req := &MuxRequest{
			ID:   id,
			Data: data[:n],
		}

		err = mux.Encoder.Encode(req)
		if err != nil {
			logger.Errorf(err.Error())
			continue
		}
	}
}

func (mux *Mux) NewConn(host string) (id uint64, err error) {
	id = rand.Uint64()
	req := &MuxRequest{
		ID:     id,
		Method: "new",
		Data:   []byte(host),
	}

	err = mux.Encoder.Encode(req)
	if err != nil {
		return
	}

	logger.Debugf("mux dialed %s", host)
	return
}
