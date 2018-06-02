package core

import (
	"io"
	"net"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	"github.com/juju/loggo"
	"github.com/xtaci/smux"
)

var logger = loggo.GetLogger("core")

type Client struct {
	LogLevel   loggo.Level
	ListenAddr *net.TCPAddr
	URL        *url.URL

	Mux        bool
	Opened     int
	StreamChan chan *smux.Stream

	Dialer *websocket.Dialer

	CreatedAt time.Time
}

func (client *Client) Listen() (err error) {
	logger.SetLogLevel(client.LogLevel)

	listener, err := net.ListenTCP("tcp", client.ListenAddr)
	if err != nil {
		return err
	}

	logger.Infof("Start to listen at %s", client.ListenAddr.String())

	defer listener.Close()

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			logger.Debugf(err.Error())
			continue
		}

		go client.handleConn(conn)
	}

	return nil
}

func (client *Client) handleConn(conn *net.TCPConn) {
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

	if client.Mux {
		l := len(client.StreamChan)
		c := cap(client.StreamChan)
		logger.Debugf("%d %d", l, c)
		if l != c {
			go func() {
				err := client.OpenSession()
				if err != nil {
					logger.Errorf(err.Error())
					return
				}
			}()
		}

		stream, err := client.GetStream(host)
		if err != nil {
			logger.Errorf(err.Error())
			return
		}

		go func() {
			_, err = io.Copy(stream, conn)
			if err != nil {
				logger.Debugf(err.Error())
				stream.Close()
				return
			}
			return
		}()

		_, err = io.Copy(conn, stream)
		if err != nil {
			logger.Errorf(err.Error())
			stream.Close()
			return
		}

		return
	}

	wsConn, _, err := client.Dialer.Dial(client.URL.String(), map[string][]string{
		"WebSocks-Host": {host},
	})

	if err != nil {
		return
	}

	ws := &WebSocket{
		conn: wsConn,
	}

	go func() {
		_, err = io.Copy(ws, conn)
		if err != nil {
			logger.Debugf(err.Error())
			return
		}
		return
	}()

	_, err = io.Copy(conn, ws)
	if err != nil {
		return
	}

	return
}
