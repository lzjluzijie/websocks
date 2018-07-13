package server

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/gorilla/websocket"
	"github.com/urfave/cli"
)

func GetServer(path string) (server *WebSocksServer, err error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	//read config
	config := &Config{}
	err = json.Unmarshal(data, config)
	if err != nil {
		return
	}

	server = &WebSocksServer{
		Config: config,
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

	config := &Config{
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
