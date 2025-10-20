package commands

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestInit_WritesConfigWithServer(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("perm checks differ on windows; covered in config tests")
	}
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", "")

	cmd := NewRootCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{"init", "--server", "http://localhost:8080"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("init: %v; out=%s", err, out.String())
	}

	want := filepath.Join(home, ".config", "poppo-press", "config.yaml")
	if _, err := os.Stat(want); err != nil {
		t.Fatalf("config not written at %s: %v", want, err)
	}
}

func TestConfigTZ_GetSet(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", "")

	// init first
	root := NewRootCmd()
	root.SetArgs([]string{"init", "--server", "http://localhost:8080"})
	if err := root.Execute(); err != nil {
		t.Fatalf("init: %v", err)
	}

	// set tz
	set := NewRootCmd()
	set.SetArgs([]string{"config", "tz", "set", "Asia/Tokyo"})
	if err := set.Execute(); err != nil {
		t.Fatalf("tz set: %v", err)
	}

	// get tz
	var out bytes.Buffer
	get := NewRootCmd()
	get.SetOut(&out)
	get.SetArgs([]string{"config", "tz"})
	if err := get.Execute(); err != nil {
		t.Fatalf("tz get: %v", err)
	}
	if got := out.String(); got != "Asia/Tokyo\n" {
		t.Fatalf("got %q, want %q", got, "Asia/Tokyo\n")
	}
}
