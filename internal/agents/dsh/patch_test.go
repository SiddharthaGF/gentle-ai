package dsh

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestEnsureMCPClientIdempotent(t *testing.T) {
	dir := t.TempDir()
	if err := EnsureMCPClient(dir, "engram", "http://127.0.0.1:9999/mcp"); err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(ProfilePatchPath(dir))
	if !strings.Contains(string(data), "gentle-engram") {
		t.Fatalf("entry not written: %s", data)
	}
	// second call must not duplicate
	if err := EnsureMCPClient(dir, "engram", "http://127.0.0.1:9999/mcp"); err != nil {
		t.Fatal(err)
	}
	data2, _ := os.ReadFile(ProfilePatchPath(dir))
	var ops []map[string]any
	if err := yaml.Unmarshal(data2, &ops); err != nil {
		t.Fatal(err)
	}
	count := 0
	for _, op := range ops {
		for _, e := range op["insert"].([]any) {
			m := e.(map[string]any)
			if m["id"] == "gentle-engram" {
				count++
			}
		}
	}
	if count != 1 {
		t.Fatalf("expected 1 entry, got %d", count)
	}
	// preserve user entries
	path := filepath.Join(dir, "cordis.patch.yml")
	_ = path
}
func TestEnsurePluginIdempotent(t *testing.T) {
	dir := t.TempDir()
	if err := EnsurePlugin(dir, "gentle-skill", "@deepseek-ai/dsh-skill-filesystem", map[string]any{"customSkillDirs": []string{"/tmp/skills"}}); err != nil {
		t.Fatal(err)
	}
	if err := EnsurePlugin(dir, "gentle-skill", "@deepseek-ai/dsh-skill-filesystem", nil); err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(ProfilePatchPath(dir))
	if strings.Count(string(data), "gentle-skill") != 1 {
		t.Fatalf("duplicate or missing entry: %s", data)
	}
}
