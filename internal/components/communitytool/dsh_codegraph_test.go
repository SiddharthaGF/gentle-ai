package communitytool

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestProvisionDSHCodeGraph(t *testing.T) {
	dir := t.TempDir()
	if err := ProvisionDSHCodeGraph(dir); err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(filepath.Join(dir, "cordis.patch.yml"))
	if !strings.Contains(string(data), "codegraph") {
		t.Fatalf("entry not written: %s", data)
	}
	// idempotent
	if err := ProvisionDSHCodeGraph(dir); err != nil {
		t.Fatal(err)
	}
}