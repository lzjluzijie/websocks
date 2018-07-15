package server

import (
	"time"

	"github.com/gorilla/websocket"
	"github.com/juju/loggo"
	"github.com/lzjluzijie/websocks/core"
)

type Config struct {
	ListenAddr   string
	Pattern      string
	TLS          bool
	CertPath     string
	KeyPath      string
	ReverseProxy string
}

func (config *Config) GetServer() (server *WebSocksServer) {
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

	logger.SetLogLevel(loggo.DEBUG)
	return
}
