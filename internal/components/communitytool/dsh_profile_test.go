package communitytool

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestProvisionDSHProfile(t *testing.T) {
	dir := t.TempDir()
	err := ProvisionDSHProfile(DSHProfileProvisionOptions{
		ProfileDir:   dir,
		CodeGraphURL: "http://127.0.0.1:8080/mcp",
	})
	if err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(filepath.Join(dir, "cordis.patch.yml"))
	if !strings.Contains(string(data), "codegraph") {
		t.Fatalf("entry not written: %s", data)
	}
	// idempotent
	if err := ProvisionDSHProfile(DSHProfileProvisionOptions{
		ProfileDir:   dir,
		CodeGraphURL: "http://127.0.0.1:8080/mcp",
	}); err != nil {
		t.Fatal(err)
	}
}