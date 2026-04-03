## Parent PRD

#1

## What to build

Enforce that callees are declared in the same order they are invoked. If `foo()` calls `bar()` then `baz()`, then `bar` should be declared before `baz`. This ensures reading the callee declarations mirrors the caller's execution flow.

## Implementation

All changes go in `analyzer.go` only — extend the existing call collection. No new source files.

Changes to the existing algorithm:
- When walking a function body for calls, record the **invocation order** (slice of callee names, preserving first-occurrence order)
- After collecting, compare invocation order against declaration order of those callees
- If callee B is invoked before callee A, but A is declared before B → report on B
- Use a different diagnostic message to distinguish from caller-before-callee violations

Diagnostic message format:
```
function "second" is called by "main" before "first" but declared after it (stepdown rule)
```

## Acceptance criteria

- [ ] Callee declaration order violations are detected when callees are declared in a different order than their invocation order within the caller
- [ ] Only the first invocation of each callee determines its position (duplicate calls are ignored)
- [ ] Plain function and method callee ordering are both checked
- [ ] No false positives when callees are declared in invocation order

## TDD cycles

Each cycle is a commit: RED (failing test) → GREEN (minimal code to pass) → REFACTOR (clean up).

### Cycle 1: Basic invocation order violation

**RED** — Create `testdata/src/calleeorder/calleeorder.go`:
```go
package calleeorder

func main() {
	first()
	second()
}

func second() {} // want `function "second" is called by "main" before "first" but declared after it \(stepdown rule\)`

func first() {}
```
Add test case in `analyzer_test.go`. Test fails.

**GREEN** — In `run()`, track invocation order per caller. After collecting all calls, for each pair of callees (A invoked before B), check if B is declared before A. Report on B.

**REFACTOR** — Clean up invocation order tracking.

### Cycle 2: Valid invocation order

**RED** — Create `testdata/src/calleeorder/valid.go`:
```go
package calleeorder

func run() {
	step1()
	step2()
	step3()
}

func step1() {}
func step2() {}
func step3() {}
```
Verify no diagnostics.

### Cycle 3: Duplicate calls don't affect ordering

**RED** — Create `testdata/src/calleeorder/duplicate_calls.go`:
```go
package calleeorder

func process() {
	validate()
	transform()
	validate() // second call — should not affect ordering
}

func validate() {}
func transform() {}
```
Verify no diagnostics (invocation order matches declaration order by first occurrence).

## Blocked by

- Blocked by issue 01 (basic caller-before-callee)

## User stories addressed

- User story 2: callee ordering enforced to match invocation order
