package runner

import (
	"regexp"

	messages "github.com/cucumber/cucumber-messages-go/v3"
	tagexpressions "github.com/cucumber/tag-expressions-go"
)

// PickleFilter filters pickles
type PickleFilter struct {
	nameRegexps   []*regexp.Regexp
	lines         map[string][]uint64
	tagExpression tagexpressions.Evaluatable
}

// NewPickleFilter returns a PickleFilter
func NewPickleFilter(config *messages.SourcesFilterConfig) (*PickleFilter, error) {
	tagExpression, err := tagexpressions.Parse(config.TagExpression)
	if err != nil {
		return nil, err
	}
	nameRegexps := make([]*regexp.Regexp, len(config.GetNameRegularExpressions()))
	for i, name := range config.GetNameRegularExpressions() {
		nameRegexps[i] = regexp.MustCompilePOSIX(name)
	}
	lines := map[string][]uint64{}
	for _, uriToLines := range config.GetUriToLinesMapping() {
		lines[uriToLines.GetAbsolutePath()] = uriToLines.GetLines()
	}
	return &PickleFilter{
		nameRegexps:   nameRegexps,
		lines:         lines,
		tagExpression: tagExpression,
	}, nil
}

// Matches returns whether the pickle matches the filters
func (p *PickleFilter) Matches(pickle *messages.Pickle) bool {
	return p.matchesAnyLine(pickle) &&
		p.matchesAnyName(pickle) &&
		p.matchesTagExpression(pickle)
}

func (p *PickleFilter) matchesAnyLine(pickle *messages.Pickle) bool {
	uriLines, ok := p.lines[pickle.Uri]
	if !ok || len(uriLines) == 0 {
		return true
	}
	for _, line := range uriLines {
		for _, location := range pickle.Locations {
			if line == uint64(location.Line) {
				return true
			}
		}
	}
	return false
}

func (p *PickleFilter) matchesAnyName(pickle *messages.Pickle) bool {
	if len(p.nameRegexps) == 0 {
		return true
	}
	for _, nameRegexp := range p.nameRegexps {
		if nameRegexp.MatchString(pickle.Name) {
			return true
		}
	}
	return false
}

func (p *PickleFilter) matchesTagExpression(pickle *messages.Pickle) bool {
	tagNames := make([]string, len(pickle.Tags))
	for i, tag := range pickle.Tags {
		tagNames[i] = tag.Name
	}
	return p.tagExpression.Evaluate(tagNames)
}
