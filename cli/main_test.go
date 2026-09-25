package cli

import (
	"os"
	"slices"
	"testing"

	cli "github.com/urfave/cli/v3"

	"github.com/umatare5/twelvedata-exporter/config"
)

// flagsByName indexes the registered flags so a case can look one up by its long name.
func flagsByName(t *testing.T) map[string]cli.Flag {
	t.Helper()

	flags := registerFlags()
	byName := make(map[string]cli.Flag, len(flags))

	for _, flag := range flags {
		byName[flag.Names()[0]] = flag
	}

	if len(byName) != len(flags) {
		t.Fatalf("registerFlags() declares %d flags under %d names", len(flags), len(byName))
	}

	return byName
}

// TestRegisterFlagsDeclaresEveryFlag checks each flag carries the short form the README documents.
func TestRegisterFlagsDeclaresEveryFlag(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		alias    string
		required bool
	}{
		{config.WebListenAddressFlagName, "I", false},
		{config.WebListenPortFlagName, "P", false},
		{config.WebScrapePathFlagName, "p", false},
		{config.TwelvedataAPIKeyFlagName, "a", true},
	}

	byName := flagsByName(t)

	if len(byName) != len(tests) {
		t.Fatalf("registerFlags() = %d flags, want %d", len(byName), len(tests))
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			flag, ok := byName[tt.name]
			if !ok {
				t.Fatalf("registerFlags() omits %s", tt.name)
			}

			if !slices.Contains(flag.Names(), tt.alias) {
				t.Errorf("names = %v, want the %q alias", flag.Names(), tt.alias)
			}

			required, ok := flag.(cli.RequiredFlag)
			if !ok {
				t.Fatalf("%s does not report whether it is required", tt.name)
			}

			if required.IsRequired() != tt.required {
				t.Errorf("IsRequired() = %v, want %v", required.IsRequired(), tt.required)
			}
		})
	}
}

// TestRegisterFlagsDefaults checks an omitted flag falls back to the value the README documents.
func TestRegisterFlagsDefaults(t *testing.T) {
	t.Parallel()

	byName := flagsByName(t)

	strings := map[string]string{
		config.WebListenAddressFlagName: "0.0.0.0",
		config.WebScrapePathFlagName:    "/price",
	}

	for name, want := range strings {
		flag, ok := byName[name].(*cli.StringFlag)
		if !ok {
			t.Fatalf("%s is not a string flag", name)
		}

		if flag.Value != want {
			t.Errorf("%s default = %q, want %q", name, flag.Value, want)
		}
	}

	port, ok := byName[config.WebListenPortFlagName].(*cli.IntFlag)
	if !ok {
		t.Fatalf("%s is not an int flag", config.WebListenPortFlagName)
	}

	if port.Value != 10016 {
		t.Errorf("%s default = %d, want 10016", config.WebListenPortFlagName, port.Value)
	}
}

// TestRegisterAPIKeyFlagReadsEnvironment checks the key can arrive without appearing in the process list.
func TestRegisterAPIKeyFlagReadsEnvironment(t *testing.T) {
	t.Parallel()

	flag, ok := flagsByName(t)[config.TwelvedataAPIKeyFlagName].(*cli.StringFlag)
	if !ok {
		t.Fatalf("%s is not a string flag", config.TwelvedataAPIKeyFlagName)
	}

	const want = `environment variable "TWELVEDATA_API_KEY"`
	if got := flag.Sources.String(); got != want {
		t.Errorf("Sources = %q, want %q", got, want)
	}
}

// TestGetVersion checks an unstamped build reports the placeholder the Makefile replaces.
func TestGetVersion(t *testing.T) {
	t.Parallel()

	if got := getVersion(); got != "dev" {
		t.Errorf("getVersion() = %q, want %q", got, "dev")
	}
}

// TestStart drives the command from the arguments, covering the run that prints and the run that reaches the server.
// Each case names the key, because the flag reads an environment variable a developer shell may already export.
func TestStart(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{"version", []string{"--version"}, false},
		{"unbindable port", []string{"--twelvedata.api-key", "secret", "--web.listen-port", "-1"}, true},
	}

	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	for _, tt := range tests {
		os.Args = append([]string{"twelvedata-exporter"}, tt.args...)

		if err := Start(); (err != nil) != tt.wantErr {
			t.Errorf("%s: Start() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}
