package golinters

import (
	"github.com/nishanths/predeclared/passes/predeclared"
	"golang.org/x/tools/go/analysis"

	"github.com/chenfeining/golangci-lint/pkg/config"
	"github.com/chenfeining/golangci-lint/pkg/golinters/goanalysis"
)

func NewPredeclared(settings *config.PredeclaredSettings) *goanalysis.Linter {
	a := predeclared.Analyzer

	var cfg map[string]map[string]any
	if settings != nil {
		cfg = map[string]map[string]any{
			a.Name: {
				predeclared.IgnoreFlag:    settings.Ignore,
				predeclared.QualifiedFlag: settings.Qualified,
			},
		}
	}

	return goanalysis.NewLinter(a.Name, a.Doc, []*analysis.Analyzer{a}, cfg).
		WithLoadMode(goanalysis.LoadModeSyntax)
}
