package golinters

import (
	"github.com/GaijinEntertainment/go-exhaustruct/v3/analyzer"
	"golang.org/x/tools/go/analysis"

	"github.com/chenfeining/golangci-lint/pkg/config"
	"github.com/chenfeining/golangci-lint/pkg/golinters/goanalysis"
)

func NewExhaustruct(settings *config.ExhaustructSettings) *goanalysis.Linter {
	var include, exclude []string
	if settings != nil {
		include = settings.Include
		exclude = settings.Exclude
	}

	a, err := analyzer.NewAnalyzer(include, exclude)
	if err != nil {
		linterLogger.Fatalf("exhaustruct configuration: %v", err)
	}

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
