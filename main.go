package main

import (
	"net"
	"net/url"
	"os"

	"os/exec"

	"io/ioutil"

	"github.com/juju/loggo"
	"github.com/lzjluzijie/websocks/core"
	"github.com/urfave/cli"
	"k8s.io/client-go/util/cert"
)

var logger = loggo.GetLogger("websocks")

func main() {
	app := cli.NewApp()
	app.Name = "WebSocks"
	app.Version = "0.2.1"
	app.Usage = "A secure proxy based on websocket."
	app.Description = "See https://github.com/lzjluzijie/websocks"
	app.Author = "Halulu"
	app.Email = "lzjluzijie@gmail.com"

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug",
			Usage: "debug mode",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "client",
			Aliases: []string{"c"},
			Usage:   "start websocks client",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "l",
					Value: ":10801",
					Usage: "local listening port",
				},
				cli.StringFlag{
					Name:  "s",
					Value: "ws://localhost:23333/websocks",
					Usage: "server url",
				},
				cli.BoolFlag{
					Name:  "insecure",
					Usage: "InsecureSkipVerify: true",
				},
			},
			Action: func(c *cli.Context) (err error) {
				debug := c.GlobalBool("debug")
				listenAddr := c.String("l")
				serverURL := c.String("s")
				insecureCert := false

				if c.Bool("insecure") {
					insecureCert = true
				}

				if debug {
					logger.SetLogLevel(loggo.DEBUG)
				} else {
					logger.SetLogLevel(loggo.INFO)
				}

				logger.Infof("Log level %s", logger.LogLevel().String())

				u, err := url.Parse(serverURL)
				if err != nil {
					return
				}

				lAddr, err := net.ResolveTCPAddr("tcp", listenAddr)
				if err != nil {
					return
				}

				local := core.Client{
					LogLevel:     logger.LogLevel(),
					ListenAddr:   lAddr,
					URL:          u,
					InsecureCert: insecureCert,
				}

				err = local.Listen()
				if err != nil {
					return
				}

				return nil
			},
		},
		{
			Name:    "server",
			Aliases: []string{"s"},
			Usage:   "start websocks server",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "l",
					Value: ":23333",
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
			},
			Action: func(c *cli.Context) (err error) {
				debug := c.GlobalBool("debug")
				listenAddr := c.String("l")
				pattern := c.String("p")
				tls := c.Bool("tls")
				certPath := c.String("cert")
				keyPath := c.String("key")

				if debug {
					logger.SetLogLevel(loggo.DEBUG)
				} else {
					logger.SetLogLevel(loggo.INFO)
				}

				logger.Infof("Log level %s", logger.LogLevel().String())

				server := core.Server{
					LogLevel:   logger.LogLevel(),
					Pattern:    pattern,
					ListenAddr: listenAddr,
					TLS:        tls,
					CertPath:   certPath,
					KeyPath:    keyPath,
				}

				logger.Infof("Listening at %s", listenAddr)
				err = server.Listen()
				if err != nil {
					logger.Errorf(err.Error())
				}

				return
			},
		},
		{
			Name:    "github",
			Aliases: []string{"github"},
			Usage:   "open official github page",
			Action: func(c *cli.Context) (err error) {
				err = exec.Command("explorer", "https://github.com/lzjluzijie/websocks").Run()
				return
			},
		},
		{
			Name:    "key",
			Aliases: []string{"key"},
			Usage:   "generate self signed cert and key",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "host",
					Value: "baidu.com",
					Usage: "certificate host",
				},
			},
			Action: func(c *cli.Context) (err error) {
				host := c.String("host")
				certByte, keyByte, err := cert.GenerateSelfSignedCertKey(host, nil, nil)
				if err != nil {
					return err
				}

				err = ioutil.WriteFile("websocks.cer", certByte, 0600)
				if err != nil {
					return err
				}

				err = ioutil.WriteFile("websocks.key", keyByte, 0600)
				if err != nil {
					return err
				}

				return
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		logger.Errorf(err.Error())
	}

}
