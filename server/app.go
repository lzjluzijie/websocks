package server

import (
	"crypto/tls"
	"net/http"

	"gopkg.in/macaron.v1"
)

type App struct {
	//todo multiple servers
	*WebSocksServer

	WebListenAddr string
	TLS           bool
	CertPath      string
	KeyPath       string

	s http.Server
	m *macaron.Macaron
}

func (app *App) Run() (err error) {
	m := app.Macaron()
	app.m = m
	app.s = http.Server{
		Addr:    app.WebListenAddr,
		Handler: m,
	}

	if !app.TLS {
		err = app.s.ListenAndServe()
		if err != nil {
			return
		}
		return
	}

	app.s.TLSConfig = &tls.Config{
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		},
	}

	err = app.s.ListenAndServeTLS(app.CertPath, app.KeyPath)
	if err != nil {
		return err
	}
	return
}
