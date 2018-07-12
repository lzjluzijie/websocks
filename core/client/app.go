package client

import "gopkg.in/macaron.v1"

type WebSocksClientApp struct {
	//todo
	WebListenAddr string

	m macaron.Macaron

	running bool
	//todo multiple client
	*WebSocksClient
}

func (app *WebSocksClientApp) GetStatus() (stats *Stats) {
	if !app.running {
		return nil
	}

	stats = app.WebSocksClient.Status()
	return
}
