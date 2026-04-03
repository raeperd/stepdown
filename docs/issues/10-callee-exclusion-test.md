## Parent PRD

#1

## What to build

Add a test that exercises the callee exclusion path — when a non-excluded function calls an excluded function that is declared before it, no violation should be reported.

The code path exists (`analyzer.go` line 173) but is untested. Since `init` can't be called explicitly in Go, a new test package with a different exclusion name (e.g., `"excluded"`) is needed.

## Implementation

No changes to `analyzer.go`. Add a new test function in `analyzer_test.go` and a new testdata package.

## Acceptance criteria

- [ ] Callee exclusion path is exercised and verified
- [ ] Caller exclusion still works
- [ ] Non-excluded pairs still produce violations

## TDD cycles

### Cycle 1: Callee exclusion test

**RED** — Create `testdata/src/callee_exclusion/callee_exclusion.go`:
```go
package callee_exclusion

// No violation — "excluded" is excluded as callee
func excluded() {}

func caller() {
	excluded()
}

// Non-excluded pair still produces violation
func other() {} // want `function "other" is called by "another" but declared before it \(stepdown rule\)`

func another() {
	other()
}
```
Add `TestAnalyzerCalleeExclusion` with `Settings{Exclusions: []string{"excluded"}}`. Test should pass immediately since the code path already works.

If it fails, the callee exclusion logic has a bug.
