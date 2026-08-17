# The HTTP layer

The design of the `web` package. The code and its `doc.go` are authoritative for the package API;
this note holds the reasoning behind it and the intent for the parts not yet built.

## The bootstrap belongs in the library

Earlier implementations of this layer shipped health handlers but no server, so every consumer hand-wrote
the same wrapper around `http.Server` — the same fifty lines in every service, with the same defect:
`ListenAndServe` binds inside
the goroutine it serves from, so a taken port or a bad address was logged by a goroutine nobody was
watching while startup carried on and readiness reported healthy with nothing listening.

`Server` splits that in two. The bind happens on the calling goroutine — bounded by the context `Start`
takes — and its failure is a returned error; serving happens in the background. Because the composition
root registers `Start` as a lifecycle startup hook and a startup hook's error fails the coordinator's
startup, a bind failure stops startup before readiness ever flips. Binding first also makes the bound
address knowable, so a caller can ask for port 0 and read back what it got — which is what lets tests
bind without racing a fixed port.

A serve failure after startup is a value, not a log line: it goes to a buffered channel the root hands
to the coordinator as a monitored source, so the first failure ends the run and surfaces in `Run`'s
returned error. That decision holds now that `logging` exists: a library that logged the failure would
be choosing a policy the composition root owns.

`web` registers no lifecycle hooks and holds no shutdown timeout, but `Start` and `Shutdown` carry the
lifecycle hook signature, so the root registers them as bare method values. The coordinator cancels the
run context and then invokes each hook with a fresh timeout-bounded drain context, which `Shutdown(ctx)`
consumes directly — so a hand-rolled private `context.WithTimeout` and a cancellation guard have no
counterpart here. `web` imports `lifecycle` only for the `ReadinessChecker` interface `/readyz`
consumes.

## One flat package

`web` is a single package, and stays one. An earlier design made `web` its own module with its concerns as
sub-packages, but that split was organizational: `problem`, `respond`, and `health` are all near-stdlib,
and none carries weight the rest of `web` shouldn't. So the rule that separates the base library from
provider sub-modules applies one level down — a split is earned by dependency weight, not by topic.
A sub-package appears only if some part of the HTTP layer needs a dependency the rest should not carry,
the same test that makes `database/postgres` a sub-module.

The cost is that names take the prefix a sub-package would have carried: `WriteProblem` rather than
`problem.Write`, and a future `CORSConfig` rather than `middleware.Config`, since `Config` is the
server's. That is ordinary Go, and cheaper than fragmenting a cohesive capability.

## Health reports; it does not check

`/healthz` reports that the process is up and serving HTTP and checks nothing else — that is what makes an
unanswered probe the liveness signal rather than a 500. `/readyz` aggregates whatever
`lifecycle.ReadinessChecker` participants the composition root hands it; `web` contributes none of its
own. Today the only participant is the coordinator, so `/readyz` reports startup and draining; each
capability joins the list as it is built, named, so an operator can see which subsystem is holding
readiness down.

Readiness is non-monotonic because the coordinator's is: a draining process reports not-ready and stops
receiving traffic before its shutdown hooks run.

## The library defines no problem types

The package's standard tier is RFC 9110 and RFC 9457 over the stdlib transport, and it has no providers:
nothing changes on a provider swap because there is nothing to swap. The problem writers are the clearest
case of staying inside the standard.

RFC 9457's `type` is the problem's identity — the member a client branches on, with `title` advisory and
`status` an advisory copy. A type URI therefore names an *application's* vocabulary, and a library that
mints one is claiming semantics it does not own. So every problem `web` emits is `about:blank`, which has
the defined meaning "no semantics beyond the HTTP status code", and consumers bring their own URIs
through `Problem.Write` or the extras map.

