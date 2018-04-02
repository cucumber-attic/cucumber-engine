package dto

// FeaturesFilterConfig is the configuration for how to filter the features
type FeaturesFilterConfig struct {
	TagExpression string           `json:"tagExpression"`
	Names         []string         `json:"names"`
	Lines         map[string][]int `json:"lines"`
}
