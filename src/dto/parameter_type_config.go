package dto

// ParameterTypeConfig is the config for a cucumber expressions parameter type
type ParameterTypeConfig struct {
	Name                 string   `json:"name"`
	Regexps              []string `json:"regexps"`
	PreferForRegexpMatch bool     `json:"preferForRegexpMatch"`
	UseForSnippets       bool     `json:"useForSnippets"`
}
