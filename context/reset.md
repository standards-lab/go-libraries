# reset · init

- **Status:** closeout
- **Session type:** development
- **Branch:** main

## Disposition

- **Established:** the marathon-managed `go-libraries` repository — a public, multi-module Go monorepo for
  the standards-lab reference libraries (Priority 2). Scaffolded the workspace skeleton (`go.work`,
  `mise.toml`, matrix CI, the per-module release workflow, `.gitignore`, `README.md`) and the Claude
  configuration (`CLAUDE.md`, `.claude/settings.json`, `.claude/marathon.toml`).
- **Settled** (`design/`): `library-topology-and-naming` (a single Go monorepo; the module-path,
  vendor-submodule, and tag conventions), `release-and-ci` (per-module prefix-tag releases, the matrix CI,
  the `go.work`/`mise` workflow), and `conventions` (interface-in-root with vendor-in-submodule, explicit
  registration, co-located black-box tests, doc.go ownership).
- **Candidate** (`concepts/`): the module set to re-derive (`module-set`).

## Next-focus

Build the **`core`** foundation module — the Layered Composition Architecture lifecycle (cold start, hot
start, graceful shutdown) and the three-phase configuration primitives that every other module builds on.
Wire it in as the first module: add `core` to the `go.work` use-list and `mise`'s `GO_MODULES`, confirm the
`ci.yml` matrix entry activates, and add `core/CHANGELOG.md` and `core/doc.go`. Start here next session
with `marathon start`.
