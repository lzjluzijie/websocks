package main

import (
	"os"

	"github.com/juju/loggo"
	"github.com/lzjluzijie/websocks/core"
	"github.com/urfave/cli"
)

var logger = loggo.GetLogger("server")

func main() {
	app := cli.NewApp()
	app.Name = "WebSocks Server"
	app.Version = "0.1.1"
	app.Author = "Halulu"
	app.Usage = "See https://github.com/lzjluzijie/websocks"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "l",
			Value: ":23333",
			Usage: "server listening port",
		},
		cli.StringFlag{
			Name:  "p",
			Value: "/websocks",
			Usage: "server.com/pattern, like password, start with '/'",
		},
		cli.BoolFlag{
			Name:  "debug",
			Usage: "debug mode",
		},
		cli.BoolFlag{
			Name:  "tls",
			Usage: "enable built-in tls",
		},
	}

	app.Action = func(c *cli.Context) (err error) {
		debug := c.Bool("debug")
		pattern := c.String("p")
		listenAddr := c.String("l")
		tls := c.Bool("tls")

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
		}

		logger.Infof("Listening at %s", listenAddr)
		err = server.Listen()
		if err != nil {
			logger.Errorf(err.Error())
		}
		return
	}

	app.Run(os.Args)
}
