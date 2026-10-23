package dsh

import (
	"os"
	"path/filepath"
)

// ProfilesDir returns the DSH profiles directory under home.
func ProfilesDir(homeDir string) string {
	if v := os.Getenv("DSH_HOME"); v != "" {
		return filepath.Join(v, "profiles")
	}
	return filepath.Join(homeDir, ".dsh", "profiles")
}

// DetectProfiles returns the names of existing DSH profiles under home.
func DetectProfiles(homeDir string) ([]string, error) {
	dir := ProfilesDir(homeDir)
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var out []string
	for _, e := range entries {
		if e.IsDir() {
			out = append(out, e.Name())
		}
	}
	return out, nil
}