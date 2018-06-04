package core

import (
	"errors"
	"fmt"
	"net"
)

func (muxWS *MuxWebSocket) ServerListen() {
	//block and listen
	for {
		m, err := muxWS.ReceiveMessage()
		if err != nil {
			logger.Debugf(err.Error())
			return
		}

		go muxWS.ServerHandleMessage(m)
	}
	return
}

func (muxWS *MuxWebSocket) AcceptMuxConn(m *Message) (conn *MuxConn, host string, err error) {
	if m.Method != MessageMethodDial {
		err = errors.New(fmt.Sprintf("wrong message method %d", m.Method))
		return
	}

	host = string(m.Data)

	conn = &MuxConn{
		ID:            m.ConnID,
		muxWS:         muxWS,
		wait:          make(chan int),
		sendMessageID: new(uint64),
	}
	muxWS.PutMuxConn(conn)
	return
}

func (muxWS *MuxWebSocket) ServerHandleMessage(m *Message) {
	//accept new conn
	if m.Method == MessageMethodDial {
		conn, host, err := muxWS.AcceptMuxConn(m)
		if err != nil {
			logger.Debugf(err.Error())
			return
		}

		tcpAddr, err := net.ResolveTCPAddr("tcp", host)
		if err != nil {
			logger.Debugf(err.Error())
			return
		}

		tcpConn, err := net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			logger.Debugf(err.Error())
			return
		}

		logger.Debugf("Accepted mux conn %s", host)

		conn.Run(tcpConn)
		return
	}

	//get conn and send message
	conn := muxWS.GetMuxConn(m.ConnID)
	err := conn.HandleMessage(m)
	if err != nil {
		logger.Debugf(err.Error())
		return
	}
}

//func (server *Server) HandleMuxWS(ws *WebSocket) (muxWS *MuxWebSocket,err error) {
//	dec := gob.NewDecoder(ws)
//	enc := gob.NewEncoder(ws)
//
//	//receive messages
//	go func() {
//		for {
//			m := &Message{}
//			err = dec.Decode(m)
//			if err != nil {
//				logger.Debugf(err.Error())
//				return
//			}
//
//			err = server.HandleMessage(m)
//			if err != nil {
//				logger.Debugf(err.Error())
//				continue
//			}
//		}
//	}()
//
//	//send messages
//	go func() {
//		for {
//			m := <-server.MessageChan
//			err = enc.Encode(m)
//			if err != nil {
//				logger.Debugf(err.Error())
//				return
//			}
//		}
//	}()
//
//	time.Sleep(time.Minute)
//	return
//}
//
//func (server *Server) HandleMessage(m *Message) (err error) {
//	if m.Method == MessageMethodDial {
//		id := m.ConnID
//		dataChan := make(chan []byte)
//		conn := &MuxConn{
//			ID:       id,
//			DataChan: dataChan,
//		}
//
//		server.muxConnMap.Store(id, conn)
//		server.DialRemote(conn, string(m.Data))
//		return
//	}
//
//	if m.Method != MessageMethodData {
//		return errors.New("unknown method")
//	}
//
//	connID := m.ConnID
//	c, ok := server.muxConnMap.Load(connID)
//	if !ok {
//		return errors.New("can not load conn")
//	}
//
//	conn := c.(*MuxConn)
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
//func (server *Server) DialRemote(muxConn *MuxConn, host string) {
//	conn, err := net.Dial("tcp", host)
//	if err != nil {
//		logger.Debugf(err.Error())
//		return
//	}
//
//	go func() {
//		for {
//			buf := make([]byte, 32*1024)
//			n, err := conn.Read(buf)
//			if err != nil {
//				logger.Debugf(err.Error())
//				return
//			}
//
//			m := &Message{
//				Method:    MessageMethodData,
//				ConnID:    muxConn.ID,
//				MessageID: muxConn.DataID,
//				Data:      buf[:n],
//			}
//			muxConn.DataID++
//
//			server.MessageChan <- m
//		}
//	}()
//
//	return
//}
