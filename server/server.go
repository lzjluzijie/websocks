package server

import (
	"io"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/juju/loggo"
	"github.com/lzjluzijie/websocks/core"
)

//todo
var logger = loggo.GetLogger("server")

type Config struct {
	ListenAddr   string
	Pattern      string
	TLS          bool
	CertPath     string
	KeyPath      string
	ReverseProxy string
}

type WebSocksServer struct {
	*Config
	LogLevel loggo.Level

	Upgrader   *websocket.Upgrader
	muxConnMap sync.Map
	mutex      sync.Mutex

	CreatedAt time.Time
	Stats     *core.Stats
}

func (config *Config) NewWebSocksServer() (server *WebSocksServer) {
	server = &WebSocksServer{
		Config: config,
		Upgrader: &websocket.Upgrader{
			ReadBufferSize:   4 * 1024,
			WriteBufferSize:  4 * 1024,
			HandshakeTimeout: 10 * time.Second,
		},
		CreatedAt: time.Now(),
		Stats:     core.NewStats(),
	}
	return
}

func (server *WebSocksServer) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	c, err := server.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Debugf(err.Error())
		return
	}
	defer c.Close()

	ws := core.NewWebSocket(c, server.Stats)
	//todo conns

	////mux
	//if r.Header.Get("WebSocks-Mux") == "mux" {
	//	muxWS := NewMuxWebSocket(ws)
	//	muxWS.ServerListen()
	//	return
	//}

	host := r.Header.Get("WebSocks-Host")
	logger.Debugf("Dial %s", host)
	conn, err := server.DialRemote(host)
	if err != nil {
		logger.Debugf(err.Error())
		return
	}

	go func() {
		_, err = io.Copy(conn, ws)
		if err != nil {
			logger.Debugf(err.Error())
			return
		}
	}()

	_, err = io.Copy(ws, conn)
	if err != nil {
		logger.Debugf(err.Error())
		return
	}

	return
}

func (server *WebSocksServer) DialRemote(host string) (conn net.Conn, err error) {
	conn, err = net.Dial("tcp", host)
	if err != nil {
		return
	}
	return
}

func (server *WebSocksServer) Listen() (err error) {
	logger.SetLogLevel(server.LogLevel)

	s := http.Server{
		Addr:    server.ListenAddr,
		Handler: server.getMacaron(),
	}

	logger.Infof("Start to listen at %s", server.ListenAddr)

	if !server.TLS {
		err = s.ListenAndServe()
		if err != nil {
			return err
		}
		return
	} else {

	}

	return
}
