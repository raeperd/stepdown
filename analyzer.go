// Package stepdown provides a Go linter that checks callers are declared before callees.
package stepdown

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// NewAnalyzer returns a new analyzer to check for the stepdown rule.
func NewAnalyzer(s Settings) *analysis.Analyzer {
	a := &analyzer{}

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

type analyzer struct{}

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
		if funcDecl.Recv != nil {
			return // skip methods for now
		}
		pos := pass.Fset.Position(funcDecl.Pos())
		if funcs[pos.Filename] == nil {
			funcs[pos.Filename] = map[string]funcInfo{}
		}
		funcs[pos.Filename][funcDecl.Name.Name] = funcInfo{pos: funcDecl.Pos(), line: pos.Line}
	})

	// Build call graph per file: caller -> set of callees
	callGraph := map[string]map[string]map[string]bool{} // filename -> caller -> callees

	insp.Preorder([]ast.Node{(*ast.FuncDecl)(nil)}, func(n ast.Node) {
		funcDecl := n.(*ast.FuncDecl)
		if funcDecl.Recv != nil || funcDecl.Body == nil {
			return
		}
		callerName := funcDecl.Name.Name
		callerPos := pass.Fset.Position(funcDecl.Pos())
		fileFuncs := funcs[callerPos.Filename]
		if fileFuncs == nil {
			return
		}

		ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
			callExpr, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			ident, ok := callExpr.Fun.(*ast.Ident)
			if !ok {
				return true
			}
			if _, exists := fileFuncs[ident.Name]; !exists {
				return true
			}
			if callGraph[callerPos.Filename] == nil {
				callGraph[callerPos.Filename] = map[string]map[string]bool{}
			}
			if callGraph[callerPos.Filename][callerName] == nil {
				callGraph[callerPos.Filename][callerName] = map[string]bool{}
			}
			callGraph[callerPos.Filename][callerName][ident.Name] = true
			return true
		})
	})

	// Check each function's calls, skipping circular pairs
	insp.Preorder([]ast.Node{(*ast.FuncDecl)(nil)}, func(n ast.Node) {
		funcDecl := n.(*ast.FuncDecl)
		if funcDecl.Recv != nil || funcDecl.Body == nil {
			return
		}
		callerPos := pass.Fset.Position(funcDecl.Pos())
		fileFuncs := funcs[callerPos.Filename]
		if fileFuncs == nil {
			return
		}
		fileGraph := callGraph[callerPos.Filename]

		seen := map[string]bool{}
		ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
			callExpr, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			ident, ok := callExpr.Fun.(*ast.Ident)
			if !ok {
				return true
			}
			callee, exists := fileFuncs[ident.Name]
			if !exists || seen[ident.Name] {
				return true
			}
			if callee.line < callerPos.Line {
				if reachable(fileGraph, ident.Name, funcDecl.Name.Name) {
					return true
				}
				seen[ident.Name] = true
				pass.Reportf(callee.pos,
					"function %q is called by %q but declared before it (stepdown rule)",
					ident.Name, funcDecl.Name.Name,
				)
			}
			return true
		})
	})

	return nil, nil
}

// reachable returns true if there is a path from src to dst in the call graph using DFS.
func reachable(graph map[string]map[string]bool, src, dst string) bool {
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
