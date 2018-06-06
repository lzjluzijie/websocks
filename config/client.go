package config

import (
	"crypto/tls"
	"net"
	"net/url"

	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/gorilla/websocket"
	"github.com/lzjluzijie/websocks/core"
	"github.com/urfave/cli"
)

//GetClient return client from path
func GetClientConfig(path string) (client *core.Client, err error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	//read config
	config := &core.ClientConfig{}
	err = json.Unmarshal(data, config)
	if err != nil {
		return
	}

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

	client = &core.Client{
		ClientConfig: config,

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

//GenerateClientConfig create a client config from cli.Context
func GenerateClientConfig(c *cli.Context) (err error) {
	path := c.String("path")

	config := &core.ClientConfig{
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
