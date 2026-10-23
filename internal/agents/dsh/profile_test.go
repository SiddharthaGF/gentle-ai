package dsh

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEnsureGentleProfile(t *testing.T) {
	t.Setenv("DSH_HOME", "")
	home := t.TempDir()
	// Create a web template first.
	profilesDir := filepath.Join(home, ".dsh", "profiles")
	if err := os.MkdirAll(filepath.Join(profilesDir, "web"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(profilesDir, "web", "package.json"), []byte(`{"name":"dsh-profile-web","dsh":{"profile":{"bundles":[]}}}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(profilesDir, "web", "cordis.patch.yml"), []byte("[]\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	dir, err := EnsureGentleProfile(home)
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Base(dir) != "gentle" {
		t.Fatalf("expected gentle dir, got %s", dir)
	}
	if _, err := os.Stat(filepath.Join(dir, "package.json")); err != nil {
		t.Fatal("package.json missing")
	}
	// idempotent
	dir2, err := EnsureGentleProfile(home)
	if err != nil {
		t.Fatal(err)
	}
	if dir != dir2 {
		t.Fatalf("not idempotent: %s != %s", dir, dir2)
	}
}