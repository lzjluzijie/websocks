package server

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"encoding/json"
	"io/ioutil"
	"log"
	"os/exec"

	"github.com/go-macaron/pongo2"
	"github.com/gorilla/sessions"
	"gopkg.in/macaron.v1"
)

func (app *App) Macaron() (m *macaron.Macaron) {
	app.store = sessions.NewCookieStore([]byte("just a test"))

	m = macaron.New()
	m.Use(pongo2.Pongoer())

	m.Get("/", func(ctx *macaron.Context) {
		session, _ := app.store.Get(ctx.Req.Request, "cookie")
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			ctx.Error(403, "not halulu")
			return
		}

		ctx.HTML(200, "server/server")
	})

	m.Get("/login", func(ctx *macaron.Context) {
		ctx.HTML(200, "server/login")
	})

	//api v0
	m.Group("/api/v0/server", func() {
		m.Get("/stats", func(ctx *macaron.Context) {
			if app.WebSocksServer == nil {
				ctx.Error(403, "websocks server is not running")
				return
			}

			stats := app.Stats
			ctx.JSON(200, stats)
		})
		m.Post("/start", app.StartServer)
		//m.Post("/stop", app.StopServer)

		m.Post("/login", app.Login)
	})

	go func() {
		err := exec.Command("explorer", "http://127.0.0.1:23333").Run()
		if err != nil {
			log.Println(err.Error())
			return
		}
	}()

	return
}

func (app *App) StartServer(ctx *macaron.Context) {
	config := &Config{}
	data, err := ioutil.ReadAll(ctx.Req.Body().ReadCloser())
	if err != nil {
		ctx.Error(403, err.Error())
	}

	err = json.Unmarshal(data, config)
	if err != nil {
		ctx.Error(403, err.Error())
	}
	ctx.JSON(200, config)

	webSocksServer := config.NewWebSocksServer()
	app.WebSocksServer = webSocksServer
	app.m.Get(webSocksServer.Pattern, webSocksServer.HandleWebSocket)
	return
}

func (server *WebSocksServer) getMacaron() (m *macaron.Macaron) {
	m = macaron.New()
	m.Use(pongo2.Pongoer())
	m.Group(server.Pattern, func() {
		m.Get("/", server.HandleWebSocket)
	})

	if server.ReverseProxy != "" {
		m.NotFound(func(w http.ResponseWriter, r *http.Request) {
			remote, err := url.Parse(server.ReverseProxy)
			if err != nil {
				panic(err)
			}
			proxy := httputil.NewSingleHostReverseProxy(remote)
			proxy.ServeHTTP(w, r)
		})
	} else {
		m.Get("/", func(ctx *macaron.Context) {
			ctx.HTML(200, "home")
		})
	}
	return
}
