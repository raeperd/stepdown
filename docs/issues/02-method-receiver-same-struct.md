## Parent PRD

#1

## What to build

Extend the analyzer to check method receiver functions on the same struct. When `(s *S) foo()` calls `(s *S) bar()`, and `bar` is declared before `foo`, report a violation.

## Implementation

All changes go in `analyzer.go` only — extend the existing `run()` logic. No new source files.

Changes to the existing algorithm:
- When collecting declarations, include methods: use `funcDecl.Recv` to detect receivers, store as `"TypeName.MethodName"` in the position map
- When walking function bodies, resolve `*ast.SelectorExpr` calls (e.g., `s.bar()`) where `s` is the receiver parameter — map back to `"TypeName.bar"`
- Use `recvTypeIdent()` pattern from recvcheck to extract receiver type name (handles both `*T` and `T`)
- Same comparison logic: if callee line < caller line → report

Key: only match calls where the selector's object is the function's own receiver variable, not arbitrary expressions.

## Acceptance criteria

- [ ] Same-struct method caller-before-callee violations are detected
- [ ] Both pointer and value receivers are handled
- [ ] Plain function checks from issue 01 continue to work
- [ ] No false positives when methods are correctly ordered

## TDD cycles

Each cycle is a commit: RED (failing test) → GREEN (minimal code to pass) → REFACTOR (clean up).

### Cycle 1: Pointer receiver method violation

**RED** — Create `testdata/src/methods/methods.go`:
```go
package methods

type Server struct{}

func (s *Server) handle() {} // want `function "handle" is called by "run" but declared before it \(stepdown rule\)`

func (s *Server) run() {
	s.handle()
}
```
Add test case in `analyzer_test.go`: `analysistest.Run(t, ..., a, "methods")`. Test fails.

**GREEN** — In `analyzer.go`:
1. When collecting declarations, detect receiver via `funcDecl.Recv`, store as `"Server.handle"` keyed by position
2. When walking bodies, resolve `s.handle()` → check if `s` is the receiver param → look up `"Server.handle"` in the map
3. Compare positions, report if violated

**REFACTOR** — Extract receiver type resolution into a helper if needed.

### Cycle 2: Value receiver

**RED** — Create `testdata/src/methods/value_receiver.go`:
```go
package methods

type Config struct{}

func (c Config) validate() {} // want `function "validate" is called by "load" but declared before it \(stepdown rule\)`

func (c Config) load() {
	c.validate()
}
```

**GREEN** — Ensure receiver type extraction handles both `*T` and `T`. Should work if using the `recvTypeIdent()` pattern.

### Cycle 3: Valid method ordering (no false positives)

**RED** — Create `testdata/src/methods/valid.go`:
```go
package methods

type Client struct{}

func (c *Client) send() {
	c.connect()
}

func (c *Client) connect() {}
```
Verify no diagnostics emitted.

## Blocked by

- Blocked by issue 01 (basic caller-before-callee)

## User stories addressed

- User story 3: method receiver functions checked (same struct)
