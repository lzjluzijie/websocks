package core

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/go-macaron/pongo2"
	"gopkg.in/macaron.v1"
)

func (server *Server) getMacaron() (m *macaron.Macaron) {
	m = macaron.New()
	m.Use(pongo2.Pongoer())
	m.Group(server.Pattern, func() {
		m.Get("/", server.HandleWebSocket)
		m.Get("/status", server.getStatus)
	})

	if server.Proxy != "" {
		m.NotFound(func(w http.ResponseWriter, r *http.Request) {
			remote, err := url.Parse(server.Proxy)
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

func (server *Server) getStatus(ctx *macaron.Context) {
	ctx.Data["Status"] = server.status()
	ctx.HTML(200, "status")
	return
}
