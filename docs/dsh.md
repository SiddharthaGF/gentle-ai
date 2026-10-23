# DSH (DeepSeek Harness) integration

> [!NOTE]
> These docs track `main`, which may include unreleased changes.

← [Back to README](../README.md)

Gentle AI detects DSH and registers it as a supported agent (`dsh`). DSH is a
plugin-based harness: a profile is an ordered stack of bundle layers plus a user
patch layer (`cordis.patch.yml`), not a single agent home with settings.

## Current status

- **Detection + catalog**: Gentle AI detects `$DSH_HOME` (or `~/.dsh`) and lists
  profiles via `dsh.DetectProfiles`.
- **Guidance**: DSH already consumes `AGENTS.md` / `CLAUDE.md` at the project
  root and `$DSH_HOME/AGENTS.md`, so project guidance flows in automatically.
  No managed writes are performed.
- **MCP**: Gentle AI generates patch-layer entries for MCP servers via
  `dsh.EnsureMCPClient` and provisions CodeGraph MCP (stdio) into the `gentle`
  profile during install/sync.
- **Skills**: DSH loads skills from `customSkillDirs` in `dsh-skill-filesystem`.
  Adding a gentle skill dir needs a bundle or scoped replace (not a blind
  insert); pending.
- **Companion packages**: DSH has no `pi install` equivalent; plugins are
  managed via `dsh plugin` (pnpm). `dsh.EnsurePlugin` is the primitive.
- **RDD**: pending a compiled transport adapter (ACP/SDK profiles).

## Selectable profile

During install/sync, Gentle AI creates a `gentle` profile (`$DSH_HOME/profiles/gentle`)
if it does not exist, by copying the shipped `web` template. It is selectable via:

```sh
dsh --profile gentle
```

The profile's patch layer is provisioned with managed entries (CodeGraph MCP, etc.).

## Paths

- Home: `$DSH_HOME` or `~/.dsh`
- Profiles: `$DSH_HOME/profiles/<name>`
- User patch: `$DSH_HOME/profiles/<name>/cordis.patch.yml`
- User-global instructions: `$DSH_HOME/AGENTS.md`

## Next steps

See [docs/dsh-compatibility-plan.md](dsh-compatibility-plan.md) for the phased
plan and open questions.

← [Back to README](../README.md)