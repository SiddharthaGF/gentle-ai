package dsh

import (
	"fmt"
	"os"
	"path/filepath"
)

// EnsureGentleProfile creates a "gentle" DSH profile under home if it does
// not already exist, by copying the shipped "web" template profile structure.
// It returns the profile directory path. The profile is selectable via
// `dsh --profile gentle`.
func EnsureGentleProfile(homeDir string) (string, error) {
	profilesDir := ProfilesDir(homeDir)
	gentleDir := filepath.Join(profilesDir, "gentle")
	if _, err := os.Stat(gentleDir); err == nil {
		return gentleDir, nil
	}
	templateDir := filepath.Join(profilesDir, "web")
	if _, err := os.Stat(templateDir); err != nil {
		templateDir = filepath.Join(profilesDir, "headless")
		if _, err := os.Stat(templateDir); err != nil {
			return "", fmt.Errorf("no template profile (web/headless) found under %s", profilesDir)
		}
	}
	if err := os.MkdirAll(profilesDir, 0o755); err != nil {
		return "", err
	}
	if err := copyDir(templateDir, gentleDir); err != nil {
		return "", fmt.Errorf("copy template to gentle profile: %w", err)
	}
	pkgPath := filepath.Join(gentleDir, "package.json")
	data, err := os.ReadFile(pkgPath)
	if err == nil {
		s := string(data)
		for _, old := range []string{"dsh-profile-web", "dsh-profile-headless"} {
			s = replaceAll(s, old, "dsh-profile-gentle")
		}
		_ = os.WriteFile(pkgPath, []byte(s), 0o644)
	}
	_ = os.WriteFile(filepath.Join(gentleDir, "cordis.patch.yml"), []byte("[]\n"), 0o644)
	return gentleDir, nil
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if info.IsDir() {
			return os.MkdirAll(target, info.Mode())
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, data, info.Mode())
	})
}

func replaceAll(s, old, new string) string {
	result := ""
	for {
		idx := -1
		for i := 0; i+len(old) <= len(s); i++ {
			if s[i:i+len(old)] == old {
				idx = i
				break
			}
		}
		if idx < 0 {
			result += s
			break
		}
		result += s[:idx] + new
		s = s[idx+len(old):]
	}
	return result
}