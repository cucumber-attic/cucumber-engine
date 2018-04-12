package dto

// FeaturesOrder is the configuration for what order to run the features in
type FeaturesOrder struct {
	Type FeaturesOrderType `json:"type"`
	Seed int64             `json:"seed"`
}
