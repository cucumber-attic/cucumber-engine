package dto

// PatternMatch is a match from the step pattern
type PatternMatch struct {
	Captures          []string `json:"capture"`
	ParameterTypeName string   `json:"parameterTypeName"`
}
