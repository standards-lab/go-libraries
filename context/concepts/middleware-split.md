# When middleware becomes its own package

`design/web.md` settles that `web` is one flat package and that a split is earned by dependency weight,
not by topic. Middleware is the part of `web` most likely to test that rule, since it is the growth area:
the request logger is in, and CORS, `Auth`, `Authorize`, error mapping, a recovery handler, and a request
ID are all foreseeable. This note records what would actually earn the split, so a later session decides
on evidence rather than on the size of the file list.

## What does not earn it

- **Naming.** `web.RequestLogger` and a future `web.CORSConfig` read less well than `middleware.Logger`
  and `middleware.Config`. `design/web.md` weighed that cost when it settled the flat rule and accepted
  it; it names the middleware case explicitly. Ergonomics alone does not reopen the decision.
- **Count.** Ten middleware functions in one package is ordinary Go. `net/http` is larger.
- **`auth` and `auth/authz`.** The `Auth`/`Authorize` enforcement point makes `web` import `auth`, but
  `auth` is a near-stdlib interface package — the Keycloak and Entra SDKs live in sub-modules — so the
  weight a consumer picks up is close to nothing. `concepts/module-set.md` already places the enforcement
  point in `web`.

## What would earn it

A middleware whose dependency the rest of the HTTP layer should not carry. Concretely: a tracing
middleware pulling an OpenTelemetry SDK, or a rate limiter pulling a Redis client. That is the same test
that makes `database/postgres` a sub-module, applied one level down, and it is the trigger `design/web.md`
already describes in the abstract.

If it fires, the shape is `web/middleware` as a sub-package that imports `web` for problem writing — one
direction, no cycle, since `web` itself needs no middleware. A dependency heavy enough to belong in a
sub-module rather than a sub-package would go there instead, on the ordinary rule.

The move stays cheap while the library is pre-1.0, and `Middleware`/`Chain` and `RequestLogger` are kept
in `web/middleware.go` and `web/logger.go` so it would be a file move and a package clause rather than a
rewrite.

## Middleware belongs to the HTTP layer, not to the capability it logs or authenticates against

A related question, and this one is settled rather than open: a transport-agnostic capability does not
define HTTP middleware. `logging` does not ship the request logger, and `auth` will not ship the
authentication middleware. Both would invert the dependency direction — a CLI or a worker that wants a
`*slog.Logger` would compile `net/http` to get one — and both emit records or read headers in HTTP
vocabulary the transport layer owns. The capability supplies the collaborator; the HTTP layer supplies the
middleware that consumes it. This is `context/README.md`'s "dependencies flow downward only; interfaces
are defined where they are consumed" applied to the transport seam.
