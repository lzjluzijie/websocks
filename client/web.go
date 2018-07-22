package client

import (
	"encoding/json"
	"io/ioutil"

	"gopkg.in/macaron.v1"
)

func (app *App) StartClient(ctx *macaron.Context) {
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
		err = websocksClient.Run()
		if err != nil {
			log.Error(err.Error())
		}
	}()
	return
}

func (app *App) StopClient(ctx *macaron.Context) {
	app.WebSocksClient.Stop()
	ctx.WriteHeader(200)
	ctx.Write([]byte("stopped"))
	return
}
