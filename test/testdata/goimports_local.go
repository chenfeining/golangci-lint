//golangcitest:args -Egoimports
//golangcitest:config_path testdata/configs/goimports_local.yml
package testdata

import (
	"fmt"

	"github.com/chenfeining/golangci-lint/pkg/config" // want "File is not `goimports`-ed with -local github.com/chenfeining/golangci-lint"
	"golang.org/x/tools/go/analysis"
)

func GoimportsLocalPrefixTest() {
	fmt.Print("x")
	_ = config.Config{}
	_ = analysis.Analyzer{}
}
