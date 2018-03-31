package main

import (
	"flag"

	"github.com/juju/loggo"
	"github.com/lzjluzijie/websocks/core"
)

var serverAddr string
var pattern string
var debug bool
var logLevel = loggo.INFO

var logger = loggo.GetLogger("server")

func main() {
	flag.StringVar(&serverAddr, "l", ":23333", "server listening port")
	flag.StringVar(&pattern, "p", "/websocks", "server.com/pattern, like password, start with '/'")
	flag.BoolVar(&debug, "debug", false, "debug mode")
	flag.Parse()

	if debug {
		logLevel = loggo.DEBUG
	}

	logger.SetLogLevel(logLevel)
	logger.Infof("Log level %s", logger.LogLevel().String())

	server := core.Server{
		LogLevel:   loggo.DEBUG,
		Pattern:    pattern,
		ListenAddr: serverAddr,
	}
	logger.Infof("Listening at %s", serverAddr)
	err := server.Listen()
	if err != nil {
		logger.Errorf(err.Error())
	}
}
