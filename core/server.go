package core

import (
	"io"
	"net"
	"net/http"
	"sync/atomic"
	"time"

	"fmt"

	"crypto/tls"

	"sync"

	"github.com/gorilla/websocket"
	"github.com/juju/loggo"
)

type Server struct {
	LogLevel   loggo.Level
	Pattern    string
	ListenAddr string
	TLS        bool
	CertPath   string
	KeyPath    string
	Proxy      string

	Upgrader *websocket.Upgrader

	MessageChan chan *Message
	muxConnMap  sync.Map
	Mutex       sync.Mutex

	CreatedAt time.Time

	Opened     uint64
	Closed     uint64
	Uploaded   uint64
	Downloaded uint64
}

func (server *Server) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	c, err := server.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Debugf(err.Error())
		return
	}

	ws := &WebSocket{
		conn: c,
	}

	atomic.AddUint64(&server.Opened, 1)
	defer atomic.AddUint64(&server.Closed, 1)

	if r.Header.Get("WebSocks-Mux") == "mux" {
		muxWS := NewMuxWebSocket(ws)
		muxWS.ServerListen()
		return
	}

	host := r.Header.Get("WebSocks-Host")
	logger.Debugf("Dial %s", host)

	conn, err := net.Dial("tcp", host)
	if err != nil {
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
			logger.Debugf(err.Error())
			return
		}
	}()

	uploaded, err := io.Copy(ws, conn)
	atomic.AddUint64(&server.Uploaded, uint64(uploaded))
	if err != nil {
		logger.Debugf(err.Error())
		return
	}
	return
}

func (server *Server) status() string {
	return fmt.Sprintf("%ds: opened %d, closed %d, uploaded %d bytes, downloaded %d bytes", int(time.Since(server.CreatedAt).Seconds()), server.Opened, server.Closed, server.Uploaded, server.Downloaded)
}

func (server *Server) Listen() (err error) {
	logger.SetLogLevel(server.LogLevel)

	go func() {
		for {
			time.Sleep(time.Second)
			logger.Debugf("%ds: opened %d, closed %d, uploaded %d bytes, downloaded %d bytes", int(time.Since(server.CreatedAt).Seconds()), server.Opened, server.Closed, server.Uploaded, server.Downloaded)
		}
	}()

	s := http.Server{
		Addr:    server.ListenAddr,
		Handler: server.getMacaron(),
	}

	if !server.TLS {
		err = s.ListenAndServe()
		if err != nil {
			return err
		}
		return
	} else {
		tlsConfig := &tls.Config{
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			},
		}

		s.TLSConfig = tlsConfig
		err = s.ListenAndServeTLS(server.CertPath, server.KeyPath)
		if err != nil {
			return err
		}
	}

	return
}
