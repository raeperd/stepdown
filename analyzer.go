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
		if funcDecl.Recv != nil {
			return // skip methods for now
		}
		pos := pass.Fset.Position(funcDecl.Pos())
		if funcs[pos.Filename] == nil {
			funcs[pos.Filename] = map[string]funcInfo{}
		}
		funcs[pos.Filename][funcDecl.Name.Name] = funcInfo{pos: funcDecl.Pos(), line: pos.Line}
	})

	// Check each function's calls
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
				if _, excluded := a.exclusions[funcDecl.Name.Name]; excluded {
					return true
				}
				if _, excluded := a.exclusions[ident.Name]; excluded {
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
