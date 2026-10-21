package benchmark

import (
	"fmt"
	"strings"

	"github.com/mahik504/orchestra-workflow/runtime/internal/mission"
)

const seededTestFailure = "test failed: Sum(1,2) = 3, got 4"

// RecoverSeeded uses the existing classifier and retry rule.
// The first decision must change something. The identical failure must not retry.
func RecoverSeeded(budget mission.Budget) (Record, error) {
	if budget.RetryBudget < 1 {
		budget.RetryBudget = 1
	}
	fail := mission.Classify(mission.ClassifyInput{
		Message: seededTestFailure,
		Stage:   "verify",
		Signal:  string(mission.ClassTest),
	})
	first := mission.Decide(mission.DecideInput{Failure: fail, Budget: budget, Attempt: 1})
	if !first.Retry || first.Change == mission.ChangeNone {
		return Record{}, fmt.Errorf("seeded failure did not retry with a change: %+v", first)
	}
	history := mission.RecordAttempt(nil, 1, fail, first, mission.StatusRetrying)
	second := mission.Decide(mission.DecideInput{Failure: fail, Budget: budget, History: history, Attempt: 2})
	if second.Retry {
		return Record{}, fmt.Errorf("identical failure retried")
	}
	other := mission.Classify(mission.ClassifyInput{
		Message: "test failed: Sum(2,2) = 4, got 5",
		Stage:   "verify",
		Signal:  string(mission.ClassTest),
	})
	probe := mission.Decide(mission.DecideInput{Failure: other, Budget: budget, History: history, Attempt: 2})
	left := budget.RetryBudget - 1
	if !strings.Contains(probe.Reason, fmt.Sprintf("remaining=%d", left)) {
		return Record{}, fmt.Errorf("retry budget was not consumed: %s", probe.Reason)
	}
	rec := Record{
		Phase:                 PhasePilot,
		Claim:                 ClaimNone,
		Tokens:                TokensUnavailable,
		TaskExecution:         TaskCostUnavailable,
		RecoveryClass:         string(fail.Class),
		RecoveryChange:        string(first.Change),
		Retries:               1,
		IdenticalRetryStopped: !second.Retry && second.Human,
		HumanGates:            1,
		Limitations:           []string{fmt.Sprintf("one changed attempt leaves remaining %d", left)},
	}
	if fail.Class != mission.ClassTest {
		return rec, fmt.Errorf("class %s", fail.Class)
	}
	return rec, nil
}
