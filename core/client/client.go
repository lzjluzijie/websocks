package client

import (
	"io"
	"net"
	"time"

	"net/url"

	"crypto/tls"

	"sync"

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

	dialer *websocket.Dialer
	//connMutex sync.Mutex
	//wsConns []*core.WebSocket
	muxWS *core.MuxWebSocket

	//todo enable mux
	Mux bool

	stopC chan int

	//statistics
	CreatedAt  time.Time
	Downloaded uint64
	Uploaded   uint64

	downloadMutex  sync.Mutex
	DownloadSpeed  uint64
	downloadSpeedA uint64
	uploadMutex    sync.Mutex
	UploadSpeed    uint64
	uploadSpeedA   uint64
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
		dialer: &websocket.Dialer{
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

	//status
	go func() {
		t := time.NewTicker(time.Second)
		for range t.C {
			client.downloadMutex.Lock()
			client.DownloadSpeed = client.downloadSpeedA
			client.downloadSpeedA = 0
			client.downloadMutex.Unlock()
			log.Infof("Download speed: %d", client.DownloadSpeed)
		}
	}()

	go func() {
		t := time.NewTicker(time.Second)
		for range t.C {
			client.uploadMutex.Lock()
			client.UploadSpeed = client.uploadSpeedA
			client.uploadSpeedA = 0
			client.uploadMutex.Unlock()
			log.Infof("Upload speed: %d", client.UploadSpeed)
		}
	}()

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

	//todo mux
	if client.Mux {
		client.DialMuxConn(lc.Host, conn)
		return
	}

	ws, err := client.DialWebSocket(core.NewHostHeader(lc.Host))
	if err != nil {
		log.Errorf(err.Error())
		return
	}

	lc.Run(ws)
	//client.DialWSConn(lc.Host, lc)
	return
}

//todo rewrite
func (client *WebSocksClient) DialWSConn(host string, conn io.ReadWriter) {
	ws, err := client.DialWebSocket(core.NewHostHeader(host))
	if err != nil {
		log.Errorf(err.Error())
		return
	}

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

func (client *WebSocksClient) DialWebSocket(header map[string][]string) (ws *core.WebSocket, err error) {
	wsConn, _, err := client.dialer.Dial(client.ServerURL.String(), header)
	if err != nil {
		return
	}

	ws = core.NewWebSocket(wsConn)
	ws.AddDownloaded = client.AddDownloaded
	ws.AddUploaded = client.AddUploaded
	//client.connMutex.Lock()
	//client.wsConns = append(client.wsConns, ws)
	//client.connMutex.Unlock()
	return
}
