package core

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"gopkg.in/macaron.v1"
)

func (server *Server) getMacaron() (m *macaron.Macaron) {
	m = macaron.New()
	m.Group(server.Pattern, func() {
		m.Get("/", server.HandleWebSocket)
		m.Get("/status", server.Status)
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
	}
	return
}
