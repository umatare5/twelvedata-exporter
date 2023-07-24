// Package cli is responsible for the execution of the CLI.
package cli

import (
	"errors"
	"log"

	"github.com/jinzhu/configor"
	"github.com/urfave/cli/v2"
)

// Config struct
type Config struct {
	WebListenAddress    string
	WebListenPort       int
	TwelvedataAPIKey    string
	TwelvedataRateLimit int
}

// New returns Config struct
func newConfig(ctx *cli.Context) Config {
	config := Config{
		WebListenAddress:    ctx.String(webListenAddressFlagName),
		WebListenPort:       ctx.Int(webListenPortFlagName),
		TwelvedataAPIKey:    ctx.String(twelvedataAPIKeyFlagName),
		TwelvedataRateLimit: ctx.Int(twelvedataRateLimitFlagName),
	}

	err := configor.New(&configor.Config{}).Load(&config)
	if err != nil {
		log.Fatal(err)
	}

	if err := isValidWebListenAddressFlag(config.WebListenAddress); err != nil {
		log.Fatal(err)
	}

	if err := isValidWebListenPortFlag(config.WebListenPort); err != nil {
		log.Fatal(err)
	}

	if err := isValidTwelvedataAPIKeyFlag(config.TwelvedataAPIKey); err != nil {
		log.Fatal(err)
	}

	if err := isValidTwelvedataRateLimitFlag(config.TwelvedataRateLimit); err != nil {
		log.Fatal(err)
	}

	return config
}

func isValidWebListenAddressFlag(_ string) error {
	return nil
}

func isValidWebListenPortFlag(_ int) error {
	return nil
}

func isValidTwelvedataAPIKeyFlag(apikey string) error {
	if apikey == "" {
		return errors.New("Environment variable 'TWELVEDATA_API_KEY' is not set")
	}

	return nil
}

func isValidTwelvedataRateLimitFlag(_ int) error {
	return nil
}
