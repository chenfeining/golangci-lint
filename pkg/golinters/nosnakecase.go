package golinters

import (
	"github.com/sivchari/nosnakecase"
	"golang.org/x/tools/go/analysis"

	"github.com/chenfeining/golangci-lint/pkg/golinters/goanalysis"
)

func NewNoSnakeCase() *goanalysis.Linter {
	a := nosnakecase.Analyzer

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
