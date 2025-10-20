package commands

import (
	"bytes"
	"strings"
	"testing"
)

func TestRootHelp_NoCompletionAndBasicSections(t *testing.T) {
	cmd := NewRootCmd()

	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{"--help"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute --help: %v", err)
	}

	s := out.String()
	if !strings.Contains(s, "Usage:") {
		t.Fatalf("expected help to contain 'Usage:', got:\n%s", s)
	}
	if !strings.Contains(s, "Flags:") {
		t.Fatalf("expected help to contain 'Flags:', got:\n%s", s)
	}
	if strings.Contains(s, "completion") {
		t.Fatalf("help should not include 'completion' command, got:\n%s", s)
	}
}
