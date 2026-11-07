package benchmark

// ComparisonTask is a paired eval that has not been run.
// Enabled stays false until a real Orchestra arm and a raw host arm both exist.
type ComparisonTask struct {
	ID       string
	Class    string
	Enabled  bool
	Claim    string
	Measures []string
	Request  string
}

// ComparisonSuite names the six classes from the 4.1 brief.
// It does not score them and it does not change the T1/T3 pilot.
func ComparisonSuite() []ComparisonTask {
	measures := []string{
		"requirement_completion",
		"visual_quality",
		"interaction_correctness",
		"code_quality",
		"architectural_correctness",
		"performance",
		"accessibility",
		"regression_count",
		"time",
		"tool_cost",
		"token_cost",
	}
	requests := []struct {
		id, class, request string
	}{
		{"c1-premium-landing", "premium-landing", "A premium landing page with a stated audience and one reference principle."},
		{"c2-immersive-3d", "immersive-3d", "One purposeful 3D scene with a still fallback and a DPR cap."},
		{"c3-fullstack-saas", "fullstack-saas", "A full-stack SaaS slice with accounts, one schema, and one authenticated route."},
		{"c4-fastapi-ai", "fastapi-ai", "A FastAPI inference route with validation, tests, and no second runtime."},
		{"c5-backend-only", "backend-only", "A backend change with a migration, a constraint, and no UI."},
		{"c6-bugfix", "bugfix", "A bug fix in an existing project with a regression test and no new framework."},
	}
	out := make([]ComparisonTask, 0, len(requests))
	for _, item := range requests {
		copied := append([]string{}, measures...)
		out = append(out, ComparisonTask{
			ID:       item.id,
			Class:    item.class,
			Enabled:  false,
			Claim:    ClaimNone,
			Measures: copied,
			Request:  item.request,
		})
	}
	return out
}
