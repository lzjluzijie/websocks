package server

import (
	"crypto/tls"
	"net/http"

	"encoding/json"
	"io/ioutil"

	"github.com/gorilla/sessions"
	"gopkg.in/macaron.v1"
)

type App struct {
	//todo multiple servers
	*WebSocksServer

	WebListenAddr string
	TLS           bool
	CertPath      string
	KeyPath       string

	s     http.Server
	store sessions.Store
	m     *macaron.Macaron
}

func LoadApp() (app *App, err error) {
	app = &App{}
	data, err := ioutil.ReadFile("server.json")
	if err != nil {
		return
	}

	err = json.Unmarshal(data, app)
	if err != nil {
		return
	}
	return
}

func NewApp() (app *App) {
	app = &App{
		WebListenAddr: ":23333",
		TLS:           false,
		CertPath:      "websocks.cer",
		KeyPath:       "websocks.key",
	}
	return
}

func (app *App) Save() (err error) {
	data, err := json.MarshalIndent(app, "", "    ")
	if err != nil {
		return
	}

	err = ioutil.WriteFile("server.json", data, 0600)
	return
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
