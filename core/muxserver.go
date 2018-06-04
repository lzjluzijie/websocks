package core

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
