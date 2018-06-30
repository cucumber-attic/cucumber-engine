package runner

import (
	"regexp"

	"github.com/cucumber/cucumber-engine/src/dto"
	gherkin "github.com/cucumber/gherkin-go"
	tagexpressions "github.com/cucumber/tag-expressions-go"
)

// PickleFilter filters pickles
type PickleFilter struct {
	nameRegexps   []*regexp.Regexp
	lines         map[string][]int
	tagExpression tagexpressions.Evaluatable
}

// NewPickleFilter returns a PickleFilter
func NewPickleFilter(config *dto.FeaturesFilterConfig) (*PickleFilter, error) {
	tagExpression, err := tagexpressions.Parse(config.TagExpression)
	if err != nil {
		return nil, err
	}
	nameRegexps := make([]*regexp.Regexp, len(config.Names))
	for i, name := range config.Names {
		nameRegexps[i] = regexp.MustCompilePOSIX(name)
	}
	return &PickleFilter{
		nameRegexps:   nameRegexps,
		lines:         config.Lines,
		tagExpression: tagExpression,
	}, nil
}

// Matches returns whether the pickle matches the filters
func (p *PickleFilter) Matches(pickleEvent *gherkin.PickleEvent) bool {
	return p.matchesAnyLine(pickleEvent) &&
		p.matchesAnyName(pickleEvent) &&
		p.matchesTagExpression(pickleEvent)
}

func (p *PickleFilter) matchesAnyLine(pickleEvent *gherkin.PickleEvent) bool {
	uriLines, ok := p.lines[pickleEvent.URI]
	if !ok || len(uriLines) == 0 {
		return true
	}
	for _, line := range uriLines {
		for _, location := range pickleEvent.Pickle.Locations {
			if line == location.Line {
				return true
			}
		}
	}
	return false
}

func (p *PickleFilter) matchesAnyName(pickleEvent *gherkin.PickleEvent) bool {
	if len(p.nameRegexps) == 0 {
		return true
	}
	for _, nameRegexp := range p.nameRegexps {
		if nameRegexp.MatchString(pickleEvent.Pickle.Name) {
			return true
		}
	}
	return false
}

func (p *PickleFilter) matchesTagExpression(pickleEvent *gherkin.PickleEvent) bool {
	tagNames := make([]string, len(pickleEvent.Pickle.Tags))
	for i, tag := range pickleEvent.Pickle.Tags {
		tagNames[i] = tag.Name
	}
	return p.tagExpression.Evaluate(tagNames)
}
