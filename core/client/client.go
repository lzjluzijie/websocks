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

	dialer *websocket.Dialer
	//connMutex sync.Mutex
	//wsConns []*core.WebSocket
	muxWS *core.MuxWebSocket

	//todo enable mux
	Mux bool

	stopC chan int

	//statistics
	CreatedAt     time.Time
	Downloaded    uint64
	Uploaded      uint64
	DownloadSpeed uint64
	UploadSpeed   uint64
	downloadedC   chan uint64
	uploadedC     chan uint64
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

		CreatedAt:   time.Now(),
		downloadedC: make(chan uint64, 0),
		uploadedC:   make(chan uint64, 0),
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
			downloadSpeed := uint64(0)
			for i := len(client.uploadedC); i > 0; i-- {
				downloadSpeed += <-client.uploadedC
			}

			client.DownloadSpeed = downloadSpeed
			log.Infof("Download speed: %d", downloadSpeed)
		}
	}()

	go func() {
		t := time.NewTicker(time.Second)
		for range t.C {
			uploadSpeed := uint64(0)
			for i := len(client.uploadedC); i > 0; i-- {
				uploadSpeed += <-client.uploadedC
			}

			client.UploadSpeed = uploadSpeed
			log.Infof("Upload speed: %d", uploadSpeed)
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
