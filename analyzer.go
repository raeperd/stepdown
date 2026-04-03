// Package stepdown provides a Go linter that checks callers are declared before callees.
package stepdown

import (
	"go/ast"
	"go/token"
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

		// Determine receiver info for method declarations
		var recvVar, recvType string
		if funcDecl.Recv != nil {
			recvType = recvTypeName(funcDecl)
			if len(funcDecl.Recv.List) > 0 && len(funcDecl.Recv.List[0].Names) > 0 {
				recvVar = funcDecl.Recv.List[0].Names[0].Name
			}
		}

		seen := map[string]bool{}
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
				if ident, ok := fun.X.(*ast.Ident); ok && recvVar != "" && ident.Name == recvVar {
					calleeKey = recvType + "." + fun.Sel.Name
				}
			}
			if calleeKey == "" {
				return true
			}

			callee, exists := fileFuncs[calleeKey]
			if !exists || seen[calleeKey] {
				return true
			}
			if callee.line < callerPos.Line {
				seen[calleeKey] = true
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
	})

	return nil, nil
}

// recvTypeName extracts the receiver type name from a method declaration,
// handling both value receivers (T) and pointer receivers (*T).
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
