package config

import (
	"github.com/juju/loggo"
	"github.com/urfave/cli"
)

type ServerConfig struct {
	LogLevel   loggo.Level
	ListenAddr string
	Pattern    string
	TLS        bool
	CertPath   string
	KeyPath    string
	Proxy      string
}

//GetServerConfig create a server config from cli.Context
func GetServerConfig(c *cli.Context) (config *ServerConfig, err error) {
	debug := c.GlobalBool("debug")
	listenAddr := c.String("l")
	pattern := c.String("p")
	tls := c.Bool("tls")
	certPath := c.String("cert")
	keyPath := c.String("key")
	proxy := c.String("proxy")

	if debug {
		logger.SetLogLevel(loggo.DEBUG)
	}

	logger.Infof("Log level %s", logger.LogLevel().String())

	config = &ServerConfig{
		LogLevel:   logger.LogLevel(),
		Pattern:    pattern,
		ListenAddr: listenAddr,
		TLS:        tls,
		CertPath:   certPath,
		KeyPath:    keyPath,
		Proxy:      proxy,
	}
	return
}
