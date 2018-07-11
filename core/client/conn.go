package client

import (
	"errors"
	"io"
	"net"
	"time"
)

type LocalConn struct {
	Host string

	conn *net.TCPConn

	//stats
	createdAt time.Time
	closed    bool
}

func NewLocalConn(conn *net.TCPConn) (lc *LocalConn, err error) {
	conn.SetLinger(0)
	err = handShake(conn)
	if err != nil {
		return
	}

	_, host, err := getRequest(conn)
	if err != nil {
		return
	}

	_, err = conn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x08, 0x43})
	if err != nil {
		log.Debugf(err.Error())
		return
	}

	lc = &LocalConn{
		Host: host,

		conn:      conn,
		createdAt: time.Now(),
	}
	return
}

func (lc *LocalConn) Run(conn io.ReadWriter) {
	go func() {
		_, err := io.Copy(lc, conn)
		if err != nil {
			log.Debugf(err.Error())
			return
		}
		return
	}()

	go func() {
		_, err := io.Copy(conn, lc)
		if err != nil {
			log.Debugf(err.Error())
			return
		}
	}()
	return
}

func (lc *LocalConn) Read(p []byte) (n int, err error) {
	if lc.closed {
		return 0, errors.New("local conn closed")
	}

	n, err = lc.conn.Read(p)
	if err != nil {
		lc.closed = true
	}
	return
}

func (lc *LocalConn) Write(p []byte) (n int, err error) {
	if lc.closed {
		return 0, errors.New("local conn closed")
	}

	n, err = lc.conn.Write(p)
	if err != nil {
		lc.closed = true
	}
	return
}
