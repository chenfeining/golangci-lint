package golinters

import (
	"github.com/lufeee/execinquery"
	"golang.org/x/tools/go/analysis"

	"github.com/chenfeining/golangci-lint/pkg/golinters/goanalysis"
)

func NewExecInQuery() *goanalysis.Linter {
	a := execinquery.Analyzer

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
