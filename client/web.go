package client

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"

	"github.com/go-macaron/pongo2"
	"github.com/sirupsen/logrus"
	"gopkg.in/macaron.v1"
)

func (app *WebSocksClientApp) RunWeb() {
	//log setup
	buf := make([]byte, 0)
	buffer := bytes.NewBuffer(buf)
	log.Out = io.MultiWriter(os.Stdout, buffer)
	log.SetLevel(logrus.DebugLevel)

	m := macaron.New()
	m.Use(pongo2.Pongoer())
	m.Get("/", func(ctx *macaron.Context) {
		ctx.HTML(200, "client")
		return
	})

	//todo pac
	m.Get("/pac", func(ctx *macaron.Context) {
		return
	})

	//api v0
	m.Group("/api/v0/client", func() {
		m.Get("/log", func(ctx *macaron.Context) {
			ctx.WriteHeader(200)
			ctx.Write(buffer.Bytes())
			return
		})
		m.Get("/stats", func(ctx *macaron.Context) {
			if app.WebSocksClient == nil {
				ctx.Error(403, "websocks client is not running")
				return
			}
			ctx.JSON(200, app.WebSocksClient.Stats)
		})
		m.Post("/start", app.StartClient)
		m.Post("/stop", app.StopClient)
	})

	go func() {
		err := exec.Command("explorer", "http://127.0.0.1:10801").Run()
		if err != nil {
			log.Debug(err.Error())
			return
		}
	}()

	log.Infof("web start to listen at :10801")
	err := http.ListenAndServe(":10801", m)
	if err != nil {
		log.Error(err.Error())
		return
	}
	return
}

func (app *WebSocksClientApp) StartClient(ctx *macaron.Context) {
	config := &Config{}
	data, err := ioutil.ReadAll(ctx.Req.Body().ReadCloser())
	if err != nil {
		ctx.Error(403, err.Error())
	}

	err = json.Unmarshal(data, config)
	if err != nil {
		ctx.Error(403, err.Error())
	}

	websocksClient, err := config.GetClient()
	if err != nil {
		ctx.Error(403, err.Error())
	}

	app.WebSocksClient = websocksClient
	app.running = true

	go func() {
		err = websocksClient.Listen()
		if err != nil {
			log.Error(err.Error())
		}
	}()
	return
}

func (app *WebSocksClientApp) StopClient(ctx *macaron.Context) {
	app.WebSocksClient.Stop()
	ctx.WriteHeader(200)
	ctx.Write([]byte("stopped"))
	return
}
