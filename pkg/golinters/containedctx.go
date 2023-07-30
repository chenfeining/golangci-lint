package golinters

import (
	"github.com/sivchari/containedctx"
	"golang.org/x/tools/go/analysis"

	"github.com/chenfeining/golangci-lint/pkg/golinters/goanalysis"
)

func NewContainedCtx() *goanalysis.Linter {
	a := containedctx.Analyzer

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
