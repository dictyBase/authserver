package main

import (
	"os"

	"github.com/dictyBase/authserver/commands"
	"github.com/dictyBase/authserver/validate"
	"gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()
	app.Name = "authserver"
	app.Usage = "oauth server that provides endpoints for managing authentication"
	app.Version = "3.0.0"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "log,l",
			Usage: "Name of the log file(optional), default goes to stderr",
		},
		cli.StringFlag{
			Name:  "log-format",
			Usage: "Format of the log output,could be either of text or json, default is json",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:   "run",
			Usage:  "runs the auth server",
			Action: commands.RunServer,
			Before: validate.ValidateRunArgs,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "config, c",
					Usage:  "Config file(required)",
					EnvVar: "OAUTH_CONFIG",
				},
				cli.StringFlag{
					Name:   "pkey, public-key",
					Usage:  "public key file for verifying jwt",
					EnvVar: "JWT_PUBLIC_KEY",
				},
				cli.StringFlag{
					Name:   "private-key, prkey",
					Usage:  "private key file for signning jwt",
					EnvVar: "JWT_PRIVATE_KEY",
				},
				cli.IntFlag{
					Name:  "port, p",
					Usage: "server port",
					Value: 9999,
				},
				cli.StringFlag{
					Name:   "messaging-host",
					EnvVar: "NATS_SERVICE_HOST",
					Usage:  "host address for messaging server",
				},
				cli.StringFlag{
					Name:   "messaging-port",
					EnvVar: "NATS_SERVICE_PORT",
					Usage:  "port for messaging server",
				},
			},
		},
		{
			Name:   "generate-keys",
			Usage:  "generate rsa key pairs(public and private keys) in pem format",
			Action: commands.GenerateKeys,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "private, pr",
					Usage: "output file name for private key",
				},
				cli.StringFlag{
					Name:  "public, pub",
					Usage: "output file name for public key",
				},
			},
		},
	}
	app.Run(os.Args)
}
