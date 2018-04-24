package formatter

import (
	"fmt"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/cucumber/cucumber-pickle-runner/src/dto"
	"github.com/cucumber/cucumber-pickle-runner/src/dto/event"
)

type rerunFormatter struct {
	opts      *CommonOptions
	mapping   map[string][]int
	seperator string
}

// FormatAsRerun formats the events with the rerun formatter
func FormatAsRerun(opts *CommonOptions, seperator string) {
	r := &rerunFormatter{
		mapping:   map[string][]int{},
		opts:      opts,
		seperator: seperator,
	}
	go func() {
		for ev := range opts.EventChannel {
			switch e := ev.(type) {
			case *event.TestCaseFinished:
				r.storeFailedTestCases(e)
			case *event.TestRunFinished:
				r.logFailedTestCases(e)
			}
		}
		opts.Stream.Close()
	}()
}

func (r *rerunFormatter) storeFailedTestCases(e *event.TestCaseFinished) {
	if e.Result.Status != dto.StatusPassed {
		if _, ok := r.mapping[e.SourceLocation.URI]; !ok {
			r.mapping[e.SourceLocation.URI] = []int{}
		}
		r.mapping[e.SourceLocation.URI] = append(r.mapping[e.SourceLocation.URI], e.SourceLocation.Line)
	}
}

func (r *rerunFormatter) logFailedTestCases(e *event.TestRunFinished) {
	elements := []string{}
	for uri, lines := range r.mapping {
		relURI, err := filepath.Rel(r.opts.BaseDirectory, uri)
		if err != nil {
			panic(err)
		}
		elements = append(elements, fmt.Sprintf("%s:%s", relURI, joinInts(lines, ":")))
	}
	sort.Strings(elements)
	text := strings.Join(elements, r.seperator)
	r.opts.Stream.Write([]byte(text))
}

func joinInts(a []int, sep string) string {
	if len(a) == 0 {
		return ""
	}
	b := make([]string, len(a))
	for i, v := range a {
		b[i] = strconv.Itoa(v)
	}
	return strings.Join(b, sep)
}
