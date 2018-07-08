package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"

	"github.com/go-macaron/pongo2"
	"github.com/sirupsen/logrus"
	"gopkg.in/macaron.v1"
)

func RunWeb() {
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
	m.Post("/api/client", Client)
	m.Get("/api/log", func(ctx *macaron.Context) {
		ctx.WriteHeader(200)
		ctx.Write(buffer.Bytes())
		return
	})
	//todo pac
	m.Get("/pac", func(ctx *macaron.Context) {
		return
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

func Client(ctx *macaron.Context) {
	webSocksClientConfig := &WebSocksClientConfig{}
	data, err := ioutil.ReadAll(ctx.Req.Body().ReadCloser())
	if err != nil {
		ctx.Error(403, err.Error())
	}

	err = json.Unmarshal(data, webSocksClientConfig)
	if err != nil {
		ctx.Error(403, err.Error())
	}

	websocksClient, err := GetClient(webSocksClientConfig)
	if err != nil {
		ctx.Error(403, err.Error())
	}

	ctx.WriteHeader(200)
	ctx.Write([]byte(fmt.Sprintf("%v", webSocksClientConfig)))

	go func() {
		err = websocksClient.Listen()
		if err != nil {
			log.Error(err.Error())
		}
	}()
	return
}
