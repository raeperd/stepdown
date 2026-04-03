## Parent PRD

#1

## What to build

Fix the violation-checking pass to skip `*ast.FuncLit` nodes, preventing calls inside closures from being attributed to the enclosing named function.

Currently the call graph building pass correctly skips closures (line 86-88), but the violation-checking pass's `ast.Inspect` (line 123) descends into closures. This means a closure calling `setup()` inside `register()` is reported as `register` calling `setup`.

## Implementation

All changes go in `analyzer.go` only. Add the same `*ast.FuncLit` guard from the call graph pass into the violation-checking `ast.Inspect` block.

## Acceptance criteria

- [ ] Calls inside closures do not produce violations for the enclosing function
- [ ] Direct calls in the same function body still produce violations
- [ ] Existing tests continue to pass

## TDD cycles

### Cycle 1: Closure call attributed to enclosing function — false positive

**RED** — Create `testdata/src/basic/closure.go`:
```go
package basic

func helper() {}

func register() {
	callback := func() {
		helper()
	}
	_ = callback
}
```
No `// want` — `register` does not directly call `helper`, the closure does. Test fails if the linter reports a violation on `helper`.

Wait — `helper` is declared before `register`, so `helper.line < register.line`. The linter currently descends into the closure, sees `helper()`, and reports it. But `helper` has no `// want` comment, so `analysistest` will flag the unexpected diagnostic. This is the RED.

**GREEN** — In the violation-checking `ast.Inspect` (line 123), add:
```go
if _, ok := n.(*ast.FuncLit); ok {
    return false
}
```

**REFACTOR** — Clean up if needed.
