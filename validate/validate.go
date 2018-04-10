package validate

import (
	"fmt"

	"gopkg.in/urfave/cli.v1"
)

func ValidateRunArgs(c *cli.Context) error {
	for _, p := range []string{
		"config",
		"public-key",
		"private-key",
		"messaging-host",
		"messaging-port",
	} {
		if len(c.String(p)) == 0 {
			return cli.NewExitError(
				fmt.Sprintf("argument %s is missing", p),
				2,
			)
		}
	}
}
