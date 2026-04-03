## Parent PRD

#1

## What to build

Extend the analyzer to check method calls across different struct types within the same file. When a method on struct A calls a function or method on struct B (declared in the same file), the callee should be declared after the caller.

## Implementation

All changes go in `analyzer.go` only — extend existing method resolution. No new source files.

Changes to the existing algorithm:
- When resolving `*ast.SelectorExpr` calls, don't restrict to receiver-only matches
- Use type information from `pass.TypesInfo.ObjectOf()` or `pass.TypesInfo.Uses` to resolve the callee's declaration
- If the resolved declaration is a `*types.Func` in the same file, compare positions
- This generalizes the method resolution from issue 02 — receiver calls become a special case

Note: this may simplify the issue 02 implementation by using type info instead of AST-only matching.

## Acceptance criteria

- [ ] Cross-struct method call ordering violations are detected
- [ ] Both plain function → method and method → method cross-struct calls are checked
- [ ] No false positives when cross-struct methods are correctly ordered
- [ ] Same-struct and plain function checks continue to work

## TDD cycles

Each cycle is a commit: RED (failing test) → GREEN (minimal code to pass) → REFACTOR (clean up).

### Cycle 1: Cross-struct method call violation

**RED** — Create `testdata/src/crossstruct/crossstruct.go`:
```go
package crossstruct

type DB struct{}

func (d *DB) query() {} // want `function "query" is called by "handle" but declared before it \(stepdown rule\)`

type Handler struct {
	db *DB
}

func (h *Handler) handle() {
	h.db.query()
}
```
Add test case in `analyzer_test.go`. Test fails.

**GREEN** — Extend `*ast.SelectorExpr` resolution to use `pass.TypesInfo` for resolving the callee, regardless of whether the receiver is the current function's own receiver.

**REFACTOR** — If issue 02's AST-only matching can be replaced by the type-info approach, simplify.

### Cycle 2: Plain function calling a method

**RED** — Create `testdata/src/crossstruct/func_to_method.go`:
```go
package crossstruct

type Logger struct{}

func (l *Logger) log() {} // want `function "log" is called by "process" but declared before it \(stepdown rule\)`

func process(l *Logger) {
	l.log()
}
```

**GREEN** — Should work with the type-info approach from cycle 1.

### Cycle 3: Valid cross-struct ordering

**RED** — Create `testdata/src/crossstruct/valid.go`:
```go
package crossstruct

type Service struct {
	repo *Repo
}

func (s *Service) serve() {
	s.repo.find()
}

type Repo struct{}

func (r *Repo) find() {}
```
Verify no diagnostics.

## Blocked by

- Blocked by issue 02 (method receiver same-struct)

## User stories addressed

- User story 9: cross-struct method calls checked within the same file
