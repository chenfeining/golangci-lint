package processors

import (
	"regexp"

	"github.com/chenfeining/golangci-lint/pkg/fsutils"
	"github.com/chenfeining/golangci-lint/pkg/logutils"
	"github.com/chenfeining/golangci-lint/pkg/result"
)

type BaseRule struct {
	Text       string
	Source     string
	Path       string
	PathExcept string
	Linters    []string
}

type baseRule struct {
	text       *regexp.Regexp
	source     *regexp.Regexp
	path       *regexp.Regexp
	pathExcept *regexp.Regexp
	linters    []string
}

func (r *baseRule) isEmpty() bool {
	return r.text == nil && r.source == nil && r.path == nil && r.pathExcept == nil && len(r.linters) == 0
}

func (r *baseRule) match(issue *result.Issue, files *fsutils.Files, log logutils.Log) bool {
	if r.isEmpty() {
		return false
	}
	if r.text != nil && !r.text.MatchString(issue.Text) {
		return false
	}
	if r.path != nil && !r.path.MatchString(files.WithPathPrefix(issue.FilePath())) {
		return false
	}
	if r.pathExcept != nil && r.pathExcept.MatchString(issue.FilePath()) {
		return false
	}
	if len(r.linters) != 0 && !r.matchLinter(issue) {
		return false
	}

	// the most heavyweight checking last
	if r.source != nil && !r.matchSource(issue, files.LineCache, log) {
		return false
	}

	return true
}

func (r *baseRule) matchLinter(issue *result.Issue) bool {
	for _, linter := range r.linters {
		if linter == issue.FromLinter {
			return true
		}
	}

	return false
}

func (r *baseRule) matchSource(issue *result.Issue, lineCache *fsutils.LineCache, log logutils.Log) bool { //nolint:interfacer
	sourceLine, errSourceLine := lineCache.GetLine(issue.FilePath(), issue.Line())
	if errSourceLine != nil {
		log.Warnf("Failed to get line %s:%d from line cache: %s", issue.FilePath(), issue.Line(), errSourceLine)
		return false // can't properly match
	}

	return r.source.MatchString(sourceLine)
}
