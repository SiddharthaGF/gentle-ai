package dsh

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestAdapterIdentityAndDetection(t *testing.T) {
	t.Setenv("DSH_HOME", "")
	home := t.TempDir()
	// DSH home absent -> not detected.
	a := NewAdapter()
	installed, _, configPath, configFound, err := a.Detect(context.Background(), home)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}
	if installed || configFound {
		t.Fatalf("expected not detected, got installed=%v configFound=%v", installed, configFound)
	}
	if configPath != filepath.Join(home, ".dsh") {
		t.Fatalf("configPath = %q", configPath)
	}

	// Create DSH home -> detected.
	dshHome := filepath.Join(home, ".dsh")
	if err := os.MkdirAll(dshHome, 0o755); err != nil {
		t.Fatal(err)
	}
	installed, _, _, configFound, err = a.Detect(context.Background(), home)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}
	if !installed || !configFound {
		t.Fatalf("expected detected, got installed=%v configFound=%v", installed, configFound)
	}
	if a.Agent() == "" {
		t.Fatal("empty agent id")
	}
}