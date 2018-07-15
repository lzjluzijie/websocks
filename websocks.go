package main

import (
	"os"

	"encoding/json"
	"errors"
	"io/ioutil"
	"os/exec"
	"runtime"

	"github.com/juju/loggo"
	"github.com/lzjluzijie/websocks/client"
	"github.com/lzjluzijie/websocks/core"
	"github.com/lzjluzijie/websocks/server"
	"github.com/urfave/cli"
)

func main() {
	logger := loggo.GetLogger("websocks")
	logger.SetLogLevel(loggo.INFO)

	app := cli.App{
		Name:        "WebSocks",
		Version:     "0.11.1",
		Usage:       "A secure proxy based on WebSocket. Click to start web client.",
		Description: "websocks.org",
		Author:      "Halulu",
		Email:       "lzjluzijie@gmail.com",
		Action: func(c *cli.Context) (err error) {
			app := client.WebSocksClientApp{}
			app.RunWeb()
			return
		},
		Commands: []cli.Command{
			{
				Name:    "client",
				Aliases: []string{"c"},
				Usage:   "start websocks client(cli)",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "l",
						Value: "127.0.0.1:10801",
						Usage: "local listening port",
					},
					cli.StringFlag{
						Name:  "s",
						Value: "ws://localhost:23333/websocks",
						Usage: "server url",
					},
					cli.BoolFlag{
						Name:  "mux",
						Usage: "mux mode",
					},
					cli.StringFlag{
						Name:  "sni",
						Value: "",
						Usage: "fake server name for tls client hello, leave blank to disable",
					},
					cli.BoolFlag{
						Name:  "insecure",
						Usage: "InsecureSkipVerify: true",
					},
				},
				Action: func(c *cli.Context) (err error) {
					listenAddr := c.String("l")
					serverURL := c.String("s")
					mux := c.Bool("mux")
					sni := c.String("sni")
					insecureCert := false
					if c.Bool("insecure") {
						insecureCert = true
					}

					config := &client.Config{
						ListenAddr:   listenAddr,
						ServerURL:    serverURL,
						SNI:          sni,
						InsecureCert: insecureCert,
						Mux:          mux,
					}

					wc, err := config.GetClient()
					if err != nil {
						return
					}

					err = wc.Run()
					return
				},
			},
			{
				Name: "server",
				//Aliases: []string{"s"},
				Usage: "start websocks server",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "l",
						Value: "0.0.0.0:23333",
						Usage: "local listening port",
					},
					cli.StringFlag{
						Name:  "p",
						Value: "/websocks",
						Usage: "server.com/pattern, like password, start with '/'",
					},
					cli.BoolFlag{
						Name:  "tls",
						Usage: "enable built-in tls",
					},
					cli.StringFlag{
						Name:  "cert",
						Value: "websocks.cer",
						Usage: "tls cert path",
					},
					cli.StringFlag{
						Name:  "key",
						Value: "websocks.key",
						Usage: "tls key path",
					},
					cli.StringFlag{
						Name:  "reverse-proxy",
						Value: "",
						Usage: "reverse proxy url, leave blank to disable",
					},
				},
				Action: func(c *cli.Context) (err error) {
					debug := c.GlobalBool("debug")
					listenAddr := c.String("l")
					pattern := c.String("p")
					tls := c.Bool("tls")
					certPath := c.String("cert")
					keyPath := c.String("key")
					reverseProxy := c.String("reverse-proxy")

					if debug {
						logger.SetLogLevel(loggo.DEBUG)
					}

					logger.Infof("Log level %s", logger.LogLevel().String())

					config := server.Config{
						Pattern:      pattern,
						ListenAddr:   listenAddr,
						TLS:          tls,
						CertPath:     certPath,
						KeyPath:      keyPath,
						ReverseProxy: reverseProxy,
					}

					websocksServer := config.GetServer()
					logger.Infof("Listening at %s", listenAddr)
					err = websocksServer.Run()
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
				Name:    "webserver",
				Aliases: []string{"s"},
				Usage:   "web ui server",
				Action: func(c *cli.Context) (err error) {
					app := &server.App{}
					data, err := ioutil.ReadFile("server.json")
					if err != nil {
						return
					}

					err = json.Unmarshal(data, app)
					if err != nil {
						return
					}

					err = app.Run()
					return
				},
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		logger.Errorf(err.Error())
	}
}
