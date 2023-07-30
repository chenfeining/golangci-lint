package processors

import (
	"github.com/chenfeining/golangci-lint/pkg/result"
)

type Processor interface {
	Process(issues []result.Issue) ([]result.Issue, error)
	Name() string
	Finish()
}
