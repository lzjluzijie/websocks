package core

import (
	"encoding/gob"
	"io"
	"net"
	"net/http"

	"time"

	"crypto/tls"

	"github.com/juju/loggo"
	"golang.org/x/net/websocket"
	"k8s.io/client-go/util/cert"
)

type Server struct {
	LogLevel   loggo.Level
	Pattern    string
	ListenAddr string
	TLS        bool
}

var opened = 0
var closed = 0
var t time.Time

func handler(ws *websocket.Conn) {
	opened++
	var err error
	defer func() {
		closed++
		if err != nil {
			logger.Debugf(err.Error())
		}
	}()

	defer ws.Close()

	dec := gob.NewDecoder(ws)
	req := Request{}
	err = dec.Decode(&req)
	if err != nil {
		return
	}

	logger.Debugf("Dial %s from %s", req.Addr, ws.RemoteAddr().String())
	conn, err := net.Dial("tcp", req.Addr)
	if err != nil {
		return
	}

	defer conn.Close()

	go func() {
		_, err = io.Copy(conn, ws)
		if err != nil {
			logger.Debugf(err.Error())
			return
		}
	}()

	_, err = io.Copy(ws, conn)
	if err != nil {
		return
	}
}

func (server *Server) Listen() (err error) {
	logger.SetLogLevel(server.LogLevel)

	t = time.Now()

	go func() {
		for {
			time.Sleep(time.Second)
			logger.Debugf("%s: opened%d, closed%d", time.Since(t), opened, closed)
		}
	}()

	if !server.TLS {
		http.Handle(server.Pattern, websocket.Handler(handler))
		err = http.ListenAndServe(server.ListenAddr, nil)
		if err != nil {
			return err
		}
		return
	}

	println("tls")
	c, k, err := cert.GenerateSelfSignedCertKey("baidu.com", nil, nil)
	if err != nil {
		return err
	}

	certificate, err := tls.X509KeyPair(c, k)
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle(server.Pattern, websocket.Handler(handler))

	s := http.Server{
		Addr: server.ListenAddr,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{certificate},
		},
		Handler: mux,
	}

	err = s.ListenAndServeTLS("", "")
	if err != nil {
		return err
	}
	return
}
