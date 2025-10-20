package config

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestPath_UnixDefault(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("unix-only path test")
	}
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", "")
	p, err := Path()
	if err != nil {
		t.Fatalf("Path(): %v", err)
	}
	want := filepath.Join(home, ".config", "poppo-press", "config.yaml")
	if p != want {
		t.Fatalf("got %q, want %q", p, want)
	}
}

func TestSaveLoad_RoundtripAndPerms_Unix(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("perm bits vary on windows")
	}
	cfgHome := filepath.Join(t.TempDir(), ".config")
	t.Setenv("XDG_CONFIG_HOME", cfgHome)
	t.Setenv("HOME", t.TempDir())

	cfg := &Config{Server: "http://localhost:8080", Token: "", Timezone: "", Output: Output{Pager: "auto"}}
	if err := Save(cfg); err != nil {
		t.Fatalf("Save: %v", err)
	}

	p, err := Path()
	if err != nil {
		t.Fatalf("Path: %v", err)
	}
	st, err := os.Stat(p)
	if err != nil {
		t.Fatalf("stat file: %v", err)
	}
	if got, want := st.Mode().Perm(), os.FileMode(0o600); got != want {
		t.Fatalf("file perms = %v, want %v", got, want)
	}
	dst := filepath.Dir(p)
	dstSt, err := os.Stat(dst)
	if err != nil {
		t.Fatalf("stat dir: %v", err)
	}
	if got, want := dstSt.Mode().Perm(), os.FileMode(0o700); got != want {
		t.Fatalf("dir perms = %v, want %v", got, want)
	}

	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.Server != cfg.Server || loaded.Output.Pager != cfg.Output.Pager {
		t.Fatalf("loaded mismatch: %#v", loaded)
	}
}
