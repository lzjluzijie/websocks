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

	client = &core.Client{
		ClientConfig: config,

		Dialer: &websocket.Dialer{
			ReadBufferSize:   4 * 1024,
			WriteBufferSize:  4 * 1024,
			HandshakeTimeout: 10 * time.Second,
			TLSClientConfig:  config.TLSConfig,
		},
		CreatedAt: time.Now(),
	}
	return
}

//GenerateClientConfig create a client config from cli.Context
func GenerateClientConfig(c *cli.Context) (err error) {
	path := c.String("path")
	listenAddr := c.String("l")
	serverURL := c.String("s")
	mux := c.Bool("mux")
	serverName := c.String("n")
	insecureCert := false
	if c.Bool("insecure") {
		insecureCert = true
	}

	u, err := url.Parse(serverURL)
	if err != nil {
		return
	}

	lAddr, err := net.ResolveTCPAddr("tcp", listenAddr)
	if err != nil {
		return
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: insecureCert,
	}

	if serverName != "" {
		tlsConfig.ServerName = serverName
	}

	config := &core.ClientConfig{
		ListenAddr: lAddr,
		URL:        u,
		Mux:        mux,
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
