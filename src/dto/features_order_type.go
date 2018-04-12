package dto

// FeaturesOrderType is an enumeration of the available values for the Type field in the FeaturesOrder struct
type FeaturesOrderType string

const (
	// FeaturesOrderTypeDefined is for the order of definition
	FeaturesOrderTypeDefined = FeaturesOrderType("defined")
	// FeaturesOrderTypeRandom is for a random order (with a given seed)
	FeaturesOrderTypeRandom = FeaturesOrderType("random")
)
