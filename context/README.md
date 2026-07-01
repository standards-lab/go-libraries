# go-libraries

The standards-lab organization's Go reference libraries. This is Priority 2 of the reference-architecture
effort: a worked example of how to design, layer, and independently version shared libraries, built with
the marathon workflow so the workflow itself is battle-tested in the process.

The libraries are the capability boundary. Each capability is an independently versioned Go module; the
module is the unit of reuse. The standard first materializes here as code — the abstractions live in the
libraries, and the reference architecture (Priority 3) consumes them.

## What we're building toward

- The lowest practical level of abstraction, no frameworks by default. Dependencies flow downward only;
  interfaces are defined where they are consumed.
- Within a capability, implementations are providers selected by configuration: a self-hosted provider
  alongside managed ones behind one interface (the interface in the module root, each provider a nested
  submodule).
- Independent, per-module semantic-version releases.

## Capability map

Broad and shallow; detail is added when a module is about to be built.

- **core** — the foundation every other module builds on: the Layered Composition Architecture lifecycle
  (cold start, hot start, graceful shutdown) and the three-phase configuration primitives. First to be
  built; see `reset.md`.
- **auth** — authentication behind one interface, with providers (a self-hosted Keycloak provider
  alongside managed ones) as nested submodules. Authorization (RBAC/ABAC) as an in-module package.
- **database** — data access behind one interface, with SQL providers (postgres, sql server) as nested
  submodules.
- **storage** — object storage behind one interface (minio / azurite ↔ s3 / azure blob).
- **web** — the HTTP layer: RFC 9457 problem responses, middleware, the policy-enforcement point.

The set is a candidate, not a commitment — see `concepts/module-set.md`. Each module is settled when it is
built.

## How this repository works

- **Topology and naming** — the monorepo structure; the module-path, vendor, and tag conventions. See
  `design/library-topology-and-naming.md`.
- **Releases and CI** — per-module prefix tags, the matrix CI, the `go.work` and `mise` workflow. See
  `design/release-and-ci.md`.
- **Module conventions** — interface-in-root with vendor-in-submodule, explicit registration, co-located
  black-box tests, doc.go ownership. See `design/conventions.md`.
