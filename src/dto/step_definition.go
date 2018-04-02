package dto

// StepDefinition is the implementation of a step
type StepDefinition struct {
	ID      string  `json:"id"`
	Pattern Pattern `json:"pattern"`
	URI     string  `json:"uri"`
	Line    int     `json:"line"`
}
