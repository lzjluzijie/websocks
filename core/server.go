package core

import (
	"encoding/gob"
	"io"
	"net"
	"net/http"

	"time"

	"net/http/httputil"
	"net/url"

	"sync/atomic"

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

	CreatedAt time.Time

	Opened     uint64
	Closed     uint64
	Uploaded   uint64
	Downloaded uint64
}

func (server *Server) HandleWebSocket(ws *websocket.Conn) {
	server.Opened++

	defer ws.Close()

	dec := gob.NewDecoder(ws)
	req := Request{}
	err := dec.Decode(&req)
	if err != nil {
		server.Closed++
		if err != nil {
			logger.Debugf(err.Error())
		}
		return
	}

	logger.Debugf("Dial %s from %s", req.Addr, ws.RemoteAddr().String())
	conn, err := net.Dial("tcp", req.Addr)
	if err != nil {
		server.Closed++
		if err != nil {
			logger.Debugf(err.Error())
		}
		return
	}

	defer conn.Close()

	go func() {
		downloaded, err := io.Copy(conn, ws)
		atomic.AddUint64(&server.Downloaded, uint64(downloaded))
		if err != nil {
			server.Closed++
			if err != nil {
				logger.Debugf(err.Error())
			}
			return
		}
	}()

	uploaded, err := io.Copy(ws, conn)
	atomic.AddUint64(&server.Uploaded, uint64(uploaded))
	if err != nil {
		server.Closed++
		if err != nil {
			logger.Debugf(err.Error())
		}
		return
	}
}

func (server *Server) Listen() (err error) {
	logger.SetLogLevel(server.LogLevel)

	go func() {
		for {
			time.Sleep(time.Second)
			logger.Debugf("%ds: opened %d, closed %d, uploaded %d bytes, downloaded %d bytes", int(time.Since(server.CreatedAt).Seconds()), server.Opened, server.Closed, server.Uploaded, server.Downloaded)
		}
	}()

	mux := http.NewServeMux()
	mux.Handle(server.Pattern, websocket.Handler(server.HandleWebSocket))
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
