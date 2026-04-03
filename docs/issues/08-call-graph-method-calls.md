## Parent PRD

#1

## What to build

Fix the call graph to include method calls (`*ast.SelectorExpr`) so that circular method call chains are correctly detected and suppressed.

Currently the call graph (used for cycle detection) only captures plain function calls via `*ast.Ident`. Method calls like `s.stop()` are invisible to the graph, causing false positives when methods form cycles.

## Implementation

All changes go in `analyzer.go` only. In the call graph building pass (second `Preorder`), extend the `ast.Inspect` to also resolve `*ast.SelectorExpr` calls using `pass.TypesInfo.Selections`, matching the pattern already used in the violation-checking pass.

This requires passing `pass` into the call graph building logic since it currently only uses the inspector.

## Acceptance criteria

- [ ] Circular method calls (same-struct and cross-struct) produce no violations
- [ ] Non-circular method calls still produce violations
- [ ] Existing plain function circular detection still works

## TDD cycles

### Cycle 1: Circular same-struct methods — false positive

**RED** — Create `testdata/src/circular/methods.go`:
```go
package circular

type Worker struct{}

func (w *Worker) start() {
	w.stop()
}

func (w *Worker) stop() {
	w.start()
}
```
Add to existing `circular` test. Test fails — reports `start` as violation because call graph misses `w.stop()` and `w.start()`.

**GREEN** — In the call graph building pass, add `*ast.SelectorExpr` resolution via `pass.TypesInfo.Selections`. Store method callees as `"TypeName.MethodName"` keys matching the existing declaration keys.

**REFACTOR** — Clean up if needed.
