package golinters

import (
	"github.com/Djarvur/go-err113"
	"golang.org/x/tools/go/analysis"

	"github.com/chenfeining/golangci-lint/pkg/golinters/goanalysis"
)

func NewGoerr113() *goanalysis.Linter {
	return goanalysis.NewLinter(
		"goerr113",
		"Go linter to check the errors handling expressions",
		[]*analysis.Analyzer{err113.NewAnalyzer()},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