The trade-off this accepts: `/readyz` attaches a `checks` extension member (RFC 9457's term) to an
`about:blank` problem,
and RFC 9457 means extension members to be defined by the problem type. If a consumer needs readiness
failures under its own vocabulary, the answer is a type hook on `Readiness`, not a library-owned URI —
deferred until something asks for it.

An empty title defaults to the status phrase. That is what the RFC asks of an `about:blank` problem, and
it keeps a hand-typed title from drifting away from the status code it accompanies.

## Middleware

`Middleware` is `func(http.Handler) http.Handler` — the signature the ecosystem already uses — and `Chain`
composes a set of them in argument order, so the first argument sees the request first. Both live in
`web` under the one-flat-package rule; `concepts/middleware-split.md` records what would earn them a
package of their own.

A middleware belongs to the transport that consumes it, not to the capability it collaborates with. The
request logger is in `web` and takes a stdlib `*slog.Logger` rather than living in `logging` and taking a
`Middleware`; the same reasoning will keep `Auth` here rather than in `auth`. The alternative inverts the
dependency direction — a worker that wants a logger would compile `net/http` to get one — and the record
the request logger emits is HTTP vocabulary in any case.

Two things earlier request loggers got wrong, both from wrapping the `ResponseWriter`. They swallowed
a second `WriteHeader` instead of delegating it, hiding the standard library's superfluous-header warning;
and it implemented no `Unwrap`, so wrapping silently cost a handler flushing and hijacking. The wrapper
here records the first status, always delegates, and implements `Unwrap` so `http.ResponseController`
reaches through it. Seeding the recorded status with 200 covers the handler that writes a body without
calling `WriteHeader`, which removes the need to intercept `Write` at all.

A request logs at info on the normal path, with one carve-out: a successful probe request logs at
debug, because an orchestrator hits `/healthz` and `/readyz` every few seconds forever and that
heartbeat would otherwise dominate production logs — while a failing probe stays at info, since
readiness flapping is exactly what an operator greps for. A panicking handler logs its record at error
level with the panic value attached, then the panic continues to net/http's recovery. Beyond the probe
carve-out the middleware does not judge status codes: treating a 5xx as an error record means knowing
whether the status came from the application's own failure, and that judgment belongs to the error
mapping below, not to the middleware. The duration attribute is a `slog.Duration`, which the JSON handler renders as nanoseconds; a
`duration_ms` float would suit dashboards better and waits for an observability consumer to ask.

## Route groups and modules

The routing layer is three types — `Group`, `Module`, `Router` — re-derived from the prior reference
material (ref-go-libraries' module package, exercised by two consumer services) with its catalogued
flaws corrected. The code and `doc.go` are authoritative for the API; this note holds the decisions.

- **Compose once, at wiring time.** `NewModule` compiles a group tree into full mux patterns with
  middleware baked in; nothing recomposes per request (the reference rebuilt the middleware chain on
  every request) and nothing rewrites request paths except `NewHandlerModule`'s stdlib `StripPrefix`
  (the reference mutated the caller's request). Compile-once creates its own hazard — a route added
  after compilation would be silently dead — so compiling seals the group and later mutation panics.
- **One canonical path.** A group becomes servable one way (`NewModule`) and mounts one way
  (`Router.Mount`); the reference carried a flatten-vs-bridge duality, and its single-level prefix
  restriction meant `/api/v1` existed only through a workaround. Multi-level prefixes are first-class
  here.
- **Probes stay outside the modules structurally.** `Router.Handle` mirrors `ServeMux.Handle` on the
  router's native fallback mux, satisfying `Mounter`, so `RegisterHealth` mounts the probes beyond
  every module's middleware: they need no auth exemption because they never pass through a module.
  Router-level middleware (the request logger) still wraps the whole dispatch.
- **Registration mistakes panic at wiring time** — a malformed prefix, a duplicate pattern or mount, a
  sealed-group mutation — mirroring lifecycle's misuse panics: the composition root is written once
  and runs at boot, so a loud failure there beats a silent one in production.

## What the HTTP layer still needs

Deferred deliberately, each waiting on a consumer that would validate its API:

- **The rest of the middleware set** — `Auth`/`Authorize` wait on `auth` and `auth/authz`; CORS waits on a
  browser client; a recovery handler and a request ID wait for a service to need them.
- **Error mapping** — the domain-error-to-status matchers that turn a returned error into a problem
  response. Additive to the problem writers; needs a domain handler to exercise it.
- **The success envelope and the page response** — the JSON structure a handler returns on success, and
  the HTTP-side pagination that pairs with `database`'s query vocabulary. The page contract carries no
  engine detail: the query builder renders paging in the standard SQL form, so nothing a dialect does
  reaches the HTTP side.

None of these revise the current API; all of them add to it.
