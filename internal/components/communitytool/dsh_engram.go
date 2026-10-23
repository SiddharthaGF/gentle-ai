package communitytool

import (
	"fmt"

	dshagent "github.com/gentleman-programming/gentle-ai/v4/internal/agents/dsh"
)

// ProvisionDSHEngram is a no-op marker for DSH: Engram on DSH uses native
// tools registered by gentle-engram, not MCP. It exists so the install flow
// can treat DSH like other agents without special-casing.
func ProvisionDSHEngram(profileDir string) error {
	if profileDir == "" {
		return fmt.Errorf("profileDir is required")
	}
	// Engram is native on DSH; nothing to patch. Keep for API parity.
	return nil
}

// RegisterDSHMCP is a convenience wrapper to register an arbitrary MCP server
// (e.g. a community tool) in a DSH profile.
func RegisterDSHMCP(profileDir, serverName, url string) error {
	return dshagent.EnsureMCPClient(profileDir, serverName, url)
}