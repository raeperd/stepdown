# stepdown

[![Build Status](https://github.com/raeperd/stepdown/actions/workflows/build.yaml/badge.svg)](https://github.com/raeperd/stepdown/actions/workflows/build.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/raeperd/stepdown)](https://goreportcard.com/report/github.com/raeperd/stepdown)
[![Coverage Status](https://coveralls.io/repos/github/raeperd/stepdown/badge.svg?branch=main)](https://coveralls.io/github/raeperd/stepdown?branch=main)
[![Go Reference](https://pkg.go.dev/badge/github.com/raeperd/stepdown.svg)](https://pkg.go.dev/github.com/raeperd/stepdown)

Go linter that keeps your code reading top-to-bottom like a newsletter.

```
main.go:20:1: function "bar" is called by "foo" but declared before it (stepdown rule)
```

## What is the Stepdown Rule?

Robert C. Martin's *Clean Code* calls it the **Stepdown Rule**. Kent Beck calls it **Reading Order** in *Tidy First?*. Same idea — functions should be ordered so that each function appears above the functions it calls.

## Install

```bash
go install github.com/raeperd/stepdown/cmd/stepdown@latest
```

Or with golangci-lint:

```yaml
# .golangci.yml
linters:
  enable:
    - stepdown
```

## Run

```bash
stepdown ./...
# or
go vet -vettool=$(which stepdown) ./...
# or
golangci-lint run
```

## Configure

```yaml
# .golangci.yml
linters-settings:
  stepdown:
    exclusions:
      - "init"
      - "main"
```

## Contributing

```bash
make build  # Build binary
make test   # Run tests
make lint   # Run linter
```

## License

MIT
