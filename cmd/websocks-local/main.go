package main

import (
	"flag"
	"net"

	"github.com/juju/loggo"
	"github.com/lzjluzijie/websocks/core"
)

var l = flag.String("l", ":10801", "local listening port")
var url = flag.String("u", "ws://localhost:23333/ws", "server url")
var origin = flag.String("o", "http://localhost/", "server origin")

var logger = loggo.GetLogger("local")

func main() {
	flag.Parse()

	laddr, err := net.ResolveTCPAddr("tcp", *l)
	if err != nil {
		logger.Errorf(err.Error())
	}

	local := core.Local{
		LogLevel:   loggo.DEBUG,
		ListenAddr: laddr,
		URL:        *url,
		Origin:     *origin,
	}

	err = local.Listen()
	if err != nil {
		logger.Errorf(err.Error())
	}

}
