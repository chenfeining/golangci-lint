//golangcitest:args -Egci
//golangcitest:config_path testdata/configs/gci.yml
package testdata

import (
	"errors"
	"fmt"
	"github.com/chenfeining/golangci-lint/pkg/config"
	gcicfg "github.com/daixiang0/gci/pkg/config" // want "File is not \\`gci\\`-ed with --skip-generated -s standard -s prefix\\(github.com/chenfeining/golangci-lint,github.com/daixiang0/gci\\) -s default --custom-order"
	"golang.org/x/tools/go/analysis"             // want "File is not \\`gci\\`-ed with --skip-generated -s standard -s prefix\\(github.com/chenfeining/golangci-lint,github.com/daixiang0/gci\\) -s default --custom-order"
)

func GoimportsLocalTest() {
	fmt.Print(errors.New("x"))
	_ = config.Config{}
	_ = analysis.Analyzer{}
	_ = gcicfg.BoolConfig{}
}
