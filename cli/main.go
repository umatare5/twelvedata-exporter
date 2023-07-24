// Package cli is responsible for the execution of the CLI.
package cli

import (
	"fmt"
	"os"

	server "github.com/umatare5/twelvedata-exporter/internal"
	"github.com/urfave/cli/v2"
)

const (
	webListenAddressFlagName    = "web.listen-address"
	webListenPortFlagName       = "web.listen-port"
	twelvedataAPIKeyFlagName    = "twelvedata.api-key"
	twelvedataRateLimitFlagName = "twelvedata.rate-limit"
)

// Start is a entrypoint of this command
func Start() {
	cmd := &cli.App{
		Name:      "twelvedata-exporter",
		HelpName:  "Fetch metrics from Twelvedata API",
		Usage:     "twelvedata-exporter",
		UsageText: "twelvedata-exporter COMMAND [options...]",
		Version:   "0.1.0",
		Flags:     registerFlags(),
		Action: func(ctx *cli.Context) error {
			config := newConfig(ctx)
			server := server.New(
				config.WebListenAddress,
				config.WebListenPort,
				config.TwelvedataAPIKey,
				config.TwelvedataRateLimit,
			)
			server.Boot()

			return nil
		},
	}

	err := cmd.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// registerFlags returns global flags
func registerFlags() []cli.Flag {
	flags := []cli.Flag{}
	flags = append(flags, registerWebListenAddressFlag()...)
	flags = append(flags, registerWebListenPortFlag()...)
	flags = append(flags, registerAPIKeyFlag()...)
	flags = append(flags, registerRateLimitFlag()...)
	return flags
}

// registerWebListenAddressFlag
func registerWebListenAddressFlag() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    webListenAddressFlagName,
			Usage:   "Set IP address",
			Aliases: []string{"I"},
			Value:   "0.0.0.0",
		},
	}
}

// registerWebListenPortFlag
func registerWebListenPortFlag() []cli.Flag {
	return []cli.Flag{
		&cli.IntFlag{
			Name:    webListenPortFlagName,
			Usage:   "Set port number",
			Aliases: []string{"P"},
			Value:   9341,
		},
	}
}

// registerAPIKeyFlag
func registerAPIKeyFlag() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     twelvedataAPIKeyFlagName,
			Usage:    "Set key to use twelvedata API",
			Aliases:  []string{"a"},
			EnvVars:  []string{"TWELVEDATA_API_KEY"},
			Required: true,
		},
	}
}

// registerRateLimitFlag
func registerRateLimitFlag() []cli.Flag {
	return []cli.Flag{
		&cli.IntFlag{
			Name:    twelvedataRateLimitFlagName,
			Usage:   "Set rate limit per minute to use twelvedata API",
			Aliases: []string{"l"},
			Value:   0,
		},
	}
}
