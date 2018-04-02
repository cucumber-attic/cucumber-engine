package dto

// FeaturesConfig is the configuration for what features to run
type FeaturesConfig struct {
	AbsolutePaths []string              `json:"absolutePaths"`
	Filters       *FeaturesFilterConfig `json:"filters"`
}
