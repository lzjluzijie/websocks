package client

import (
	"crypto/tls"
	"net"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	"github.com/lzjluzijie/websocks/core"
)

type Config struct {
	ListenAddr string
	ServerURL  string

	SNI          string
	InsecureCert bool

	Mux bool
}

//GetClient return client from path
func (config *Config) GetClient() (client *WebSocksClient, err error) {
	//tackle config
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

		//todo mux

		CreatedAt: time.Now(),
		Stats:     core.NewStats(),
	}
	return
}
