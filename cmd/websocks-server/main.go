package main

import (
	"flag"
	"log"

	"github.com/lzjluzijie/websocks/core"
)

var l = flag.String("l", ":23333", "listening port")

func main() {
	flag.Parse()

	server := core.Server{
		ListenAddr: *l,
	}
	log.Printf("Listening at %s", *l)
	err := server.Listen()
	if err != nil {
		log.Println(err.Error())
	}
}
