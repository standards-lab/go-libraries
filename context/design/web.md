# The HTTP layer

How the `web` package is shaped. The code and its `doc.go` are authoritative for the package surface;
this note holds the reasoning behind it and the intent for the parts not yet built.

## The bootstrap belongs in the library

The baseline shipped health handlers but no server, so every consumer hand-wrote the same wrapper around
`http.Server` — the same fifty lines in two demos, carrying the same defect: `ListenAndServe` binds inside
the goroutine it serves from, so a taken port or a bad address was logged by a goroutine nobody was
watching while startup carried on and readiness reported healthy with nothing listening.

`Server` splits that in two. The bind happens on the calling goroutine and its failure is a returned
error; serving happens in the background. Because the composition root registers `Start` as a lifecycle
startup hook, and a startup hook that cannot do its job fails the process, a bind failure now stops
startup. Binding first also makes the bound address knowable, so a caller can ask for port 0 and read back
what it got — which is what lets tests bind without racing a fixed port.

A serve failure after startup is a value, not a log line: it goes to a buffered channel the root selects
on and decides policy for. This keeps a logging dependency out of the base library while it has no
logging story, and a root that ignores the channel is no worse off than the baseline that logged and
carried on.

`web` registers no lifecycle hooks and holds no shutdown timeout. The coordinator cancels the root
context and then invokes each hook with a fresh timeout-bounded drain context, which `Shutdown(ctx)`
consumes directly — so the baseline's private `context.WithTimeout` and its `<-lc.Context().Done()` wait
have no counterpart here. `web` imports `lifecycle` only for the `ReadinessChecker` interface `/readyz`
consumes.

## One flat package

`web` is a single package, and stays one. In the baseline `web` was its own module and its concerns were
sub-packages, but that split was organizational: `problem`, `respond`, and `health` are all near-stdlib,
and none carries weight the rest of `web` shouldn't. So the rule that separates the base library from
provider sub-modules applies one level down — **a split is earned by dependency weight, not by topic**.
A sub-package appears only if some part of the HTTP layer needs a dependency the rest should not carry,
the same test that makes `database/postgres` a sub-module.

The cost is that names take the prefix a sub-package would have carried: `WriteProblem` rather than
`problem.Write`, and a future `CORSConfig` rather than `middleware.Config`, since `Config` is the
server's. That is ordinary Go, and cheaper than fragmenting a cohesive capability.

## Health is a surface, not a check

`/healthz` reports that the process is up and serving HTTP and checks nothing else — that is what makes an
unanswered probe the liveness signal rather than a 500. `/readyz` aggregates whatever
`lifecycle.ReadinessChecker` participants the composition root hands it; `web` contributes none of its
own. Today the only participant is the coordinator, so `/readyz` reports startup and draining; each
capability joins the list as it is built, named, so an operator can see which subsystem is holding
readiness down.

Readiness is non-monotonic because the coordinator's is: a draining process reports not-ready and stops
receiving traffic before its shutdown hooks run.

## The library defines no problem types

RFC 9457's `type` is the problem's identity — the member a client branches on, with `title` advisory and
`status` an advisory copy. A type URI therefore names an *application's* vocabulary, and a library that
mints one is claiming semantics it does not own. So every problem `web` emits is `about:blank`, which has
the defined meaning "no semantics beyond the HTTP status code", and consumers bring their own URIs
through `Problem.Write` or the extras map.

The trade-off this accepts: `/readyz` attaches a `checks` extension member to an `about:blank` problem,
and RFC 9457 means extension members to be defined by the problem type. If a consumer needs readiness
failures under its own vocabulary, the answer is a type hook on `Readiness`, not a library-owned URI —
deferred until something asks for it.

An empty title defaults to the status phrase. That is what the RFC asks of an `about:blank` problem, and
it keeps a hand-typed title from drifting away from the status code it accompanies.

## What the HTTP layer still needs

Deferred deliberately, each waiting on a consumer that would validate its shape:

- **Middleware** — a `func(http.Handler) http.Handler` chain composes outside the server today. The
  request logger waits on the logging story; `Auth`/`Authorize` wait on `auth` and `auth/authz`; CORS
  waits on a browser client.
- **Error mapping** — the domain-error-to-status matchers that turn a returned error into a problem
  response. Additive to the problem writers; needs a domain handler to exercise it.
- **The success envelope and the page response** — the shape a handler returns on success, and the
  HTTP-side pagination that pairs with `database`'s query vocabulary.

None of these revise the current surface; all of them add to it.
