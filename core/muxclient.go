package core

import (
	"net"
)

func (client *Client) OpenMux() (muxWS *MuxWebSocket, err error) {
	wsConn, _, err := client.Dialer.Dial(client.URL.String(), map[string][]string{
		"WebSocks-Mux": {"mux"},
	})

	if err != nil {
		return
	}

	ws := &WebSocket{
		conn: wsConn,
	}

	muxWS = NewMuxWebSocket(ws)
	return
}
func (client *Client) handleMuxConn(conn *net.TCPConn) {
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

	logger.Debugf("host: %s", host)

	muxConn := NewMuxConn(client.MuxWS)

	err = muxConn.DialMessage(host)
	if err != nil {
		logger.Errorf(err.Error())
		return
	}

	muxConn.Run(conn)
	return
}

//func (client *Client) Dial(conn *net.TCPConn, host string) {
//	muxConn := NewMuxConn(client.MuxWS)
//
//	//listen local conn and send message
//	go func() {
//		messageID := 1
//		buf := make([]byte, 32*1024)
//		for {
//			n, err := conn.Read(buf)
//			if err != nil {
//				logger.Errorf(err.Error())
//				return
//			}
//
//			println(n)
//
//			dataMessage := &Message{
//				Method:    MessageMethodData,
//				ConnID:    id,
//				MessageID: messageID,
//				Data:      buf[:n],
//			}
//
//			messageID++
//			client.MessageChan <- dataMessage
//		}
//	}()
//
//	go func() {
//		for {
//			_, err := conn.Write(<-dataChan)
//			if err != nil {
//				logger.Debugf(err.Error())
//				conn.Close()
//			}
//		}
//	}()
//
//	return
//}
//
//func (client *MuxClient) HandleMessage(m *Message) (err error) {
//	if m.Method != MessageMethodData {
//		return errors.New("unknown method")
//	}
//
//	connID := m.ConnID
//	c, ok := client.muxConnMap.Load(connID)
//	if !ok {
//		return errors.New("can not load conn")
//	}
//	conn := c.(*MuxConn)
//
//	go func() {
//		for {
//			if conn.DataID == m.MessageID {
//				conn.DataChan <- m.Data
//				return
//			}
//		}
//	}()
//
//	return
//}
//
//func (client *Client) handleConn(conn *net.TCPConn) {
//	defer conn.Close()
//
//	conn.SetLinger(0)
//
//	err := handShake(conn)
//	if err != nil {
//		logger.Errorf(err.Error())
//		return
//	}
//
//	_, host, err := getRequest(conn)
//	if err != nil {
//		logger.Errorf(err.Error())
//		return
//	}
//
//	_, err = conn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x08, 0x43})
//	if err != nil {
//		logger.Errorf(err.Error())
//		return
//	}
//
//	client.Dial(conn, host)
//	return
//}
