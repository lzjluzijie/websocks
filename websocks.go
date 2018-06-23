package main

import (
	"os"

	"io/ioutil"

	"errors"
	"runtime"

	"github.com/juju/loggo"
	"github.com/lzjluzijie/websocks/config"
	"github.com/lzjluzijie/websocks/core"
	"github.com/urfave/cli"
	"golang.org/x/sys/windows/registry"
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

				client, err := config.GetClientConfig(path)
				if err != nil {
					return
				}

				logLevel := loggo.INFO
				if debug {
					logLevel = loggo.DEBUG
				}
				client.LogLevel = logLevel

				err = client.Listen()
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

				k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Internet Settings`, registry.ALL_ACCESS)
				if err != nil {
					return
				}

				err = k.SetStringValue("AutoConfigURL", "http://127.0.0.1:10801/pac")
				return
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		logger.Errorf(err.Error())
	}
}
