# go-libraries

Go reference libraries for Standards Lab — layered building blocks for cloud-native enterprise
services.

The repository holds one base library: a single Go module rooted at
`github.com/standards-lab/go-libraries`. Each capability is a package inside that module, and all of them
are versioned and released together.

Providers that depend on a third-party SDK will be separate modules in nested directories, such as
`database/postgres` — each named for the system it targets and released on its own schedule, so a
consumer that needs one provider does not download the others. None exist yet. The base library takes at
most near-stdlib dependencies; today it has none.

## Packages

- `lifecycle` — starts subsystems concurrently, reports readiness, and shuts down within a timeout.
- `config` — loads configuration in layers: a base file, an environment overlay, and secrets.
- `logging` — builds an `*slog.Logger` from a configuration that `config` loads.
- `web` — an HTTP server, RFC 9457 problem responses, `/healthz` and `/readyz`, and a middleware chain.

## Development

The repository uses a Go workspace and [mise](https://mise.jdx.dev):

```
mise run test    # build and test every module
```

## License

[Apache License 2.0](LICENSE).
