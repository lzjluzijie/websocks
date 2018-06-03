package core

import (
	"encoding/gob"
	"errors"
	"math/rand"
	"net"
	"sync"
)

type MuxClient struct {
	*Client
	MessageChan chan *Message
	muxConnMap  sync.Map
	Mutex       sync.Mutex
}

func (client *MuxClient) Open() (err error) {
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

	dec := gob.NewDecoder(ws)
	enc := gob.NewEncoder(ws)

	go func() {
		for {
			m := &Message{}
			err = dec.Decode(m)
			if err != nil {
				logger.Debugf(err.Error())
				return
			}

			err = client.HandleMessage(m)
			if err != nil {
				logger.Debugf(err.Error())
				continue
			}
		}
	}()

	go func() {
		for {
			m := <-client.MessageChan
			err = enc.Encode(m)
			if err != nil {
				logger.Debugf(err.Error())
				return
			}
		}
	}()
	return
}

func (client *MuxClient) Dial(conn *net.TCPConn, host string) {
	dataChan := make(chan []byte)
	id := rand.Int()
	muxConn := &MuxConn{
		DataChan: dataChan,
		ID:       id,
	}

	client.muxConnMap.Store(id, muxConn)

	m := &Message{
		Method: MessageMethodDial,
		ConnID: id,
		Data:   []byte(host),
	}
	client.MessageChan <- m

	//listen local conn and send message
	go func() {
		messageID := 1
		buf := make([]byte, 32*1024)
		for {
			n, err := conn.Read(buf)
			if err != nil {
				logger.Errorf(err.Error())
				return
			}

			println(n)

			dataMessage := &Message{
				Method:    MessageMethodData,
				ConnID:    id,
				MessageID: messageID,
				Data:      buf[:n],
			}

			messageID++
			client.MessageChan <- dataMessage
		}
	}()

	go func() {
		for {
			_, err := conn.Write(<-dataChan)
			if err != nil {
				logger.Debugf(err.Error())
				conn.Close()
			}
		}
	}()

	return
}

func (client *MuxClient) HandleMessage(m *Message) (err error) {
	if m.Method != MessageMethodData {
		return errors.New("unknown method")
	}

	connID := m.ConnID
	c, ok := client.muxConnMap.Load(connID)
	if !ok {
		return errors.New("can not load conn")
	}
	conn := c.(*MuxConn)

	go func() {
		for {
			if conn.DataID == m.MessageID {
				conn.DataChan <- m.Data
				return
			}
		}
	}()

	return
}

func (client *MuxClient) handleConn(conn *net.TCPConn) {
	defer conn.Close()

	conn.SetLinger(0)

	err := handShake(conn)
	if err != nil {
		logger.Errorf(err.Error())
		return
	}

	_, host, err := getRequest(conn)
	if err != nil {
		logger.Errorf(err.Error())
		return
	}

	_, err = conn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x08, 0x43})
	if err != nil {
		logger.Errorf(err.Error())
		return
	}

	client.Dial(conn, host)
	return
}
