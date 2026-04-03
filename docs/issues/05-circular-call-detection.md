## Parent PRD

#1

## What to build

Detect circular call chains (direct and indirect) and skip reporting violations for those pairs. If A calls B and B calls A, neither ordering satisfies the stepdown rule, so both directions should be silently ignored.

## Implementation

All changes go in `analyzer.go` only — add cycle detection after collecting the call graph. No new source files.

Changes to the existing algorithm:
- After collecting all caller→callee edges per file, build a directed graph
- Before reporting a violation for pair (caller, callee), check if there's a path from callee back to caller
- Simple approach: for each violation candidate, do a DFS/BFS from callee in the call graph; if caller is reachable, skip
- The call graph is small (per-file), so performance is not a concern

## Acceptance criteria

- [ ] Direct circular calls (A→B, B→A) produce no violations
- [ ] Indirect circular calls (A→B→C→A) produce no violations for any pair in the cycle
- [ ] Non-circular pairs in the same file are still reported
- [ ] Circular detection works for both plain functions and methods

## TDD cycles

Each cycle is a commit: RED (failing test) → GREEN (minimal code to pass) → REFACTOR (clean up).

### Cycle 1: Direct circular calls ignored

**RED** — Create `testdata/src/circular/direct.go`:
```go
package circular

func a() {
	b()
}

func b() {
	a()
}
```
Add test case in `analyzer_test.go`. Test fails (currently reports `b` as violation because it's declared after `a` but `a` also calls `b` — wait, `a` is before `b` so no violation from issue 01. Actually `b` calls `a` and `a` is before `b`, so no violation either. Let me reconsider.)

Actually: `a` calls `b` (b is after a — valid), `b` calls `a` (a is before b — violation). Without cycle detection, this reports a violation on `a`. With cycle detection, it should be suppressed.

Corrected test:
```go
package circular

// No violations — a and b call each other
func a() {
	b()
}

func b() {
	a()
}
```
Currently reports: `function "a" is called by "b" but declared before it (stepdown rule)`. After this issue, no violation.

**GREEN** — Before reporting, check if callee also (directly or indirectly) calls caller. If yes, skip.

**REFACTOR** — Extract cycle detection into a helper.

### Cycle 2: Indirect circular calls ignored

**RED** — Create `testdata/src/circular/indirect.go`:
```go
package circular

func x() {
	y()
}

func y() {
	z()
}

func z() {
	x()
}
```
`z` calls `x`, and `x` is before `z` — violation without cycle detection. With cycle detection (x→y→z→x), suppressed.

**GREEN** — Extend cycle detection to handle indirect cycles via DFS.

### Cycle 3: Mixed — circular pair + non-circular pair

**RED** — Create `testdata/src/circular/mixed.go`:
```go
package circular

func helper() {} // want `function "helper" is called by "entry" but declared before it \(stepdown rule\)`

func entry() {
	helper()
	mutual()
}

func mutual() {
	entry()
}
```
`helper` is called by `entry` but declared before it — not in a cycle, so reported. `entry` and `mutual` form a cycle — suppressed.

**GREEN** — Should work if cycle detection is per-pair.

## Blocked by

- Blocked by issue 01 (basic caller-before-callee)

## User stories addressed

- User story 4: circular calls automatically ignored
