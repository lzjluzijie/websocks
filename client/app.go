package client

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"

	"encoding/json"
	"io/ioutil"

	"github.com/go-macaron/pongo2"
	"github.com/sirupsen/logrus"
	"gopkg.in/macaron.v1"
)

type App struct {
	//todo
	WebListenAddr string

	m macaron.Macaron

	running bool
	//todo multiple client
	*WebSocksClient
}

func LoadApp() (app *App, err error) {
	app = &App{}
	data, err := ioutil.ReadFile("client.json")
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
		WebListenAddr: ":10801",
	}
	return
}

func (app *App) Save() (err error) {
	data, err := json.MarshalIndent(app, "", "    ")
	if err != nil {
		return
	}

	err = ioutil.WriteFile("client.json", data, 0600)
	return
}

func (app *App) Run() (err error) {
	//log setup
	buf := make([]byte, 0)
	buffer := bytes.NewBuffer(buf)
	log.Out = io.MultiWriter(os.Stdout, buffer)
	log.SetLevel(logrus.DebugLevel)

	m := macaron.New()
	m.Use(pongo2.Pongoer())
	m.Get("/", func(ctx *macaron.Context) {
		ctx.HTML(200, "client/client")
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
		err := exec.Command("explorer", fmt.Sprintf("http://127.0.0.1%s", app.WebListenAddr)).Run()
		if err != nil {
			log.Debug(err.Error())
			return
		}
	}()

	log.Infof("web start to listen at %s", app.WebListenAddr)
	err = http.ListenAndServe(app.WebListenAddr, m)
	if err != nil {
		log.Error(err.Error())
		return
	}
	return
}
