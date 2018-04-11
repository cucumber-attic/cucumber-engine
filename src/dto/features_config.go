package dto

// FeaturesConfig is the configuration for what features to run
type FeaturesConfig struct {
	// TODO add default language
	// TODO add order (random vs defined plus seed)
	AbsolutePaths []string              `json:"absolutePaths"`
	Filters       *FeaturesFilterConfig `json:"filters"`
}
