package dsh

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectProfiles(t *testing.T) {
	t.Setenv("DSH_HOME", "")
	home := t.TempDir()
	profiles := filepath.Join(home, ".dsh", "profiles")
	if err := os.MkdirAll(filepath.Join(profiles, "web"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(profiles, "headless"), 0o755); err != nil {
		t.Fatal(err)
	}
	got, err := DetectProfiles(home)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 profiles, got %v", got)
	}
	// nonexistent home returns nil
	empty, err := DetectProfiles(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if len(empty) != 0 {
		t.Fatalf("expected empty, got %v", empty)
	}
}