package formatter

import (
	"io"

	gherkin "github.com/cucumber/gherkin-go"
)

// CommonOptions are the options passed to each FormatAs function
type CommonOptions struct {
	EventChannel  chan gherkin.CucumberEvent
	Stream        io.WriteCloser
	BaseDirectory string
}
