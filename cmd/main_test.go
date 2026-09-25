package main

import (
	"os"
	"testing"
)

// TestMainCanCall drives main through --version, which parses and prints without serving.
func TestMainCanCall(t *testing.T) {
	t.Parallel()

	originalArgs := os.Args

	defer func() {
		os.Args = originalArgs

		if r := recover(); r != nil {
			t.Fatalf("main() panic: %v", r)
		}
	}()

	os.Args = []string{"twelvedata-exporter", "--version"}

	main()
}
