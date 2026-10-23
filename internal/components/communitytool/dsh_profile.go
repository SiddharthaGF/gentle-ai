package communitytool

import (
	"fmt"

	dshagent "github.com/gentleman-programming/gentle-ai/v4/internal/agents/dsh"
)

// DSHProfileProvisionOptions holds the parameters for provisioning a DSH profile.
type DSHProfileProvisionOptions struct {
	ProfileDir    string
	CodeGraphURL  string
	CodeGraphName string
}

// ProvisionDSHProfile runs all managed DSH provisions for a profile:
// CodeGraph MCP (stdio) and any future managed plugins.
func ProvisionDSHProfile(opts DSHProfileProvisionOptions) error {
	if opts.ProfileDir == "" {
		return fmt.Errorf("profileDir is required")
	}
	// Best-effort: make the gentleman mode default so ODD guidance loads.
	_ = dshagent.EnsureGentleDefaultMode(opts.ProfileDir)
	if err := ProvisionDSHCodeGraph(opts.ProfileDir); err != nil {
		return fmt.Errorf("codegraph MCP: %w", err)
	}
	return nil
}