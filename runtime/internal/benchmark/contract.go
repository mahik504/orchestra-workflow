package benchmark

// Hard dimensions are fixed before any run. They are not an overall score.
const (
	DimRequiredBehavior = "required_behavior"
	DimBuild            = "build"
	DimTests            = "tests"
	DimAcceptance       = "acceptance"
	DimScope            = "scope"

	TokensUnavailable   = "unavailable"
	TaskCostUnavailable = "unavailable"
	ClaimNone           = "none"
	PhasePilot          = "pilot"
)

var HardDimensions = []string{
	DimRequiredBehavior,
	DimBuild,
	DimTests,
	DimAcceptance,
	DimScope,
}

type Check struct {
	Name string `json:"name"`
	Pass bool   `json:"pass"`
}

type Judgment struct {
	Note     string `json:"note,omitempty"`
	Positive bool   `json:"positive"`
}

// ApplyScore records hard checks. A positive judgment cannot clear a hard fail.
func ApplyScore(rec *Record, checks []Check, judgment *Judgment) {
	if rec == nil {
		return
	}
	rec.Deterministic = append([]Check{}, checks...)
	rec.Judgment = judgment
	rec.HardFail = false
	rec.TestFailures = 0
	rec.AcceptanceFailures = 0
	for _, c := range checks {
		if c.Pass {
			continue
		}
		rec.HardFail = true
		switch c.Name {
		case DimTests:
			rec.TestFailures++
		case DimAcceptance, DimRequiredBehavior:
			rec.AcceptanceFailures++
		case DimScope:
			rec.ScopeFailures++
		}
	}
}
