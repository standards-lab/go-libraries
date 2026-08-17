# go-libraries

The Standards Lab organization's Go reference libraries. This is the library level of the organization's
reference architecture: a worked example of how to design, layer, and independently version shared
libraries, built with the marathon workflow so the workflow itself is exercised in the process.

The libraries are the capability boundary: the standard first appears here as code, and the abstractions
live here.

## What we're building toward

- The lowest practical level of abstraction, no frameworks by default. Dependencies flow downward only;
  interfaces are defined where they are consumed. A consumer interfaces with a technology at the
  resolution its purpose requires, and the provider's own API stays reachable through the handle the
  capability exposes.
- Every capability presents two tiers: a standard tier that is exactly the technology's common standard,
  as a package in the base library, and the provider's native API, in a nested sub-module that pins its
  own SDK. A provider is one implementation of one target API, self-hosted or managed by configuration;
  the consumer selects it by typed construction at compile time and imports it only from its composition
  root and the packages that declare native use. Each capability declares what a provider swap costs. See
  `design/conventions.md`.
- A single base library, versioned and released as one artifact, plus provider sub-modules released
  independently.

## Repository topology

The repository is one base library — a single Go module rooted here — plus a set of provider sub-modules.
Every capability is a package inside the base library; the packages co-evolve and release together. The
base library takes only near-stdlib dependencies, so importing one capability pulls in no heavy SDKs and
compiles no other capability. A provider whose weight comes from a third-party SDK is a nested sub-module
with its own `go.mod`, versioned on its own schedule and selected by the consumer without pulling the
others. See `design/library-topology-and-naming.md`.

## Capability map

Broad and shallow; detail is added when a capability is about to be built.

Base library packages:

- **lifecycle** — process lifecycle for long-running consumers: concurrent startup, a readiness signal,
  and timeout-bounded graceful shutdown. Conceptually, cold start initializes objects from configuration
  so their state is valid (the `/healthz` side), and hot start brings the long-running services up until
  they are ready to receive requests (the `/readyz` side). Built; the code and its `doc.go` are
  authoritative.
- **config** — layered configuration: a base file, environment overlays, and `secrets.json`, resolved
  through a merge/finalize contract each capability's config implements.
- **logging** — the `*slog.Logger` a process writes through, built from a configuration that takes part in
  the layered load. It constructs a logger and nothing more: the level vocabulary is `slog`'s, and the
  HTTP request logger belongs to `web`. See `design/logging.md`.
- **auth** — authentication behind `Authenticator`/`TokenSource` interfaces; the standard tier is OAuth
  2.0, OpenID Connect, and JWT, with providers (Keycloak, Entra) as nested sub-modules, each usable
  locally or managed. Token verification is interchangeable; what a token's claims contain is
  interchangeable with review. Authorization (ABAC/RBAC) as an in-package `auth/authz`.
- **database** — SQL data access: the `database.DB` wrapper, the dialect interface, and a persistence
  query vocabulary in the base, with the standard tier ISO/IEC 9075 and drivers as nested provider
  sub-modules — `database/postgres` built, a second engine when a proof or a consumer earns it. Schema-
  bearing: what changes on a provider swap is the consumer's schema, migrations, and domain SQL, never
  the composition root. See `concepts/database.md` while it is being built.
- **storage** — object storage behind the `storage.Store` interface; no formal standard exists, so the
  standard tier is the minimal operation set common to the target APIs, with providers per API family (an
  S3-API provider, an Azure Blob provider) as nested sub-modules, each serving a local emulator or a
  managed cloud by configuration. Those operations are interchangeable; consistency is interchangeable
  with review.
- **web** — the HTTP layer, its standard tier RFC 9110 and RFC 9457 over the stdlib `net/http` transport,
  with no providers: a server, problem responses, a success envelope, middleware, liveness/readiness
  (`/healthz`, `/readyz`, where `/readyz` reports the `lifecycle` readiness signal), and the
  authorization enforcement point.
- **observability** — logs, metrics, and traces across the stack, with OpenTelemetry as the standard
  tier; exporters as providers. Unbuilt; `logging` covers the `*slog.Logger` until then.

The set is provisional, not a commitment — see `concepts/module-set.md`. Each capability is settled when
it is built. `lifecycle`, `config`, `logging`, `database` with its Postgres provider and `database/seed`,
and a `web` with the bootstrap, the routing layer, the health endpoints, and the middleware chain are in;
the remaining capabilities follow as they are reached, and providers are scaffolded only when built — one
per capability by default, a second when a proof or a consumer earns it, never a matrix.

## How this repository works

- **Topology and naming** — the base-library-plus-provider-sub-module structure; the module-path, naming,
  and tag conventions. See `design/library-topology-and-naming.md`.
- **Releases and CI** — the base library's root tags and each provider's prefix tags, the matrix CI, the
  `go.work` and `mise` workflow. See `design/release-and-ci.md`.
- **Module conventions** — the standard and native tiers, interface-in-a-base-package with
  vendor-in-sub-module, the near-stdlib base, the import boundary, providers selected by typed
  construction, provider-swap classes, the dialect interface's growth rule, co-located black-box tests,
  doc.go ownership. See `design/conventions.md`.
