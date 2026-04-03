## Parent PRD

#1

## What to build

Wire up the `Settings.Exclusions` field so users can exclude specific function names from stepdown checks. Common exclusions: `init`, `main`, `TestXxx`. Excluded functions are skipped both as callers and callees.

## Implementation

All changes go in `analyzer.go` only — wire `Settings.Exclusions` into the existing `run()`. No new source files.

Follow recvcheck's exclusion pattern:
- Store exclusions as `map[string]struct{}` on the `analyzer` struct (built in `NewAnalyzer()`)
- Before reporting a violation, check if either the caller or callee name is in the exclusion set
- Wire exclusions into `cmd/stepdown/main.go` via a flag (comma-separated list)

## Acceptance criteria

- [ ] Functions listed in `Settings.Exclusions` produce no violations as callers or callees
- [ ] Default behavior (no exclusions) is unchanged
- [ ] Exclusions are configurable via CLI flag and golangci-lint `linters-settings`

## TDD cycles

Each cycle is a commit: RED (failing test) → GREEN (minimal code to pass) → REFACTOR (clean up).

### Cycle 1: Excluded caller produces no violation

**RED** — Create `testdata/src/exclusions/exclusions.go`:
```go
package exclusions

func setup() {}

func init() {
	setup()
}
```
Add test case with `Settings{Exclusions: []string{"init"}}`. Without exclusion support, this would report `setup` as violation. Test should pass with no diagnostics.

**GREEN** — In `NewAnalyzer()`, build exclusion map from `Settings.Exclusions`. In `run()`, skip reporting when caller or callee is in the exclusion set.

**REFACTOR** — Clean up.

### Cycle 2: Non-excluded functions still checked

**RED** — Create `testdata/src/exclusions/non_excluded.go`:
```go
package exclusions

func callee() {} // want `function "callee" is called by "caller" but declared before it \(stepdown rule\)`

func caller() {
	callee()
}
```
Same test settings `{Exclusions: []string{"init"}}`. Verify non-excluded functions still produce violations.

### Cycle 3: Multiple exclusions

**RED** — Create `testdata/src/exclusions/main_excluded.go`:
```go
package exclusions

func run() {}

func main() {
	run()
}
```
Update test settings to `{Exclusions: []string{"init", "main"}}`. Verify no diagnostics.

## Blocked by

- Blocked by issue 01 (basic caller-before-callee)

## User stories addressed

- User story 7: exclude specific functions from checks
