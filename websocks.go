package main

import (
	"os"

	"io/ioutil"

	"errors"
	"runtime"

	"os/exec"

	"net/http"

	"github.com/go-macaron/pongo2"
	"github.com/juju/loggo"
	"github.com/lzjluzijie/websocks/config"
	"github.com/lzjluzijie/websocks/core"
	"github.com/lzjluzijie/websocks/core/client"
	"github.com/urfave/cli"
	"gopkg.in/macaron.v1"
)

func main() {
	logger := loggo.GetLogger("websocks")
	logger.SetLogLevel(loggo.INFO)

	app := cli.NewApp()
	app.Name = "WebSocks"
	app.Version = "0.9.2"
	app.Usage = "A secure proxy based on WebSocket."
	app.Description = "See websocks.org"
	app.Author = "Halulu"
	app.Email = "lzjluzijie@gmail.com"

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug",
			Usage: "debug mode",
		},
	}

	app.Commands = []cli.Command{
		config.Command,
		{
			Name:    "client",
			Aliases: []string{"c"},
			Usage:   "start websocks client",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "c",
					Value: "client.config.json",
					Usage: "client config path",
				},
			},
			Action: func(c *cli.Context) (err error) {
				path := c.String("c")
				debug := c.GlobalBool("debug")

				data, err := ioutil.ReadFile(path)
				if err != nil {
					return
				}

				webSocksClient, err := client.GetClientConfig(data)
				if err != nil {
					return
				}

				logLevel := loggo.INFO
				if debug {
					logLevel = loggo.DEBUG
				}
				webSocksClient.LogLevel = logLevel

				err = webSocksClient.Listen()
				if err != nil {
					return
				}
				return
			},
		},
		{
			Name:    "server",
			Aliases: []string{"s"},
			Usage:   "start websocks server",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "c",
					Value: "server.config.json",
					Usage: "server config path",
				},
			},
			Action: func(c *cli.Context) (err error) {
				path := c.String("c")
				debug := c.GlobalBool("debug")

				server, err := config.GetServerConfig(path)
				if err != nil {
					return
				}

				logLevel := loggo.INFO
				if debug {
					logLevel = loggo.DEBUG
				}
				server.LogLevel = logLevel

				err = server.Listen()
				if err != nil {
					return
				}
				return
			},
		},
		{
			Name:    "cert",
			Aliases: []string{"cert"},
			Usage:   "generate self signed key and cert(default rsa 2048)",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "ecdsa",
					Usage: "generate ecdsa key and cert(P-256)",
				},
				cli.StringSliceFlag{
					Name:  "hosts",
					Value: nil,
					Usage: "certificate hosts",
				},
			},
			Action: func(c *cli.Context) (err error) {
				ecdsa := c.Bool("ecdsa")
				hosts := c.StringSlice("hosts")

				var key, cert []byte
				if ecdsa {
					key, cert, err = core.GenP256(hosts)
					logger.Infof("Generated ecdsa P-256 key and cert")
				} else {
					key, cert, err = core.GenRSA2048(hosts)
					logger.Infof("Generated rsa 2048 key and cert")
				}

				err = ioutil.WriteFile("websocks.key", key, 0600)
				if err != nil {
					return
				}
				err = ioutil.WriteFile("websocks.cer", cert, 0600)
				if err != nil {
					return
				}
				return
			},
		},
		{
			Name:    "pac",
			Aliases: []string{"pac"},
			Usage:   "set pac for windows",
			Action: func(c *cli.Context) (err error) {
				if runtime.GOOS != "windows" {
					err = errors.New("not windows")
					return
				}

				err = exec.Command("REG", "ADD", `HKCU\Software\Microsoft\Windows\CurrentVersion\Internet Settings`, "/v", "AutoConfigURL", "/d", "http://127.0.0.1:10801/pac", "/f").Run()
				return
			},
		},
		{
			Name:    "webclient",
			Aliases: []string{"wc"},
			Usage:   "test webui client",
			Action: func(c *cli.Context) (err error) {
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

				err = http.ListenAndServe(":10801", m)
				return
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		logger.Errorf(err.Error())
	}
}
