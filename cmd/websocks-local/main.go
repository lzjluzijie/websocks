package main

import (
	"flag"
	"log"
	"net"

	"github.com/lzjluzijie/websocks/core"
)

var l = flag.String("l", ":10801", "listening port")
var url = flag.String("u", "ws://localhost:23333/ws", "url")
var origin = flag.String("o", "http://localhost/", "origin")

func main() {
	flag.Parse()

	laddr, err := net.ResolveTCPAddr("tcp", *l)
	if err != nil {
		log.Println(err.Error())
	}

	local := core.Local{
		ListenAddr: laddr,
		URL:        *url,
		Origin:     *origin,
	}

	err = local.Listen()
	if err != nil {
		log.Println(err.Error())
	}
}
