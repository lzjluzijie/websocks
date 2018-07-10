package client

import (
	"io"
	"net"
	"time"

	"net/url"

	"crypto/tls"

	"github.com/gorilla/websocket"
	"github.com/lzjluzijie/websocks/core"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

type WebSocksClientConfig struct {
	ListenAddr string
	ServerURL  string

	SNI          string
	InsecureCert bool

	Mux bool
}

type WebSocksClient struct {
	ServerURL  *url.URL
	ListenAddr *net.TCPAddr
	Dialer     *websocket.Dialer
	muxWS      *core.MuxWebSocket

	//todo enable mux
	Mux bool

	stopC chan int

	//statistics
	CreatedAt  time.Time
	Uploaded   int64
	Downloaded int64
}

func NewWebSocksClient(config *WebSocksClientConfig) (client *WebSocksClient) {
	serverURL, err := url.Parse(config.ServerURL)
	if err != nil {
		return
	}

	laddr, err := net.ResolveTCPAddr("tcp", config.ListenAddr)
	if err != nil {
		return
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: config.InsecureCert,
		ServerName:         config.SNI,
	}

	client = &WebSocksClient{
		ServerURL:  serverURL,
		ListenAddr: laddr,
		Dialer: &websocket.Dialer{
			ReadBufferSize:   4 * 1024,
			WriteBufferSize:  4 * 1024,
			HandshakeTimeout: 10 * time.Second,
			TLSClientConfig:  tlsConfig,
		},

		CreatedAt: time.Now(),
	}

	return
}

func (client *WebSocksClient) Listen() (err error) {
	listener, err := net.ListenTCP("tcp", client.ListenAddr)
	if err != nil {
		return err
	}

	log.Infof("Start to listen at %s", client.ListenAddr.String())

	defer listener.Close()

	if client.Mux {
		err := client.OpenMux()
		if err != nil {
			log.Debugf(err.Error())
			return err
		}

		go client.ListenMuxWS(client.muxWS)
	}

	go func() {
		client.stopC = make(chan int)
		<-client.stopC
		err = listener.Close()
		if err != nil {
			log.Errorf(err.Error())
			return
		}

		log.Infof("stopped")
	}()

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Debugf(err.Error())
			break
		}

		go client.HandleConn(conn)
	}
	return nil
}

func (client *WebSocksClient) Stop() {
	client.stopC <- 1911
	return
}

func (client *WebSocksClient) HandleConn(conn *net.TCPConn) {
	lc, err := NewLocalConn(conn)
	if err != nil {
		log.Debug(err.Error())
		return
	}

	if client.Mux {
		client.DialMuxConn(lc.Host, conn)
	} else {
		client.DialWSConn(lc.Host, lc)
	}

	return
}

func (client *WebSocksClient) DialWSConn(host string, conn io.ReadWriter) {
	wsConn, _, err := client.Dialer.Dial(client.ServerURL.String(), map[string][]string{
		"WebSocks-Host": {host},
	})

	if err != nil {
		log.Errorf(err.Error())
		return
	}
	defer wsConn.Close()

	log.Debugf("dialed ws for %s", host)

	ws := core.NewWebSocket(wsConn)

	go func() {
		_, err = io.Copy(ws, conn)
		if err != nil {
			log.Debugf(err.Error())
			return
		}
		return
	}()

	_, err = io.Copy(conn, ws)
	if err != nil {
		log.Debugf(err.Error())
		return
	}
	return
}
