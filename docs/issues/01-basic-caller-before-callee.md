## Parent PRD

#1

## What to build

Detect when a callee function is declared before its caller within the same file. This is the core stepdown rule: callers above callees. Only plain functions (not methods) are checked in this slice.

## Implementation

All changes go in `analyzer.go` only — no new source files. Implement directly in the existing `run()` method.

Follow the recvcheck/whitespace pattern:
- Use `inspector.Preorder([]ast.Node{(*ast.FuncDecl)(nil)}, ...)` to collect all function declarations and their line positions per file
- Walk each function body to extract `*ast.CallExpr` → `*ast.Ident` for same-file calls
- Compare caller vs callee line positions; report via `pass.Reportf()` on the callee's `FuncDecl` position
- Group declarations by file (use `pass.Fset.Position().Filename`)

Algorithm sketch:
```
1. For each file, build map: funcName → line number (from *ast.FuncDecl)
2. For each function, walk body for *ast.CallExpr → *ast.Ident
3. If callee is in the map AND callee.line < caller.line → report
```

Once this lands, `go vet -vettool=./stepdown ./...` in CI will start self-linting this repo's own code.

## Acceptance criteria

- [ ] Plain function caller-before-callee violations are detected
- [ ] No false positives on correctly ordered code
- [ ] Calls to functions in other packages are ignored
- [ ] Calls to functions not declared in the same file are ignored
- [ ] Diagnostic message follows format: `function "bar" is called by "foo" but declared before it (stepdown rule)`
- [ ] `go vet -vettool=./stepdown ./...` works end-to-end
- [ ] All tests pass via `analysistest`

## TDD cycles

Each cycle is a commit: RED (failing test) → GREEN (minimal code to pass) → REFACTOR (clean up).

### Cycle 1: Simplest violation — single caller-callee pair

**RED** — Create `testdata/src/basic/basic.go`:
```go
package basic

func callee() {} // want `function "callee" is called by "caller" but declared before it \(stepdown rule\)`

func caller() {
	callee()
}
```
Update `analyzer_test.go` to run `analysistest.Run(t, ..., a, "basic")`. Test fails (no diagnostics emitted).

**GREEN** — In `analyzer.go` `run()`:
1. Collect `funcName → position` map from `*ast.FuncDecl` nodes
2. Walk each function body for `*ast.CallExpr` → `*ast.Ident`
3. If callee exists in map and callee line < caller line → `pass.Reportf()`

**REFACTOR** — Clean up if needed.

### Cycle 2: No false positives on valid code

**RED** — Update `testdata/src/valid/valid.go`:
```go
package valid

func caller() {
	callee()
}

func callee() {}
```
Test should already pass (no `// want` comments = no diagnostics expected). Verify.

### Cycle 3: External package calls ignored

**RED** — Create `testdata/src/basic/external.go`:
```go
package basic

import "fmt"

func printer() {
	fmt.Println("hello")
}
```
Verify no false positives on cross-package calls. If the implementation already handles this (only matching names in the local map), this passes immediately.

### Cycle 4: Transitive violations

**RED** — Create `testdata/src/basic/transitive.go`:
```go
package basic

func c() {} // want `function "c" is called by "b" but declared before it \(stepdown rule\)`

func b() { // want `function "b" is called by "a" but declared before it \(stepdown rule\)`
	c()
}

func a() {
	b()
}
```

**GREEN** — Should already work if each caller-callee pair is checked independently. If not, fix.

### Cycle 5: Multiple callees, only some violating

**RED** — Create `testdata/src/basic/multiple.go`:
```go
package basic

func bottomCallee() {} // want `function "bottomCallee" is called by "topCaller" but declared before it \(stepdown rule\)`

func topCaller() {
	bottomCallee()
	middleHelper()
}

func middleHelper() {}
```

**GREEN** — Should already work. Verify only the violating callee is reported.

## Blocked by

None — can start immediately

## User stories addressed

- User story 1: caller-before-callee order checking
- User story 5: `go vet -vettool` usage
- User story 8: clear violation messages
- User story 10: CI catches violations
- User story 11: standard `analysis.Analyzer` interface
