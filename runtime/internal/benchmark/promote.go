package benchmark

// PromoteSkill decides whether a skill may join a default Load line.
// successDelta is the Orchestra-with-skill rate minus the control-plane rate.
// costRatio is that arm's cost divided by the control-plane cost.
// Fewer than three trials is not a promotion. Claims stay off until a record exists.
func PromoteSkill(successDelta, costRatio float64, trials int) (bool, string) {
	if trials < 3 {
		return false, "fewer than 3 trials"
	}
	if successDelta < 0 {
		return false, "success change is negative"
	}
	if costRatio <= 0 {
		return false, "cost was not measured"
	}
	if costRatio > 1.10 {
		return false, "cost rose by more than 10 percent"
	}
	return true, "promote"
}
