package main

import (
	"flag"
	"net"

	"github.com/juju/loggo"
	"github.com/lzjluzijie/websocks/core"
	"net/url"
)

var l = flag.String("l", ":10801", "local listening port")
var sURL = flag.String("u", "ws://localhost:23333/ws", "server url")

var logger = loggo.GetLogger("local")

func main() {
	logger.SetLogLevel(loggo.DEBUG)

	flag.Parse()

	u, err := url.Parse(*sURL)
	if err != nil {
		logger.Errorf(err.Error())
		return
	}

	laddr, err := net.ResolveTCPAddr("tcp", *l)
	if err != nil {
		logger.Errorf(err.Error())
	}

	local := core.Local{
		LogLevel:   loggo.DEBUG,
		ListenAddr: laddr,
		URL:        u,
	}

	err = local.Listen()
	if err != nil {
		logger.Errorf(err.Error())
	}

}
