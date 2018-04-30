package core

import (
	"encoding/gob"
	"io"
	"net"
	"net/http"

	"time"

	"net/http/httputil"
	"net/url"

	"github.com/juju/loggo"
	"golang.org/x/net/websocket"
)

type Server struct {
	LogLevel   loggo.Level
	Pattern    string
	ListenAddr string
	TLS        bool
	CertPath   string
	KeyPath    string
	Proxy      string
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
			logger.Debugf("%ds: opened%d, closed%d", int(time.Since(t).Seconds()), opened, closed)
		}
	}()

	mux := http.NewServeMux()
	mux.Handle(server.Pattern, websocket.Handler(handler))
	if server.Proxy != "" {
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			remote, err := url.Parse(server.Proxy)
			if err != nil {
				panic(err)
			}
			proxy := httputil.NewSingleHostReverseProxy(remote)
			proxy.ServeHTTP(w, r)
		})
	}

	s := http.Server{
		Addr:    server.ListenAddr,
		Handler: mux,
	}

	if !server.TLS {
		err = s.ListenAndServe()
		if err != nil {
			return err
		}
		return
	} else {
		err = s.ListenAndServeTLS(server.CertPath, server.KeyPath)
		if err != nil {
			return err
		}
	}

	return
}
