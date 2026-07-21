# Module conventions

The patterns settled for this repository's libraries.

## Interface in a base package, vendor in a sub-module

A capability defines its interface (and shared types) as a package in the base library, using near-stdlib
dependencies only. Each concrete implementation whose weight comes from a third-party SDK lives in a
nested sub-module with its own `go.mod` that pins that SDK, so a consumer that needs only the interface
never pulls the SDKs. The `database` package defines the `database.DB` interface; `database/postgres` and
`database/mssql` are separate sub-modules pinning their respective drivers. The `auth` package defines
`Authenticator`/`TokenSource`; `auth/keycloak` and `auth/entra` are separate sub-modules.

## Near-stdlib base, heavy dependencies isolated in providers

The base library may depend only on packages as idiomatic and stable as the standard library itself —
`golang.org/x/…`, `google/uuid`, and the like. Heavy or vendor-specific dependencies — cloud SDKs,
database drivers — never enter the base; they live only in provider sub-modules. This keeps the base
effectively free to depend on: importing one capability package compiles no other capability and pulls no
vendor SDK.

## Providers selected by direct typed construction

Each provider exposes a `Provider` constant and a typed constructor. The application selects a provider at
its composition root with a typed switch over that constant and a direct import of the chosen provider —
the service knows its providers at compile time. There is no runtime registry, no `Register()` call, and
no `init()` side effects; importing a package never registers anything. Adding a provider is one new
import and one new switch case at the composition root, with no change to the capability package.

## One provider per target API, not per deployment

A provider is defined for a target API or system, not for a place it runs. The self-hosted ↔ managed seam
is a configuration change — endpoint, credentials, DSN — within a single provider wherever the API is
shared, and only spans providers when the APIs genuinely differ. Object storage is one provider per API
family (an S3-API provider serving a local emulator or a managed cloud; an Azure Blob provider likewise),
not one per environment. SQL is one provider per driver/dialect, each covering local and managed
deployments of that engine through its connection string.

## Capability interfaces are named per package

Each capability names its interface(s) for what reads best in its own package, rather than a forced
uniform noun. A capability that manages a lifecycle-bound client exposes a single encapsulated interface
(`database.DB`, `storage.Store`); a behavior-only capability exposes behavior interfaces
(`auth.Authenticator`, `auth.TokenSource`). Such an interface integrates with the `lifecycle` package at
the composition root — its start/stop wired through the coordinator, its readiness satisfying the
readiness contract — rather than re-declaring the lifecycle shape itself.

## Process lifecycle and context ownership

The `lifecycle` package coordinates startup, readiness, and graceful shutdown, and it fixes the
ecosystem's context-ownership convention. The composition root owns the root context — it traps signals
(`signal.NotifyContext`) and passes the context to `lifecycle.New`, which derives the cancellable context
every subsystem observes through `Coordinator.Context`. The coordinator installs no signal handlers of its
own, so a long-running service and a short-lived command share one context-derivation story.

Shutdown is two-phase and coordinator-driven: `Shutdown` cancels the root context, then invokes each
registered hook with a fresh, timeout-bounded drain context derived from `context.Background()`, so cleanup
is not pre-cancelled. A shutdown hook needs no cancellation guard of its own — it takes the drain context
and runs its graceful drain (`http.Server.Shutdown`, an in-flight wait) against it. Long-lived work that
must observe cancellation during operation watches `Coordinator.Context` instead.

Readiness is non-monotonic: the coordinator is ready once startup completes and not-ready again once
shutdown begins, so a `/readyz` probe reports a draining process as unavailable. A capability exposes its
own readiness through `lifecycle.ReadinessChecker`; its leaf subsystems take a plain `context.Context`
rather than the coordinator, keeping them usable without it.

## Tests: co-located and black-box

Tests are `{file}_test.go` files co-located with the source they cover, in an external test package
(`package <pkg>_test`). They exercise only the public API; private infrastructure is covered transitively
through the public entry points that use it.

## doc.go and godoc

Production source is written without doc comments; the agent writes godoc. Each package has exactly one
`doc.go` holding only the package comment.
