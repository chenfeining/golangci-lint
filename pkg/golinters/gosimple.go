package golinters

import (
	"honnef.co/go/tools/simple"

	"github.com/chenfeining/golangci-lint/pkg/config"
	"github.com/chenfeining/golangci-lint/pkg/golinters/goanalysis"
)

func NewGosimple(settings *config.StaticCheckSettings) *goanalysis.Linter {
	cfg := staticCheckConfig(settings)

	analyzers := setupStaticCheckAnalyzers(simple.Analyzers, getGoVersion(settings), cfg.Checks)

	return goanalysis.NewLinter(
		"gosimple",
		"Linter for Go source code that specializes in simplifying code",
		analyzers,
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
