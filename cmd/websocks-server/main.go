package main

import (
	"flag"

	"github.com/juju/loggo"
	"github.com/lzjluzijie/websocks/core"
)

var l = flag.String("l", ":23333", "server listening port")

var logger = loggo.GetLogger("server")

func main() {
	flag.Parse()

	server := core.Server{
		LogLevel:   loggo.DEBUG,
		ListenAddr: *l,
	}
	logger.Infof("Listening at %s", *l)
	err := server.Listen()
	if err != nil {
		logger.Errorf(err.Error())
	}
}
