package config

import (
	"context"
	"testing"

	cli "github.com/urfave/cli/v3"
)

// runNewConfig builds a Config through a real command, mirroring the flags the cli package registers.
func runNewConfig(t *testing.T, args ...string) Config {
	t.Helper()

	var got Config

	cmd := &cli.Command{
		Name: "twelvedata-exporter",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: WebListenAddressFlagName, Value: "0.0.0.0"},
			&cli.IntFlag{Name: WebListenPortFlagName, Value: 10016},
			&cli.StringFlag{Name: WebScrapePathFlagName, Value: "/price"},
			&cli.StringFlag{Name: TwelvedataAPIKeyFlagName, Required: true},
		},
		Action: func(_ context.Context, cmd *cli.Command) error {
			got = NewConfig(cmd)

			return nil
		},
	}

	if err := cmd.Run(context.Background(), append([]string{"twelvedata-exporter"}, args...)); err != nil {
		t.Fatalf("cli.Command.Run() error = %v, want nil", err)
	}

	return got
}

// TestNewConfigReadsEveryFlag checks each flag lands in the field the server reads.
func TestNewConfigReadsEveryFlag(t *testing.T) {
	t.Parallel()

	got := runNewConfig(t,
		"--"+WebListenAddressFlagName, "127.0.0.1",
		"--"+WebListenPortFlagName, "9000",
		"--"+WebScrapePathFlagName, "/metrics",
		"--"+TwelvedataAPIKeyFlagName, "secret",
	)

	want := Config{
		WebListenAddress: "127.0.0.1",
		WebListenPort:    9000,
		WebScrapePath:    "/metrics",
		TwelvedataAPIKey: "secret",
	}

	if got != want {
		t.Errorf("NewConfig() = %+v, want %+v", got, want)
	}
}

// TestNewConfigAppliesDefaults checks a run naming only the key keeps the published defaults.
func TestNewConfigAppliesDefaults(t *testing.T) {
	t.Parallel()

	got := runNewConfig(t, "--"+TwelvedataAPIKeyFlagName, "secret")

	want := Config{
		WebListenAddress: "0.0.0.0",
		WebListenPort:    10016,
		WebScrapePath:    "/price",
		TwelvedataAPIKey: "secret",
	}

	if got != want {
		t.Errorf("NewConfig() = %+v, want %+v", got, want)
	}
}

// TestIsValidTwelvedataAPIKeyFlag checks an absent key is the one rejected value.
func TestIsValidTwelvedataAPIKeyFlag(t *testing.T) {
	t.Parallel()

	if err := isValidTwelvedataAPIKeyFlag("secret"); err != nil {
		t.Errorf("isValidTwelvedataAPIKeyFlag(%q) error = %v, want nil", "secret", err)
	}

	err := isValidTwelvedataAPIKeyFlag("")
	if err == nil {
		t.Fatal("isValidTwelvedataAPIKeyFlag(\"\") error = nil, want an error")
	}

	const want = "environment variable 'TWELVEDATA_API_KEY' is not set"
	if err.Error() != want {
		t.Errorf("error = %q, want %q", err.Error(), want)
	}
}

// TestWebFlagsAcceptAnyValue checks the web flags are deliberately unvalidated.
// The listener rejects what it cannot bind, so a second check here would only duplicate it.
func TestWebFlagsAcceptAnyValue(t *testing.T) {
	t.Parallel()

	if err := isValidWebListenAddressFlag("not-an-address"); err != nil {
		t.Errorf("isValidWebListenAddressFlag() error = %v, want nil", err)
	}

	if err := isValidWebListenPortFlag(-1); err != nil {
		t.Errorf("isValidWebListenPortFlag() error = %v, want nil", err)
	}

	if err := isValidWebScrapePathFlag("no-leading-slash"); err != nil {
		t.Errorf("isValidWebScrapePathFlag() error = %v, want nil", err)
	}
}
