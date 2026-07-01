# Module conventions

The patterns settled for this repository's modules.

## Interface in root, vendor in submodule

A capability whose implementations carry heavy third-party dependencies defines its interface (and shared
types and registry) at the module root, using the standard library only. Each concrete implementation
lives in a nested submodule with its own `go.mod` that pins the vendor SDK, so a consumer that needs only
the interface never pulls the SDKs. Example: the `auth` module root defines the authenticator interface;
`auth/keycloak` and `auth/entra` are separate submodules pinning their respective SDKs.

## Explicit registration, no init() side effects

Each implementation exposes a public `Register()` that registers its factory with the root module's
registry. The application calls `Register()` explicitly at its composition root. Registration never
happens as an `init()` side effect of importing a package.

## Tests: co-located and black-box

Tests are `{file}_test.go` files co-located with the source they cover, in an external test package
(`package <pkg>_test`). They exercise only the public API; private infrastructure is covered transitively
through the public entry points that use it.

## doc.go and godoc

Production source is written without doc comments; the agent writes godoc. Each package has exactly one
`doc.go` holding only the package comment.
