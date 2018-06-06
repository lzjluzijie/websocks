package config

import (
	"github.com/urfave/cli"
)

var Command = cli.Command{
	Name:    "config",
	Aliases: []string{"config"},
	Usage:   "generate configuration",
	Subcommands: []cli.Command{
		{
			Name:    "client",
			Aliases: []string{"c"},
			Usage:   "generate client config",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "path",
					Value: "client.config.json",
					Usage: "client config output path",
				},
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
					Name:  "n",
					Value: "",
					Usage: "fake server name for tls client hello, leave blank to disable",
				},
				cli.BoolFlag{
					Name:  "insecure",
					Usage: "InsecureSkipVerify: true",
				},
			},
			Action: GenerateClientConfig,
		},
		{
			Name:    "server",
			Aliases: []string{"s"},
			Usage:   "generate server config",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "path",
					Value: "server.config.json",
					Usage: "server config output path",
				},
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
			Action: GenerateServerConfig,
		},
	},
}
