package client

import (
	"net"
	"time"

	"net/url"

	"github.com/gorilla/websocket"
	"github.com/lzjluzijie/websocks/core"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

type WebSocksClient struct {
	ServerURL  *url.URL
	ListenAddr *net.TCPAddr

	dialer *websocket.Dialer
	//connMutex sync.Mutex
	//wsConns []*core.WebSocket
	muxWS *core.MuxWebSocket

	//todo enable mux
	Mux bool

	//control
	stopC chan int

	CreatedAt time.Time
	Stats     *core.Stats
}

func (client *WebSocksClient) Run() (err error) {
	listener, err := net.ListenTCP("tcp", client.ListenAddr)
	if err != nil {
		return err
	}

	log.Infof("Start to listen at %s", client.ListenAddr.String())

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
	return
}

func (client *WebSocksClient) DialWebSocket(header map[string][]string) (ws *core.WebSocket, err error) {
	wsConn, _, err := client.dialer.Dial(client.ServerURL.String(), header)
	if err != nil {
		return
	}

	ws = core.NewWebSocket(wsConn, client.Stats)
	//client.connMutex.Lock()
	//client.wsConns = append(client.wsConns, ws)
	//client.connMutex.Unlock()
	return
}
