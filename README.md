# go-libraries

Go reference libraries for the standards lab — layered, independently versioned building blocks for
modern cloud-native enterprise services.

A public, multi-module monorepo: each capability is its own Go module, versioned and released
independently, so a consumer can take one library without pulling the rest.

## Development

The repository uses a Go workspace and [mise](https://mise.jdx.dev):

```
mise run test    # build and test every module
```
