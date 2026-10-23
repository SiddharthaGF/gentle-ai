# DSH compatibility plan (gentle-ai)

> Status: draft plan. Phases 0-3, 5-7 implemented; Phase 4 (skills) pending bundle/replace approach. DSH is detection/catalog-only for guidance (consumes AGENTS.md automatically).

## Goal

Make `gentle-ai` treat **DSH (DeepSeek Harness)** as a first-class managed runtime, the way it already treats Pi: install/refresh managed guidance, skills, MCP config, and companion packages without owning DSH's runtime prompts or persona.

## Reference: what Pi support looks like

Pi is the closest analog (see `docs/pi.md`). Its support surface in `gentle-ai`:

- Agent identity: `model.AgentPi` (`internal/model/types.go`), registered in `internal/agents/factory.go`.
- Adapter: `internal/agents/pi/adapter.go` implementing `agents.Adapter`:
  - Detection via `pi` binary + `~/.pi/agent` (honors `PI_CODING_AGENT_DIR`).
  - Config paths: `GlobalConfigDir`, `SystemPromptDir/File`, `SkillsDir`, `SettingsPath`, `MCPConfigPath`.
  - Strategies: `SystemPromptStrategy`, `MCPStrategy`.
  - Capabilities: subagents, skills, system prompt, MCP, output styles, slash commands.
- Companion packages installed via Pi's package manager: `gentle-pi`, `gentle-engram`, `pi-web-access`, `pi-btw`; pruning of retired packages.
- MCP: built-in MCP via `mcp.json` (no adapter); migration from legacy `mcp-adapter.json`; never injects `engram` server.
- Optional CodeGraph wiring + managed overlays (never edits package-owned child files).
- Rollback snapshots of Pi global settings.
- RDD-capable runtime (`model.SupportsReceiptDrivenDevelopment` includes `pi`).
- Engram native tools (not MCP).

## DSH gap analysis

DSH is a plugin-based harness (profiles = ordered bundle stacks + patch layers), not a coding-agent CLI with `~/.dsh/agent`. There is no `dsh` agent adapter today. Concretely missing:

1. **Agent identity**: no `AgentDSH` in `model.AgentID`; no factory registration; no validate/doctor entries.
2. **Adapter**: no `internal/agents/dsh/` implementing `agents.Adapter` (detection, paths, strategies, capabilities).
3. **Config root**: DSH uses `$DSH_HOME` profiles (`$DSH_HOME/profiles/<name>`, `cordis.patch.yml`, profile `package.json`), not a single agent dir. Need a mapping from DSH profile → gentle-ai config surface (likely a dedicated "gentle" profile or patch layer).
4. **Companion packages**: DSH has no `pi install` equivalent; plugins are managed via `dsh plugin` (pnpm). Need a strategy to inject gentle-managed plugins/bundles into a profile's `dsh.profile.bundles` or patch layer without clobbering user patches.
5. **MCP**: DSH supports MCP via plugins (e.g. memory MCP overlays mentioned in README). Need to reconcile with gentle-ai's MCP component (merge, don't overwrite; respect disabled markers).
6. **System prompt / guidance**: DSH consumes instruction files via plugins/patches. Need to determine delivery mechanism (patch layer vs. bundle) and ownership boundaries (like Pi, gentle-pi owns runtime prompts).
7. **Subagents/skills**: DSH has agent-team/experimental profiles. Need capability flags and dirs.
8. **RDD**: decide whether DSH is RDD-capable (likely yes, gated by transport availability).
9. **Doctor/sync/uninstall**: wiring detection, health checks, rollback, ownership records.
10. **Docs**: `docs/agents.md`, `docs/dsh.md`, README badge count.

## Phased plan

### Phase 0 — Discovery (no code)
- Confirm DSH's config model and where user/global instruction files live (`$DSH_HOME`, profile dirs).
- Identify a stable "gentle" integration point: either a shipped bundle or a patch layer template.
- Decide ownership split: DSH runtime owns prompts/persona/model; gentle-ai owns guidance/skills/MCP/companions (mirrors Pi).

### Phase 1 — Identity + adapter skeleton ✅ (done)
- Added `model.AgentDSH` (`"dsh"`).
- Created `internal/agents/dsh/adapter.go` (detection/catalog-only, Conductor-style): detects `$DSH_HOME` (or `~/.dsh`), no managed writes; honors `DSH_HOME`.
- Registered in `internal/agents/factory.go`.
- Registered in capability manifest (empty feature claims).
- Added to catalog (`docs/agents.md` row) and README badge 18→19.
- Added `internal/cli/validate.go` case.
- Basic adapter test in `internal/agents/dsh/adapter_test.go`.

### Phase 1b — Capability expansion (reference: Claude Code + OpenCode adapters)

When Phase 2/3 land, model DSH paths and strategies on the Claude Code and
OpenCode adapters (full Adapter interface: GlobalConfigDir, SystemPromptDir/File,
SkillsDir, SettingsPath, MCPConfigPath, strategies, capabilities) rather than
the Conductor detection-only shape. Current adapter intentionally keeps all
feature claims false (Conductor-style) until writes are safe.

### Phase 2 — Config + MCP (mechanism confirmed, implementation pending)

