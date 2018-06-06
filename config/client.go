package config

import (
	"crypto/tls"
	"net"
	"net/url"

	"github.com/juju/loggo"
	"github.com/urfave/cli"
)

type ClientConfig struct {
	LogLevel   loggo.Level
	ListenAddr *net.TCPAddr
	URL        *url.URL
	TLSConfig  *tls.Config
	Mux        bool
}

//GetClientConfig create a client config from cli.Context
func GetClientConfig(c *cli.Context) (config *ClientConfig, err error) {
	debug := c.GlobalBool("debug")
	listenAddr := c.String("l")
	serverURL := c.String("s")
	mux := c.Bool("mux")
	serverName := c.String("n")
	insecureCert := false
	if c.Bool("insecure") {
		insecureCert = true
	}

	if debug {
		logger.SetLogLevel(loggo.DEBUG)
	}

	logger.Infof("Log level %s", logger.LogLevel().String())

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

	config = &ClientConfig{
		LogLevel:   logger.LogLevel(),
		ListenAddr: lAddr,
		URL:        u,
		Mux:        mux,
	}
	return
}
