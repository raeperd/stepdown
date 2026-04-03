// Package stepdown provides a Go linter that checks callers are declared before callees.
package stepdown

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// NewAnalyzer returns a new analyzer to check for the stepdown rule.
func NewAnalyzer(s Settings) *analysis.Analyzer {
	exclusions := make(map[string]struct{}, len(s.Exclusions))
	for _, name := range s.Exclusions {
		exclusions[name] = struct{}{}
	}
	a := &analyzer{exclusions: exclusions}

	return &analysis.Analyzer{
		Name:     "stepdown",
		Doc:      "checks that callers are declared before callees (the stepdown rule)",
		Run:      a.run,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
}

// Settings is the configuration for the analyzer.
type Settings struct {
	// Exclusions is a list of function names to exclude from checks (e.g. "init", "main").
	Exclusions []string
}

type analyzer struct {
	exclusions map[string]struct{}
}

func (a *analyzer) run(pass *analysis.Pass) (any, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Collect all function declarations grouped by file
	type funcInfo struct {
		pos  token.Pos
		line int
	}
	funcs := map[string]map[string]funcInfo{} // filename -> funcName -> info

	insp.Preorder([]ast.Node{(*ast.FuncDecl)(nil)}, func(n ast.Node) {
		funcDecl := n.(*ast.FuncDecl)
		pos := pass.Fset.Position(funcDecl.Pos())
		if funcs[pos.Filename] == nil {
			funcs[pos.Filename] = map[string]funcInfo{}
		}
		key := funcDecl.Name.Name
		if funcDecl.Recv != nil {
			if typeName := recvTypeName(funcDecl); typeName != "" {
				key = typeName + "." + funcDecl.Name.Name
			}
		}
		funcs[pos.Filename][key] = funcInfo{pos: funcDecl.Pos(), line: pos.Line}
	})

	// Build call graph per file for cycle detection (skip closures to avoid false edges)
	callGraph := map[string]map[string]map[string]struct{}{} // filename -> caller -> callees
	insp.Preorder([]ast.Node{(*ast.FuncDecl)(nil)}, func(n ast.Node) {
		funcDecl := n.(*ast.FuncDecl)
		if funcDecl.Body == nil {
			return
		}
		pos := pass.Fset.Position(funcDecl.Pos())
		callerKey := funcDecl.Name.Name
		if funcDecl.Recv != nil {
			if typeName := recvTypeName(funcDecl); typeName != "" {
				callerKey = typeName + "." + funcDecl.Name.Name
			}
		}
		fileFuncs := funcs[pos.Filename]
		if fileFuncs == nil {
			return
		}
		callees := map[string]struct{}{}
		ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
			// Don't descend into closures — their calls belong to the closure, not the enclosing func
			if _, ok := n.(*ast.FuncLit); ok {
				return false
			}
			callExpr, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			switch fun := callExpr.Fun.(type) {
			case *ast.Ident:
				if _, exists := fileFuncs[fun.Name]; exists {
					callees[fun.Name] = struct{}{}
				}
			case *ast.SelectorExpr:
				if sel, ok := pass.TypesInfo.Selections[fun]; ok {
					if fn, ok := sel.Obj().(*types.Func); ok {
						fnPos := pass.Fset.Position(fn.Pos())
						if fnPos.Filename == pos.Filename {
							recv := sel.Recv()
							if ptr, ok := recv.(*types.Pointer); ok {
								recv = ptr.Elem()
							}
							if named, ok := recv.(*types.Named); ok {
								key := named.Obj().Name() + "." + fun.Sel.Name
								if _, exists := fileFuncs[key]; exists {
									callees[key] = struct{}{}
								}
							}
						}
					}
				}
			}
			return true
		})
		if len(callees) > 0 {
			if callGraph[pos.Filename] == nil {
				callGraph[pos.Filename] = map[string]map[string]struct{}{}
			}
			callGraph[pos.Filename][callerKey] = callees
		}
	})

	// Check each function's calls
	insp.Preorder([]ast.Node{(*ast.FuncDecl)(nil)}, func(n ast.Node) {
		funcDecl := n.(*ast.FuncDecl)
		if funcDecl.Body == nil {
			return
		}
		callerPos := pass.Fset.Position(funcDecl.Pos())
		fileFuncs := funcs[callerPos.Filename]
		if fileFuncs == nil {
			return
		}
		callerKey := funcDecl.Name.Name
		if funcDecl.Recv != nil {
			if typeName := recvTypeName(funcDecl); typeName != "" {
				callerKey = typeName + "." + funcDecl.Name.Name
			}
		}

		seen := map[string]bool{}
		var invocationOrder []string
		ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
			callExpr, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			var calleeKey string
			switch fun := callExpr.Fun.(type) {
			case *ast.Ident:
				calleeKey = fun.Name
			case *ast.SelectorExpr:
				// Use type info to resolve method calls (handles cross-struct)
				if sel, ok := pass.TypesInfo.Selections[fun]; ok {
					if fn, ok := sel.Obj().(*types.Func); ok {
						fnPos := pass.Fset.Position(fn.Pos())
						if fnPos.Filename == callerPos.Filename {
							recv := sel.Recv()
							if ptr, ok := recv.(*types.Pointer); ok {
								recv = ptr.Elem()
							}
							if named, ok := recv.(*types.Named); ok {
								calleeKey = named.Obj().Name() + "." + fun.Sel.Name
							}
						}
					}
				}
			}
			if calleeKey == "" {
				return true
			}

			callee, exists := fileFuncs[calleeKey]
			if !exists || seen[calleeKey] {
				return true
			}
			seen[calleeKey] = true
			invocationOrder = append(invocationOrder, calleeKey)
			if callee.line < callerPos.Line {
				// Skip circular calls — if callee can reach back to caller, neither ordering works
				if fileGraph := callGraph[callerPos.Filename]; fileGraph != nil {
					if reachable(fileGraph, calleeKey, callerKey) {
						return true
					}
				}
				// Use short name (without type prefix) for the diagnostic message
				_, calleeName, _ := strings.Cut(calleeKey, ".")
				if calleeName == "" {
					calleeName = calleeKey
				}
				callerName := funcDecl.Name.Name
				// Skip excluded functions (as caller or callee)
				if _, ok := a.exclusions[callerName]; ok {
					return true
				}
				if _, ok := a.exclusions[calleeName]; ok {
					return true
				}
				pass.Reportf(callee.pos,
					"function %q is called by %q but declared before it (stepdown rule)",
					calleeName, callerName,
				)
			}
			return true
		})

		// Check callee invocation order: callees should be declared in the order they are invoked
		maxLine := 0
		var maxKey string
		for _, calleeKey := range invocationOrder {
			callee := fileFuncs[calleeKey]
			if callee.line < maxLine {
				_, calleeName, _ := strings.Cut(calleeKey, ".")
				if calleeName == "" {
					calleeName = calleeKey
				}
				_, maxName, _ := strings.Cut(maxKey, ".")
				if maxName == "" {
					maxName = maxKey
				}
				pass.Reportf(callee.pos,
					"function %q is called by %q before %q but declared after it (stepdown rule)",
					calleeName, funcDecl.Name.Name, maxName,
				)
			}
			if callee.line > maxLine {
				maxLine = callee.line
				maxKey = calleeKey
			}
		}
	})

	return nil, nil
}

func recvTypeName(funcDecl *ast.FuncDecl) string {
	if funcDecl.Recv == nil || len(funcDecl.Recv.List) == 0 {
		return ""
	}
	t := funcDecl.Recv.List[0].Type
	if star, ok := t.(*ast.StarExpr); ok {
		t = star.X
	}
	if ident, ok := t.(*ast.Ident); ok {
		return ident.Name
	}
	return ""
}

func reachable(graph map[string]map[string]struct{}, src, dst string) bool {
	visited := map[string]bool{}
	stack := []string{src}
	for len(stack) > 0 {
		node := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if node == dst {
			return true
		}
		if visited[node] {
			continue
		}
		visited[node] = true
		for next := range graph[node] {
			stack = append(stack, next)
		}
	}
	return false
}
