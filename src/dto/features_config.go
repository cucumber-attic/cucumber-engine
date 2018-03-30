package dto

// FeaturesFilterConfig is the configuration for how to filter the features
type FeaturesFilterConfig struct {
	TagExpression string   `json:"tagExpression"`
	Names         []string `json:"names"`
}

// FeaturesConfig is the configuration for what features to run
type FeaturesConfig struct {
	AbsolutePaths []string             `json:"absolutePaths"`
	Filters       FeaturesFilterConfig `json:"filters"`
}
