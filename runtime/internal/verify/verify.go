package verify

type VerificationResult struct {
	Passed  bool
	Metrics map[string]float64
	Report  string
}

type VerificationAdapter interface {
	Name() string
	Run(target string) (*VerificationResult, error)
}
