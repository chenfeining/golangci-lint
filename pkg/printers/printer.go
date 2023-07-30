package printers

import (
	"github.com/chenfeining/golangci-lint/pkg/result"
)

type Printer interface {
	Print(issues []result.Issue) error
}
