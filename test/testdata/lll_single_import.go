//golangcitest:args -Elll
//golangcitest:config_path testdata/configs/lll_import.yml
//golangcitest:expected_exitcode 0
package testdata

import veryLongImportAliasNameForTest "github.com/chenfeining/golangci-lint/internal/golinters"

func LllSingleImport() {
	_ = veryLongImportAliasNameForTest.NewLLL(nil)
}
