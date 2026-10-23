package dsh

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// ProfilePatchPath returns the user patch layer path for a DSH profile dir.
func ProfilePatchPath(profileDir string) string {
	return filepath.Join(profileDir, "cordis.patch.yml")
}

// EnsureMCPClient idempotently inserts a @deepseek-ai/dsh-mcp-client entry
// for serverName into the profile's cordis.patch.yml. It never overwrites an
// existing entry for the same serverName and preserves unrelated user entries.
func EnsureMCPClient(profileDir, serverName, url string) error {
	path := ProfilePatchPath(profileDir)
	data, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	var ops []map[string]any
	if len(data) > 0 {
		if err := yaml.Unmarshal(data, &ops); err != nil {
			return fmt.Errorf("parse %s: %w", path, err)
		}
	}

	for _, op := range ops {
		entries, _ := op["insert"].([]any)
		for _, e := range entries {
			m, _ := e.(map[string]any)
			if m == nil {
				continue
			}
			cfg, _ := m["config"].(map[string]any)
			if cfg != nil && cfg["serverName"] == serverName {
				return nil
			}
		}
	}

	entry := map[string]any{
		"id":   "gentle-" + serverName,
		"name": "@deepseek-ai/dsh-mcp-client",
		"config": map[string]any{
			"serverName":         serverName,
			"transport":          "streamable-http",
			"url":                url,
			"failOnStartupError": false,
		},
	}

	placed := false
	for _, op := range ops {
		if _, ok := op["insert"]; ok {
			op["insert"] = append(op["insert"].([]any), entry)
			placed = true
			break
		}
	}
	if !placed {
		ops = append(ops, map[string]any{"insert": []any{entry}})
	}

	out, err := yaml.Marshal(ops)
	if err != nil {
		return err
	}
	return os.WriteFile(path, out, 0o644)
}
// EnsurePlugin idempotently inserts a plugin row (id/name/config) into the
// profile's cordis.patch.yml. It never overwrites an existing row with the
// same id and preserves unrelated user entries.
func EnsurePlugin(profileDir, id, name string, config map[string]any) error {
	path := ProfilePatchPath(profileDir)
	data, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	var ops []map[string]any
	if len(data) > 0 {
		if err := yaml.Unmarshal(data, &ops); err != nil {
			return fmt.Errorf("parse %s: %w", path, err)
		}
	}

	for _, op := range ops {
		entries, _ := op["insert"].([]any)
		for _, e := range entries {
			m, _ := e.(map[string]any)
			if m != nil && m["id"] == id {
				return nil
			}
		}
	}

	entry := map[string]any{"id": id, "name": name}
	if config != nil {
		entry["config"] = config
	}

	placed := false
	for _, op := range ops {
		if _, ok := op["insert"]; ok {
			op["insert"] = append(op["insert"].([]any), entry)
			placed = true
			break
		}
	}
	if !placed {
		ops = append(ops, map[string]any{"insert": []any{entry}})
	}

	out, err := yaml.Marshal(ops)
	if err != nil {
		return err
	}
	return os.WriteFile(path, out, 0o644)
}

// EnsureMCPStdio idempotently inserts a @deepseek-ai/dsh-mcp-client entry
// using stdio transport (command + args) for a server into the profile's
// cordis.patch.yml. It never overwrites an existing entry for the same
// serverName and preserves unrelated user entries.
func EnsureMCPStdio(profileDir, serverName, command string, args []string) error {
	path := ProfilePatchPath(profileDir)
	data, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	var ops []map[string]any
	if len(data) > 0 {
		if err := yaml.Unmarshal(data, &ops); err != nil {
			return fmt.Errorf("parse %s: %w", path, err)
		}
	}

	for _, op := range ops {
		entries, _ := op["insert"].([]any)
		for _, e := range entries {
			m, _ := e.(map[string]any)
			if m == nil {
				continue
			}
			cfg, _ := m["config"].(map[string]any)
			if cfg != nil && cfg["serverName"] == serverName {
				return nil
			}
		}
	}

	cfg := map[string]any{
		"serverName":         serverName,
		"transport":          "stdio",
		"command":            command,
		"failOnStartupError": false,
	}
	if len(args) > 0 {
		cfg["args"] = args
	}

	entry := map[string]any{
		"id":     "gentle-" + serverName,
		"name":   "@deepseek-ai/dsh-mcp-client",
		"config": cfg,
	}

	placed := false
	for _, op := range ops {
		if _, ok := op["insert"]; ok {
			op["insert"] = append(op["insert"].([]any), entry)
			placed = true
			break
		}
	}
	if !placed {
		ops = append(ops, map[string]any{"insert": []any{entry}})
	}

	out, err := yaml.Marshal(ops)
	if err != nil {
		return err
	}
	return os.WriteFile(path, out, 0o644)
}