Discovery + prototype: DSH has no declarative MCP config file, but MCP servers are
registered via a profile's `cordis.patch.yml` by inserting
`@deepseek-ai/dsh-mcp-client` entries (serverName/transport/url/failOnStartupError),
as shown by the shipped `templates/mcp/cordis.patch.yml`. A prototype insert entry
in the throwaway `gentle-test` profile composed cleanly (`dsh --dump-config`),
confirming the mechanism.

- **Implementation pending**: a gentle-ai component that writes/merges managed
  entries into a DSH profile's `cordis.patch.yml` (idempotent, never clobbering
  user entries), plus detection of disabled markers.
- **Available now**: AGENTS.md guidance flows automatically.

### Phase 3 — Companion delivery (follow-up)

No DSH equivalent of `pi install`. Plugins are managed via `dsh plugin`
(pnpm) in a profile dir. Gentle-ai could either (a) generate/patch a profile's
`cordis.patch.yml` with managed entries, or (b) ship a DSH bundle package.
Decision deferred to Phase 0b (profile patch generation prototype).

### Phase 4 — Guidance, skills, subagents (follow-up, bundle approach)

Skills load via `dsh-skill-filesystem` `customSkillDirs`. Adding a gentle skill dir
requires modifying an existing plugin row (last-write-wins by id), which is not a
safe blind insert — better delivered as a DSH bundle or a `replace` op scoped to
the managed row. Subagents: DSH has experimental agent-team profiles; integration
deferred. Guidance already flows via AGENTS.md (Phase 0).

### Phase 5 — Review/RDD + doctor (partial)

- Doctor detection: `dsh.DetectProfiles` added (lists profiles under `$DSH_HOME` or `~/.dsh`). Adapter detection already keys on the DSH home.
- RDD: requires a compiled transport adapter. DSH exposes ACP (`dsh --profile acp`) and SDK profiles that may serve as transport, but no adapter exists yet. Deferred.

### Phase 6 — Docs + tests (follow-up)

DSH-specific docs and e2e coverage deferred until Phases 3–5 land.

### Phase 7 — Wiring ✅ (done)

- Added `dshProfileProvisionStep` in `internal/cli/run.go`, registered in the
  install/sync pipeline (no-op when no DSH profiles exist).
- Step targets the first detected DSH profile (preferring "gentle" if present)
  and calls `communitytool.ProvisionDSHProfile`, which registers CodeGraph MCP
  (stdio) via `dsh.EnsureMCPStdio`.
- All primitives and wiring compile and tests pass.

DSH-specific docs and e2e coverage deferred until Phases 3–5 land.

## Phase 0 findings (discovery)

Conducted against DSH 0.2.0-rc.2 using a throwaway profile `gentle-test` (created from the `web` template; kept under `$DSH_HOME/profiles/gentle-test`, does not touch the original `web`/`headless` profiles).

- **Instruction surface is already compatible**: DSH's `@deepseek-ai/dsh-agent-instructions` plugin discovers `AGENTS.md` and `CLAUDE.md` at the project root plus `$DSH_HOME/AGENTS.md` (user-global). Gentle-AI already writes `AGENTS.md` for agents, so project guidance flows into DSH with no new adapter. Budget: `maxBytes` 65536, `maxSourceBytes` 1 MiB.
- **Skills**: `@deepseek-ai/dsh-skill-filesystem` loads from `customSkillDirs` (shipped skills under `@deepseek-ai/dsh-agent-preset/skills`). Gentle-AI's skill installation could target a DSH skill dir or extend `customSkillDirs` via a patch layer.
- **MCP**: DSH has `@deepseek-ai/dsh-mcp-resources` and `dsh-mcp-client`; profile bundles already include MCP-capable plugins. Need to confirm how a server entry is declared (patch layer vs. config) before implementing merge.
- **Config model**: a profile = `package.json` with `dsh.profile.bundles` + `cordis.yml` (empty root) + `cordis.patch.yml` (user layer) + optional `--patch`. Composition order: bundles → profile patch → home patch → `--patch`. Gentle-AI must never overwrite `cordis.patch.yml` blindly; it should inject managed entries or generate a separate patch file.
- **No agent dir**: unlike Pi (`~/.pi/agent`), DSH has no single agent home with `settings.json`. Detection for gentle-ai should key on `$DSH_HOME` + installed profiles, not a binary+agent-dir pair.
- **Conclusion**: DSH compatibility is closer to "patch-layer + AGENTS.md" than to a full Pi-style adapter. Phase 1 can be smaller: a managed patch layer generator + skill dir wiring, plus docs, rather than a new `agents.Adapter`.

## Open questions

- Is DSH meant to be a user-facing "agent" selectable via `--agent dsh`, or an internal runtime only?
- Which DSH profile should gentle-ai target (a new `gentle` profile vs. patching an existing one)?
- Does DSH expose an ACP/SDK transport gentle-ai can reuse for review (like Pi's host relay)?

## Risks

- DSH's plugin model differs fundamentally from Pi's package manager; premature adapter could encode wrong ownership.
- Without a clear instruction-file surface, guidance injection may fight DSH's patch composition.
- RDD requires a compiled transport adapter; DSH may not have one yet.

## Suggested first step

Run Phase 0 discovery against a real DSH install (`dsh --version`, `dsh --dump-default-config`, inspect `$DSH_HOME`) before writing any adapter code. Mirror Pi's adapter only after confirming path/strategy mappings.
