package dto

// PatternType is an enumeration of the available values for the Type field in the Pattern struct
type PatternType string

const (
	// PatternTypeCucumberExpression is for a cucumber expression
	PatternTypeCucumberExpression = PatternType("cucumber_expression")
	// PatternTypeRegularExpression is for a regular expression
	PatternTypeRegularExpression = PatternType("regular_expression")
)
