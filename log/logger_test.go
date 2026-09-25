package log

import (
	"bytes"
	"strings"
	"testing"
)

// TestLoggerWritesEveryLevel checks each wrapper reaches the shared logger at its level.
// One test owns the logger, because the package exposes a single instance every wrapper writes to.
func TestLoggerWritesEveryLevel(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	logger.SetOutput(&buf)
	logger.ExitFunc = func(int) {}

	tests := []struct {
		level string
		write func()
		want  string
	}{
		{"info", func() { Info("started") }, "started"},
		{"infof", func() { Infof("port %d", 10016) }, "port 10016"},
		{"errorf", func() { Errorf("failed: %s", "timeout") }, "failed: timeout"},
		{"fatal", func() { Fatal("unrecoverable") }, "unrecoverable"},
	}

	for _, tt := range tests {
		buf.Reset()
		tt.write()

		got := buf.String()
		if !strings.Contains(got, tt.want) {
			t.Errorf("%s wrote %q, want it to contain %q", tt.level, got, tt.want)
		}

		if !strings.Contains(got, "level="+strings.TrimSuffix(tt.level, "f")) {
			t.Errorf("%s wrote %q, want the %s level", tt.level, got, tt.level)
		}
	}
}
