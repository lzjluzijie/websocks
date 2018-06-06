package core

import (
	"io"
	"net"
	"net/http"
	"time"

	"fmt"

	"crypto/tls"

	"sync"

	"github.com/gorilla/websocket"
	"github.com/juju/loggo"
)

type ServerConfig struct {
	ListenAddr   string
	Pattern      string
	TLS          bool
	CertPath     string
	KeyPath      string
	ReverseProxy string
}

type Server struct {
	*ServerConfig
	LogLevel loggo.Level

	Upgrader   *websocket.Upgrader
	muxConnMap sync.Map
	mutex      sync.Mutex

	//statistics
	CreatedAt       time.Time
	statMutex       sync.Mutex
	openedConn      uint64
	closedConn      uint64
	downloadedBytes uint64
	uploadedBytes   uint64
}

func (server *Server) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	c, err := server.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Debugf(err.Error())
		return
	}
	defer c.Close()

	ws := &WebSocket{
		conn: c,
	}

	server.openConn()
	defer server.closeConn()

	//mux
	if r.Header.Get("WebSocks-Mux") == "mux" {
		muxWS := NewMuxWebSocket(ws)
		muxWS.ServerListen()
		return
	}

	host := r.Header.Get("WebSocks-Host")
	logger.Debugf("Dial %s", host)

	conn, err := net.Dial("tcp", host)
	if err != nil {
		logger.Debugf(err.Error())
		return
	}
	defer conn.Close()

	go func() {
		downloaded, err := io.Copy(conn, ws)
		server.downloaded(uint64(downloaded))
		if err != nil {
			logger.Debugf(err.Error())
			return
		}
	}()

	uploaded, err := io.Copy(ws, conn)
	server.uploaded(uint64(uploaded))
	if err != nil {
		logger.Debugf(err.Error())
		return
	}

	return
}

func (server *Server) openConn() {
	server.statMutex.Lock()
	server.openedConn++
	server.statMutex.Unlock()
}

func (server *Server) closeConn() {
	server.statMutex.Lock()
	server.closedConn++
	server.statMutex.Unlock()
}

func (server *Server) downloaded(d uint64) {
	server.statMutex.Lock()
	server.downloadedBytes += d
	server.statMutex.Unlock()
}

func (server *Server) uploaded(u uint64) {
	server.statMutex.Lock()
	server.uploadedBytes += u
	server.statMutex.Unlock()
}

func (server *Server) status() string {
	return fmt.Sprintf("%ds: opened %d, closed %d, uploaded %d bytes, downloaded %d bytes", int(time.Since(server.CreatedAt).Seconds()), server.openedConn, server.closedConn, server.uploadedBytes, server.downloadedBytes)
}

func (server *Server) Listen() (err error) {
	logger.SetLogLevel(server.LogLevel)

	go func() {
		for {
			time.Sleep(time.Second)
			logger.Debugf("%ds: opened %d, closed %d, uploaded %d bytes, downloaded %d bytes", int(time.Since(server.CreatedAt).Seconds()), server.openedConn, server.closedConn, server.uploadedBytes, server.downloadedBytes)
		}
	}()

	s := http.Server{
		Addr:    server.ListenAddr,
		Handler: server.getMacaron(),
	}

	logger.Infof("Start to listen at %s", server.ListenAddr)

	if !server.TLS {
		err = s.ListenAndServe()
		if err != nil {
			return err
		}
		return
	} else {
		s.TLSConfig = &tls.Config{
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			},
		}

		err = s.ListenAndServeTLS(server.CertPath, server.KeyPath)
		if err != nil {
			return err
		}
	}

	return
}
