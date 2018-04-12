package dto

// FeaturesConfig is the configuration for what features to run
type FeaturesConfig struct {
	// TODO add default language
	Order         *FeaturesOrder        `json:"order"`
	AbsolutePaths []string              `json:"absolutePaths"`
	Filters       *FeaturesFilterConfig `json:"filters"`
}
