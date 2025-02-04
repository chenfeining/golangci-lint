package golinters

import (
	"github.com/sivchari/tenv"
	"golang.org/x/tools/go/analysis"

	"github.com/chenfeining/golangci-lint/pkg/config"
	"github.com/chenfeining/golangci-lint/pkg/golinters/goanalysis"
)

func NewTenv(settings *config.TenvSettings) *goanalysis.Linter {
	a := tenv.Analyzer

	var cfg map[string]map[string]any
	if settings != nil {
		cfg = map[string]map[string]any{
			a.Name: {
				tenv.A: settings.All,
			},
		}
	}

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		cfg,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
