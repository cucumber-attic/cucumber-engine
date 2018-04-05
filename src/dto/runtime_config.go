package dto

// RuntimeConfig is the configuration for the run
type RuntimeConfig struct {
	IsFailFast bool `json:"isFailFast"`
	IsDryRun   bool `json:"isDryRun"`
	IsStrict   bool `json:"isStrict"`
}
