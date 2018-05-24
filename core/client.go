package core

import (
	"io"
	"net"
	"net/url"
	"errors"
	"crypto/tls"

	"github.com/juju/loggo"
	"golang.org/x/net/websocket"
)

var logger = loggo.GetLogger("core")

type Client struct {
	LogLevel     loggo.Level
	ListenAddr   *net.TCPAddr
	URL          *url.URL
	Origin       string
	ServerName   string
	InsecureCert bool
	WSConfig     websocket.Config
}

func (client *Client) Listen() (err error) {
	logger.SetLogLevel(client.LogLevel)

	switch client.URL.Scheme {
	case "ws":
		client.Origin = "http://" + client.URL.Host
	case "wss":
		client.Origin = "https://" + client.URL.Host
	default:
		return errors.New("unknown scheme")
	}

	logger.Debugf(client.Origin)

	config, err := websocket.NewConfig(client.URL.String(), client.Origin)
	if err != nil {
		return
	}

	config.TlsConfig = &tls.Config{
		InsecureSkipVerify: client.InsecureCert,
	}
	if client.ServerName != "" {
		config.TlsConfig.ServerName = client.ServerName
	}
	client.WSConfig = *config

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

func (client *Client) handleConn(conn *net.TCPConn) (err error) {
	defer func() {
		if err != nil {
			logger.Debugf("Handle connection error: %s", err.Error())
		}
	}()

	defer conn.Close()

	conn.SetLinger(0)

	err = handShake(conn)
	if err != nil {
		return
	}

	_, host, err := getRequest(conn)
	if err != nil {
		return
	}

	logger.Debugf("Host: %s", host)

	_, err = conn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x08, 0x43})
	if err != nil {
		return
	}

	config := client.WSConfig
	config.Header = map[string][]string{
		"WebSocks-Host": {host},
	}

	ws, err := websocket.DialConfig(&config)
	if err != nil {
		return
	}

	defer ws.Close()

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
