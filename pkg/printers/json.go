package printers

import (
	"encoding/json"
	"io"

	"github.com/chenfeining/golangci-lint/pkg/report"
	"github.com/chenfeining/golangci-lint/pkg/result"
)

type JSON struct {
	rd *report.Data
	w  io.Writer
}

func NewJSON(rd *report.Data, w io.Writer) *JSON {
	return &JSON{
		rd: rd,
		w:  w,
	}
}

type JSONResult struct {
	Issues []result.Issue
	Report *report.Data
}

func (p JSON) Print(issues []result.Issue) error {
	res := JSONResult{
		Issues: issues,
		Report: p.rd,
	}
	if res.Issues == nil {
		res.Issues = []result.Issue{}
	}

	return json.NewEncoder(p.w).Encode(res)
}
