package golinters

import (
	"github.com/kunwardeep/paralleltest/pkg/paralleltest"
	"golang.org/x/tools/go/analysis"

	"github.com/chenfeining/golangci-lint/pkg/config"
	"github.com/chenfeining/golangci-lint/pkg/golinters/goanalysis"
)

func NewParallelTest(settings *config.ParallelTestSettings) *goanalysis.Linter {
	a := paralleltest.Analyzer

	var cfg map[string]map[string]any
	if settings != nil {
		cfg = map[string]map[string]any{
			a.Name: {
				"i": settings.IgnoreMissing,
			},
		}
	}

	return goanalysis.NewLinter(
		"paralleltest",
		"paralleltest detects missing usage of t.Parallel() method in your Go test",
		[]*analysis.Analyzer{a},
		cfg,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
