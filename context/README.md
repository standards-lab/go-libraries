# go-libraries

The standards-lab organization's Go reference libraries. This is the library level of the organization's
reference architecture: a worked example of how to design, layer, and independently version shared
libraries, built with the marathon workflow so the workflow itself is battle-tested in the process.

The libraries are the capability boundary. The standard first materializes here as code; the abstractions
live in the libraries.

## What we're building toward

- The lowest practical level of abstraction, no frameworks by default. Dependencies flow downward only;
  interfaces are defined where they are consumed.
- Within a capability, implementations are providers selected by configuration: a self-hosted provider
  alongside managed ones behind one interface. The interface is a package in the base library; each
  provider is a nested sub-module that pins its own SDK and is selected by the consumer at compile time.
- A single base library, versioned and released as one artifact, plus provider sub-modules released
  independently.

## How the repository is shaped

The repository is one base library — a single Go module rooted here — plus a set of provider sub-modules.
Every capability is a package inside the base library; the packages co-evolve and release together. The
base library takes only near-stdlib dependencies, so importing one capability pulls in no heavy SDKs and
compiles no other capability. A provider whose weight comes from a third-party SDK is a nested sub-module
with its own `go.mod`, versioned on its own schedule and selected by the consumer without pulling the
others. See `design/library-topology-and-naming.md`.

## Capability map

Broad and shallow; detail is added when a capability is about to be built.

Base library packages:

- **lifecycle** — the process-lifecycle foundation every long-running consumer builds on: concurrent
  startup, a readiness signal, and timeout-bounded graceful shutdown. Its phases are cold start (boot and
  registration), hot start (subsystems warming while readiness is still false), and graceful shutdown.
  First to be built; see `reset.md`.
- **config** — layered configuration: a base file, environment overlays, and `secrets.json`, resolved
  through a merge/finalize contract each capability's config implements.
- **auth** — authentication behind `Authenticator`/`TokenSource` interfaces, with providers (a
  self-hosted Keycloak provider alongside managed ones) as nested sub-modules. Authorization (ABAC/RBAC)
  as an in-package `auth/authz`.
- **database** — SQL data access behind the `database.DB` interface, with a persistence query vocabulary;
  drivers (postgres, sql server) as nested provider sub-modules.
- **storage** — object storage behind the `storage.Store` interface; providers per API family (an S3-API
  provider, an Azure Blob provider) as nested sub-modules, each serving a local emulator or a managed
  cloud by configuration.
- **web** — the HTTP layer: a stdlib `net/http` server, RFC 9457 problem responses, a success envelope,
  middleware, liveness/readiness (`/healthz`, `/readyz`, where `/readyz` surfaces the `lifecycle`
  readiness signal), and the authorization enforcement point.

The set is provisional, not a commitment — see `concepts/module-set.md`. Each capability is settled when
it is built. Build order starts with `lifecycle`, then `config`, then a minimal `web`; the remaining
capabilities follow as they are reached, and providers are scaffolded only when built.

## How this repository works

- **Topology and naming** — the base-library-plus-provider-sub-module structure; the module-path, naming,
  and tag conventions. See `design/library-topology-and-naming.md`.
- **Releases and CI** — the base library's root tags and each provider's prefix tags, the matrix CI, the
  `go.work` and `mise` workflow. See `design/release-and-ci.md`.
- **Module conventions** — interface-in-a-base-package with vendor-in-sub-module, the near-stdlib base,
  providers selected by typed construction, co-located black-box tests, doc.go ownership. See
  `design/conventions.md`.
