## Parent PRD

#1

## What to build

Merge the first two `Preorder` passes (collect declarations + build call graph) into a single pass. Both iterate over all `*ast.FuncDecl` nodes and could be combined to reduce overhead.

## Implementation

All changes go in `analyzer.go` only. Combine the declaration collection (pass 1, line 51-63) and call graph building (pass 2, line 68-107) into one `Preorder` callback.

## Acceptance criteria

- [ ] Only two `Preorder` passes remain (collect + graph, then check)
- [ ] All existing tests pass unchanged
- [ ] Self-lint passes

## TDD cycles

### Cycle 1: Merge passes — behavior-preserving refactor

**RED** — No new test needed. Run existing tests to establish baseline.

**GREEN** — Merge the two `Preorder` callbacks into one:
```go
insp.Preorder([]ast.Node{(*ast.FuncDecl)(nil)}, func(n ast.Node) {
    funcDecl := n.(*ast.FuncDecl)
    // ... collect declaration (existing pass 1 logic) ...
    // ... build call graph (existing pass 2 logic) ...
})
```

**REFACTOR** — Remove the now-empty second pass. Clean up variable scoping.

Note: the call graph pass currently depends on `funcs[pos.Filename]` being populated. Since both now run in the same callback and declarations are processed in file order, the current function's file may not have all declarations yet when building the call graph. Solution: build the call graph using a temporary callee-name set, then filter against `fileFuncs` when checking violations instead of when building the graph.
