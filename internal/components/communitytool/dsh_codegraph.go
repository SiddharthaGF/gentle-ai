package communitytool

import (
	"fmt"

	dshagent "github.com/gentleman-programming/gentle-ai/v4/internal/agents/dsh"
)

// ProvisionDSHCodeGraph registers the CodeGraph MCP server in a DSH profile's
// patch layer via dsh.EnsureMCPStdio (stdio transport: command "codegraph",
// args ["mcp"]). It is the DSH analog of the Pi CodeGraph wiring, reduced to
// the MCP-registration primitive.
func ProvisionDSHCodeGraph(profileDir string) error {
	if profileDir == "" {
		return fmt.Errorf("profileDir is required")
	}
	return dshagent.EnsureMCPStdio(profileDir, "codegraph", "codegraph", []string{"mcp"})
}