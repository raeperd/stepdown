// Package stepdown provides a Go linter that checks callers are declared before callees.
package stepdown

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// NewAnalyzer returns a new analyzer to check for the stepdown rule.
func NewAnalyzer(s Settings) *analysis.Analyzer {
	exclusions := make(map[string]struct{}, len(s.Exclusions))
	for _, name := range s.Exclusions {
		exclusions[name] = struct{}{}
	}
	a := &analyzer{exclusions: exclusions}

	return &analysis.Analyzer{
		Name: "stepdown",
		Doc:  "checks that callers are declared before callees (the stepdown rule)",
		Run:  a.run,
	}
}

// Settings is the configuration for the analyzer.
type Settings struct {
	// Exclusions is a list of names to exclude from checks.
	// Supported forms:
	//   - plain function names (e.g. "init", "main")
	//   - exact qualified method names (e.g. "Server.handle")
	//   - short method names that match across receiver types (e.g. "handle")
	Exclusions []string
}

type analyzer struct {
	exclusions map[string]struct{}
}

func (a *analyzer) run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		a.checkFile(pass, file)
	}
	return nil, nil
}

func (a *analyzer) checkFile(pass *analysis.Pass, file *ast.File) {
	filename := pass.Fset.Position(file.Pos()).Filename

	// Collect function declarations and build call graph (deduplicated, in invocation order)
	funcs := map[string]token.Pos{} // funcKey -> declaration pos
	calls := map[string][]string{}  // caller -> unique callees in invocation order
	for _, decl := range file.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		key := funcKey(funcDecl)
		funcs[key] = funcDecl.Pos()
		if funcDecl.Body == nil {
			continue
		}
		seen := map[string]struct{}{}
		var callees []string
		ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
			switch n := n.(type) {
			case *ast.FuncLit:
				return false
			case *ast.CallExpr:
				// Find the called function or method if it's declared in this file.
				var fn *types.Func
				switch fun := n.Fun.(type) {
				case *ast.Ident: // plain call: foo()
					if f, ok := pass.TypesInfo.Uses[fun].(*types.Func); ok {
						fn = f
					}
				case *ast.SelectorExpr: // method call: s.foo()
					if sel, ok := pass.TypesInfo.Selections[fun]; ok {
						if f, ok := sel.Obj().(*types.Func); ok {
							fn = f
						}
					}
				}
				if fn == nil || pass.Fset.Position(fn.Pos()).Filename != filename {
					return true
				}
				calleeKey := funcName(fn)
				if _, ok := seen[calleeKey]; !ok {
					seen[calleeKey] = struct{}{}
					callees = append(callees, calleeKey)
				}
			}
			return true
		})
		if len(callees) > 0 {
			calls[key] = callees
		}
	}

	// Report caller-before-callee violations and callee invocation order violations
	for _, decl := range file.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok || funcDecl.Body == nil {
			continue
		}
		callerKey := funcKey(funcDecl)
		if a.isExcluded(callerKey) {
			continue
		}
		callerLine := pass.Fset.Position(funcDecl.Pos()).Line

		maxLine := 0
		var maxKey string
		for _, calleeKey := range calls[callerKey] {
			// Skip if caller and callee form a cycle (e.g. a→b→a).
			// In a cycle, at least one edge must go backward — moving the callee
			// after the caller would just create a new violation elsewhere in the cycle.
			if inCycle(calls, calleeKey, callerKey) {
				continue
			}
			if a.isExcluded(calleeKey) {
				continue
			}

			calleePos := funcs[calleeKey]
			calleeLine := pass.Fset.Position(calleePos).Line

			// Violation: callee declared before caller
			if calleeLine < callerLine {
				pass.Reportf(calleePos,
					"function %q is called by %q but declared before it (stepdown rule)",
					calleeKey, callerKey,
				)
			}

			// Violation: callees declared in different order than invoked
			if calleeLine < maxLine {
				pass.Reportf(calleePos,
					"function %q is called by %q before %q but declared after it (stepdown rule)",
					calleeKey, callerKey, maxKey,
				)
			}
			if calleeLine > maxLine {
				maxLine = calleeLine
				maxKey = calleeKey
			}
		}
	}
}

func funcName(fn *types.Func) string {
	sig, ok := fn.Type().(*types.Signature)
	if !ok || sig.Recv() == nil {
		return fn.Name()
	}
	t := sig.Recv().Type()
	if ptr, ok := t.(*types.Pointer); ok {
		t = ptr.Elem()
	}
	if named, ok := t.(*types.Named); ok {
		return named.Obj().Name() + "." + fn.Name()
	}
	return fn.Name()
}

func funcKey(funcDecl *ast.FuncDecl) string {
	if funcDecl.Recv != nil && len(funcDecl.Recv.List) > 0 {
		t := funcDecl.Recv.List[0].Type
		if star, ok := t.(*ast.StarExpr); ok {
			t = star.X
		}
		if ident, ok := t.(*ast.Ident); ok {
			return ident.Name + "." + funcDecl.Name.Name
		}
	}
	return funcDecl.Name.Name
}

func (a *analyzer) isExcluded(key string) bool {
	if _, ok := a.exclusions[key]; ok {
		return true
	}
	_, ok := a.exclusions[shortName(key)]
	return ok
}

func shortName(key string) string {
	_, name, _ := strings.Cut(key, ".")
	if name == "" {
		return key
	}
	return name
}

func inCycle(graph map[string][]string, src, dst string) bool {
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
		stack = append(stack, graph[node]...)
	}
	return false
}
