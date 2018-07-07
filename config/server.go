package config

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/gorilla/websocket"
	"github.com/lzjluzijie/websocks/core"
	"github.com/urfave/cli"
)

//GenerateServerConfig create a client config from path
func GetServerConfig(path string) (server *core.Server, err error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	//read config
	config := &core.ServerConfig{}
	err = json.Unmarshal(data, config)
	if err != nil {
		return
	}

	server = &core.Server{
		ServerConfig: config,
		Upgrader: &websocket.Upgrader{
			ReadBufferSize:   4 * 1024,
			WriteBufferSize:  4 * 1024,
			HandshakeTimeout: 10 * time.Second,
		},
		CreatedAt: time.Now(),
	}
	return
}

//GenerateServerConfig create a server config from cli.Context
func GenerateServerConfig(c *cli.Context) (err error) {
	path := c.String("path")
	listenAddr := c.String("l")
	pattern := c.String("pattern")
	tls := c.Bool("tls")
	certPath := c.String("cert")
	keyPath := c.String("key")
	reverseProxy := c.String("reverse-proxy")

	//if []byte(pattern)[0] != '/' {
	//	err = errors.New("pattern does not start with '/'")
	//	return
	//}

	config := &core.ServerConfig{
		Pattern:      pattern,
		ListenAddr:   listenAddr,
		TLS:          tls,
		CertPath:     certPath,
		KeyPath:      keyPath,
		ReverseProxy: reverseProxy,
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
