package logger

import "testing"

func TestParseLevel(t *testing.T) {
	cases := []string{"debug", "info", "warn", "warning", "error", "fatal", ""}
	for _, in := range cases {
		_ = parseLevel(in)
	}
}
