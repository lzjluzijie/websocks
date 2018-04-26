package main

import (
	"net"

	"net/url"

	"os"

	"github.com/juju/loggo"
	"github.com/lzjluzijie/websocks/core"
	"github.com/urfave/cli"
)

var logger = loggo.GetLogger("local")

func main() {
	app := cli.NewApp()
	app.Name = "WebSocks Server"
	app.Version = "0.1.1"
	app.Author = "Halulu"
	app.Usage = "See https://github.com/lzjluzijie/websocks"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "u",
			Value: "ws://localhost:23333/websocks",
			Usage: "server url",
		},
		cli.StringFlag{
			Name:  "l",
			Value: ":10801",
			Usage: "local listening port",
		},
		cli.BoolFlag{
			Name:  "debug",
			Usage: "debug mode",
		},
	}

	app.Action = func(c *cli.Context) (err error) {
		debug := c.Bool("debug")
		serverURL := c.String("u")
		localAddr := c.String("l")

		if debug {
			logger.SetLogLevel(loggo.DEBUG)
		} else {
			logger.SetLogLevel(loggo.INFO)
		}

		logger.Infof("Log level %s", logger.LogLevel().String())

		u, err := url.Parse(serverURL)
		if err != nil {
			logger.Errorf(err.Error())
			return
		}

		lAddr, err := net.ResolveTCPAddr("tcp", localAddr)
		if err != nil {
			logger.Errorf(err.Error())
		}

		local := core.Local{
			LogLevel:   logger.LogLevel(),
			ListenAddr: lAddr,
			URL:        u,
		}

		err = local.Listen()
		if err != nil {
			logger.Errorf(err.Error())
		}
		return
	}

	app.Run(os.Args)

}
