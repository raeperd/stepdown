# stepdown

[![Build Status](https://github.com/raeperd/stepdown/actions/workflows/build.yaml/badge.svg)](https://github.com/raeperd/stepdown/actions/workflows/build.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/raeperd/stepdown)](https://goreportcard.com/report/github.com/raeperd/stepdown)
[![Coverage Status](https://coveralls.io/repos/github/raeperd/stepdown/badge.svg?branch=main)](https://coveralls.io/github/raeperd/stepdown?branch=main)
[![Go Reference](https://pkg.go.dev/badge/github.com/raeperd/stepdown.svg)](https://pkg.go.dev/github.com/raeperd/stepdown)

Go linter that checks callers are declared before callees — the **Stepdown Rule** from Clean Code.

## Why

Robert C. Martin's *Clean Code* introduces the **Stepdown Rule**: functions should be ordered so that callers appear before callees, and code reads top-to-bottom like a narrative. Kent Beck echoes this as **Reading Order** in *Tidy First?*.

No Go linter currently enforces this. `stepdown` fills that gap.

## Installation

```bash
# Standalone
go install github.com/raeperd/stepdown/cmd/stepdown@latest

# With golangci-lint (recommended)
# Add to .golangci.yml:
linters:
  enable:
    - stepdown
```

## Usage

```bash
stepdown ./...
# or
go vet -vettool=$(which stepdown) ./...
# or
golangci-lint run
```

Output:
```
main.go:20:1: function "bar" is called by "foo" but declared before it (stepdown rule)
```

## Configuration

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
make build # Build binary
make test  # Run tests
make lint  # Run linter
```

## License

MIT
