package communitytool

import (
	"fmt"
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
	if err := ProvisionDSHCodeGraph(opts.ProfileDir); err != nil {
		return fmt.Errorf("codegraph MCP: %w", err)
	}
	return nil
}