# stepdown

[![Build Status](https://github.com/raeperd/stepdown/actions/workflows/build.yaml/badge.svg)](https://github.com/raeperd/stepdown/actions/workflows/build.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/raeperd/stepdown)](https://goreportcard.com/report/github.com/raeperd/stepdown)
[![Coverage Status](https://coveralls.io/repos/github/raeperd/stepdown/badge.svg?branch=main)](https://coveralls.io/github/raeperd/stepdown?branch=main)
[![Go Reference](https://pkg.go.dev/badge/github.com/raeperd/stepdown.svg)](https://pkg.go.dev/github.com/raeperd/stepdown)

Go linter that keeps your code reading top-to-bottom like a newsletter.

```
main.go:20:1: function "bar" is called by "foo" but declared before it (stepdown rule)
```

## What is the stepdown rule?

Robert C. Martin's *Clean Code* calls it the **Stepdown Rule**. Kent Beck calls it **Reading Order** in *Tidy First?*. Both describe the same idea: functions should be ordered so that each function appears above the functions it calls.

## Install

```bash
go install github.com/raeperd/stepdown/cmd/stepdown@latest
```

## Run

```bash
stepdown ./...
# or
go vet -vettool=$(which stepdown) ./...
```

## Use with golangci-lint

stepdown is available as a module plugin. Add it to `.custom-gcl.yml`:

```yaml
version: v2.13.1
plugins:
  - module: github.com/raeperd/stepdown
    import: github.com/raeperd/stepdown/plugin
    version: v0.1.0
```

Build a custom golangci-lint binary:

```bash
golangci-lint custom
```

Register and configure stepdown in `.golangci.yml`:

```yaml
version: "2"
linters:
  enable:
    - stepdown
  settings:
    custom:
      stepdown:
        type: module
        description: Checks that callers are declared before callees.
        settings:
          exclusions:
            - init
            - main
```

Run the custom binary:

```bash
./custom-gcl run
```

## Development

AI tools were used during development. I personally reviewed every line of code in this repository, understand how it works, and take responsibility for its design, implementation, and maintenance.

## Contributing

```bash
make build  # Build binary
make test   # Run tests
make lint   # Run linter
```

## License

MIT
