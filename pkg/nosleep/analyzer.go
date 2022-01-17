package nosleep

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "nosleep",
	Doc:      "Checks for usages of time.Sleep.",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	nodeFilter := []ast.Node{ // filter needed nodes: visit only them
		(*ast.CallExpr)(nil),
	}

	pass.ResultOf[inspect.Analyzer].(*inspector.Inspector).Preorder(nodeFilter, func(node ast.Node) {
		call := node.(*ast.CallExpr) // node is always a CallExpr thanks to nodeFilter

		if expr, ok := call.Fun.(*ast.SelectorExpr); ok {
			if ident, ok := expr.X.(*ast.Ident); ok {
				if ident.Name == "time" && expr.Sel.Name == "Sleep" {
					pass.Reportf(node.Pos(), "time.Sleep detected")
				}
			}
		}

		return
	})

	return nil, nil
}
