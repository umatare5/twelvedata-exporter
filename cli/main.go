// Package cli is responsible for the execution of the CLI.
package cli

import (
	"fmt"
	"os"

	"github.com/umatare5/twelvedata-exporter/config"
	"github.com/umatare5/twelvedata-exporter/internal"
	"github.com/urfave/cli/v2"
)

// Start is the entrypoint of this CLI
func Start() {
	cmd := &cli.App{
		Name:      "twelvedata-exporter",
		HelpName:  "Fetch quotes from Twelvedata API",
		Usage:     "twelvedata-exporter",
		UsageText: "twelvedata-exporter COMMAND [options...]",
		Version:   "1.0.0",
		Flags:     registerFlags(),
		Action: func(ctx *cli.Context) error {
			config := config.NewConfig(ctx)
			server, _ := internal.NewServer(&config)

			server.Start()

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
	flags = append(flags, registerWebScrapePathFlag()...)
	flags = append(flags, registerAPIKeyFlag()...)
	return flags
}

// registerWebListenAddressFlag
func registerWebListenAddressFlag() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    config.WebListenAddressFlagName,
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
			Name:    config.WebListenPortFlagName,
			Usage:   "Set port number",
			Aliases: []string{"P"},
			Value:   10016,
		},
	}
}

// registerWebScrapePathFlag
func registerWebScrapePathFlag() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    config.WebScrapePathFlagName,
			Usage:   "Set the path to expose metrics",
			Aliases: []string{"p"},
			Value:   "/price",
		},
	}
}

// registerAPIKeyFlag
func registerAPIKeyFlag() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     config.TwelvedataAPIKeyFlagName,
			Usage:    "Set key to use twelvedata API",
			Aliases:  []string{"a"},
			EnvVars:  []string{"TWELVEDATA_API_KEY"},
			Required: true,
		},
	}
}
