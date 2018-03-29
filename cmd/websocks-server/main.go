package main

import (
	"flag"

	"github.com/juju/loggo"
	"github.com/lzjluzijie/websocks/core"
)

var serverAddr string
var logLevel = loggo.INFO
var debug bool

var logger = loggo.GetLogger("server")

func main() {
	flag.StringVar(&serverAddr, "l", ":23333", "server listening port")
	flag.BoolVar(&debug, "debug", false, "debug mode")
	if debug {
		logLevel = loggo.DEBUG
	}
	flag.Parse()

	logger.SetLogLevel(logLevel)
	logger.Infof("Log level %s", logger.LogLevel().String())

	server := core.Server{
		LogLevel:   loggo.DEBUG,
		ListenAddr: serverAddr,
	}
	logger.Infof("Listening at %s", serverAddr)
	err := server.Listen()
	if err != nil {
		logger.Errorf(err.Error())
	}
}
