# reset · library-level-web-baseline

- **Status:** closeout
- **Session:** plan
- **Branch:** library-level-web-baseline

## Disposition

- **Relabeled to the library level** (`context/README.md`, `CLAUDE.md`): this repository is the library
  level of the organization's reference architecture, replacing the earlier "Priority 2" framing. The
  context stays self-contained — it describes the libraries and their own boundaries, not the consumers
  above them.
- **`web` module** (`context/README.md`, `concepts/module-set.md`): the HTTP layer now names a stdlib
  `net/http` server and liveness/readiness (`/healthz`, `/readyz`, where `/readyz` surfaces `core`'s
  hot-start lifecycle) alongside RFC 9457 responses, middleware, and the enforcement point.
- **Resolved the `web` split** (`concepts/module-set.md`): the `net/http` bootstrap, problem-response
  scaffolding, middleware, and the probes are library-level; request routing and domain handlers belong to
  the consuming service. Recorded the build order — `core`, then a minimal `web` — with the set growing one
  module at a time as each pattern and integration is established.

## Next-focus

Build the `core` foundation module — the Layered Composition Architecture lifecycle (cold start, hot start,
graceful shutdown) and the three-phase configuration primitives that every other module builds on. Wire it
in as the first module: add `core` to the `go.work` use-list and `mise`'s `GO_MODULES`, confirm the
`ci.yml` matrix entry activates, and add `core/CHANGELOG.md` and `core/doc.go`. Then build the minimal
`web` — the `net/http` server with `/healthz` and `/readyz` — as the smallest runnable HTTP surface. Start
here next session with `marathon start`.
