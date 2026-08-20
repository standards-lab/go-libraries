# go-libraries

**This repository is archived.** Its capabilities moved to dedicated repositories, one per
capability, so each is versioned and released on its own. The history and tags here are retained,
and released `go-libraries` module versions remain fetchable through the Go proxy, but nothing
further is released from this module.

## Where each package went

| Package here | Successor |
|---|---|
| `config`, `lifecycle`, `logging` | [`github.com/standards-lab/go-core`](https://github.com/standards-lab/go-core) |
| `database`, `database/seed`, `database/postgres` | [`github.com/standards-lab/go-database`](https://github.com/standards-lab/go-database) (`postgres` stays a nested sub-module) |
| `web` | [`github.com/standards-lab/go-web-sdk`](https://github.com/standards-lab/go-web-sdk) (middleware implementations now live in its `middleware` package) |

Each successor started from a snapshot of this repository with its module path rewritten; blame
begins at that import commit, and this repository keeps the record that precedes it.

## License

[Apache License 2.0](LICENSE).
