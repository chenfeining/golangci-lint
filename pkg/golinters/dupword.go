package golinters

import (
	"strings"

	"github.com/Abirdcfly/dupword"
	"golang.org/x/tools/go/analysis"

	"github.com/chenfeining/golangci-lint/pkg/config"
	"github.com/chenfeining/golangci-lint/pkg/golinters/goanalysis"
)

func NewDupWord(setting *config.DupWordSettings) *goanalysis.Linter {
	a := dupword.NewAnalyzer()

	cfgMap := map[string]map[string]any{}
	if setting != nil {
		cfgMap[a.Name] = map[string]any{
			"keyword": strings.Join(setting.Keywords, ","),
		}
	}

	return goanalysis.NewLinter(
		a.Name,
		"checks for duplicate words in the source code",
		[]*analysis.Analyzer{a},
		cfgMap,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
