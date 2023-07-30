//golangcitest:args -Egci
//golangcitest:config_path testdata/configs/gci.yml
//golangcitest:expected_exitcode 0
package gci

import (
	"github.com/chenfeining/golangci-lint/pkg/config"
	"golang.org/x/tools/go/analysis"

	"fmt"

	gcicfg "github.com/daixiang0/gci/pkg/config"
)

func GoimportsLocalTest() {
	fmt.Print("x")
	_ = config.Config{}
	_ = analysis.Analyzer{}
	_ = gcicfg.BoolConfig{}
}
