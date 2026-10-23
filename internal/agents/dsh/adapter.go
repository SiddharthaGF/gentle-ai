// Package dsh provides DeepSeek Harness (DSH) integration.
//
// DSH is a plugin-based harness: a profile is an ordered stack of bundle
// layers plus a user patch layer, not a single agent home with settings. DSH
// already consumes AGENTS.md/CLAUDE.md at the project root and
// $DSH_HOME/AGENTS.md, so project guidance flows in automatically. Gentle AI
// therefore starts as a detection/catalog-only integration: it registers DSH
// and detects its home directory, but performs no DSH-specific writes yet.
// Patch-layer generation and skill wiring are follow-up phases.
package dsh

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gentleman-programming/gentle-ai/v4/internal/agents/capabilitymanifest"
	"github.com/gentleman-programming/gentle-ai/v4/internal/model"
	"github.com/gentleman-programming/gentle-ai/v4/internal/system"
)

type statResult struct {
	isDir bool
	err   error
}

type Adapter struct {
	statPath func(string) statResult
}

func NewAdapter() *Adapter {
	return &Adapter{statPath: defaultStat}
}

func (a *Adapter) Agent() model.AgentID { return model.AgentDSH }
func (a *Adapter) Tier() model.SupportTier { return model.TierFull }

func (a *Adapter) Detect(_ context.Context, homeDir string) (bool, string, string, bool, error) {
	configPath := ConfigPath(homeDir)
	stat := a.statPath(configPath)
	if stat.err != nil {
		if os.IsNotExist(stat.err) {
			return false, "", configPath, false, nil
		}
		return false, "", "", false, stat.err
	}
	return stat.isDir, "", configPath, stat.isDir, nil
}

func (a *Adapter) CapabilityManifest() capabilitymanifest.AgentCapabilityManifest {
	return capabilitymanifest.MustForAgent(model.AgentDSH)
}

func (a *Adapter) InstallCommand(_ system.PlatformProfile) ([][]string, error) {
	return nil, AgentNotInstallableError{Agent: a.Agent()}
}

func (a *Adapter) GlobalConfigDir(homeDir string) string { return ConfigPath(homeDir) }
func (a *Adapter) SystemPromptDir(homeDir string) string { return ConfigPath(homeDir) }
func (a *Adapter) SystemPromptFile(homeDir string) string { return filepath.Join(ConfigPath(homeDir), "AGENTS.md") }
func (a *Adapter) SkillsDir(_ string) string { return "" }
func (a *Adapter) SettingsPath(_ string) string { return "" }

func (a *Adapter) SystemPromptStrategy() model.SystemPromptStrategy {
	return model.StrategyMarkdownSections
}
func (a *Adapter) MCPStrategy() model.MCPStrategy {
	return model.StrategySeparateMCPFiles
}
func (a *Adapter) MCPConfigPath(_ string, _ string) string { return "" }

func (a *Adapter) SupportsOutputStyles() bool { return a.CapabilityManifest().Features.OutputStyles }
func (a *Adapter) OutputStyleDir(_ string) string { return "" }
func (a *Adapter) SupportsSlashCommands() bool { return a.CapabilityManifest().Features.SlashCommands }
func (a *Adapter) CommandsDir(_ string) string { return "" }
func (a *Adapter) SupportsSubAgents() bool { return a.CapabilityManifest().Features.FileSubAgents }
func (a *Adapter) SubAgentsDir(_ string) string { return "" }
func (a *Adapter) EmbeddedSubAgentsDir() string { return "" }
func (a *Adapter) SupportsSkills() bool { return a.CapabilityManifest().Features.Skills }
func (a *Adapter) SupportsSystemPrompt() bool { return a.CapabilityManifest().Features.SystemPrompt }
func (a *Adapter) SupportsMCP() bool { return a.CapabilityManifest().Features.MCP }

func defaultStat(path string) statResult {
	info, err := os.Stat(path)
	if err != nil {
		return statResult{err: err}
	}
	return statResult{isDir: info.IsDir()}
}

// ConfigPath returns $DSH_HOME (or ~/.dsh) — DSH's harness home.
func ConfigPath(homeDir string) string {
	if v := os.Getenv("DSH_HOME"); v != "" {
		return v
	}
	return filepath.Join(homeDir, ".dsh")
}

type AgentNotInstallableError struct {
	Agent model.AgentID
}
func (e AgentNotInstallableError) Error() string {
	return fmt.Sprintf("agent %q must be installed manually before Gentle AI can configure it", e.Agent)
}
