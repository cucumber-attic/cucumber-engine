package dto

// Pattern is how the step definition matches text
type Pattern struct {
	Source string      `json:"source"`
	Type   PatternType `json:"type"`
}
