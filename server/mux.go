package server

//import (
//	"errors"
//	"fmt"
//	"net"
//	"time"
//)
//
//func (muxWS *MuxWebSocket) ServerListen() {
//	//block and listen
//	for {
//		m, err := muxWS.ReceiveMessage()
//		if err != nil {
//			logger.Debugf(err.Error())
//			return
//		}
//
//		go muxWS.ServerHandleMessage(m)
//	}
//	return
//}
//
//func (muxWS *MuxWebSocket) ServerHandleMessage(m *Message) {
//	//check message
//	if m.Data == nil {
//		return
//	}
//
//	//accept new conn
//	if m.Method == MessageMethodDial {
//		conn, host, err := muxWS.AcceptMuxConn(m)
//		if err != nil {
//			logger.Debugf(err.Error())
//			return
//		}
//
//		tcpAddr, err := net.ResolveTCPAddr("tcp", host)
//		if err != nil {
//			logger.Debugf(err.Error())
//			return
//		}
//
//		tcpConn, err := net.DialTCP("tcp", nil, tcpAddr)
//		if err != nil {
//			logger.Debugf(err.Error())
//			return
//		}
//
//		logger.Debugf("Accepted mux conn %s", host)
//
//		conn.Run(tcpConn)
//		return
//	}
//
//	//get conn and send message
//	conn := muxWS.GetMuxConn(m.ConnID)
//	if conn == nil {
//		time.Sleep(time.Second)
//		conn = muxWS.GetMuxConn(m.ConnID)
//		if conn == nil {
//			logger.Debugf("conn %d do not exist", m.ConnID)
//			return
//		}
//	}
//	err := conn.HandleMessage(m)
//	if err != nil {
//		logger.Debugf(err.Error())
//		return
//	}
//}
//
//func (muxWS *MuxWebSocket) AcceptMuxConn(m *Message) (conn *MuxConn, host string, err error) {
//	if m.Method != MessageMethodDial {
//		err = errors.New(fmt.Sprintf("wrong message method %d", m.Method))
//		return
//	}
//
//	host = string(m.Data)
//
//	conn = &MuxConn{
//		ID:            m.ConnID,
//		MuxWS:         muxWS,
//		wait:          make(chan int),
//		sendMessageID: new(uint64),
//	}
//	muxWS.PutMuxConn(conn)
//	return
//}
