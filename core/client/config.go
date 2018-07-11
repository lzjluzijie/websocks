package client

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	"github.com/urfave/cli"
)

//GetClient return client from path
func GetClient(config *WebSocksClientConfig) (client *WebSocksClient, err error) {
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
	}
	return
}

//GenerateClientConfig create a client config from cli.Context
func GenerateClientConfig(c *cli.Context) (err error) {
	path := c.String("path")

	config := &WebSocksClientConfig{
		ListenAddr:   c.String("l"),
		ServerURL:    c.String("s"),
		SNI:          c.String("sni"),
		InsecureCert: c.Bool("insecure"),
		Mux:          c.Bool("mux"),
	}

	data, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return
	}

	err = ioutil.WriteFile(path, data, 600)
	if err != nil {
		return
	}
	return
}
