package server

import (
	"io"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/lzjluzijie/websocks/core/mux"

	"net/http/httputil"
	"net/url"

	"crypto/tls"

	"github.com/gorilla/websocket"
	"github.com/juju/loggo"
	"github.com/julienschmidt/httprouter"
	"github.com/lzjluzijie/websocks/core"
)

type WebSocksServer struct {
	*Config
	LogLevel loggo.Level

	Upgrader *websocket.Upgrader

	//todo multiple clients
	group *mux.Group

	CreatedAt time.Time
	Stats     *core.Stats
}

func (server *WebSocksServer) HandleWebSocket(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	wsConn, err := server.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf(err.Error())
		return
	}
	defer wsConn.Close()

	ws := core.NewWebSocket(wsConn, server.Stats)
	//todo conns

	//mux
	//todo multiple clients
	if r.Header.Get("WebSocks-Mux") == "v0.15" {
		muxWS := mux.NewMuxWebSocket(ws)
		server.group.AddMuxWS(muxWS)
		return
	}

	host := r.Header.Get("WebSocks-Host")
	log.Printf("Dial %s", host)
	conn, err := server.DialRemote(host)
	if err != nil {
		log.Printf(err.Error())
		return
	}

	go func() {
		_, err = io.Copy(conn, ws)
		if err != nil {
			log.Printf(err.Error())
			return
		}
	}()

	_, err = io.Copy(ws, conn)
	if err != nil {
		log.Printf(err.Error())
		return
	}

	return
}

func (server *WebSocksServer) DialRemote(host string) (conn net.Conn, err error) {
	conn, err = net.Dial("tcp", host)
	if err != nil {
		return
	}
	return
}

func (server *WebSocksServer) Run() (err error) {
	r := httprouter.New()
	r.GET(server.Pattern, server.HandleWebSocket)

	if server.ReverseProxy != "" {
		remote, err := url.Parse(server.ReverseProxy)
		if err != nil {
			panic(err)
		}
		proxy := httputil.NewSingleHostReverseProxy(remote)
		r.NotFound = proxy
	}

	s := http.Server{
		Addr:    server.ListenAddr,
		Handler: r,
	}

	log.Printf("Start to listen at %s", server.ListenAddr)

	if !server.TLS {
		err = s.ListenAndServe()
		if err != nil {
			return err
		}
		return
	}

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

	err = s.ListenAndServeTLS(server.Config.CertPath, server.Config.KeyPath)
	if err != nil {
		return err
	}
	return
}
