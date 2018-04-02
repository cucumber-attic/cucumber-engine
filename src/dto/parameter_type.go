package dto

// ParameterType is the a configuration for cucumber expressions
type ParameterType struct {
	Name    string   `json:"name"`
	Regexps []string `json:"regexps"`
}
