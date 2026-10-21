package benchmark

type Task struct {
	ID      string `json:"id"`
	Class   string `json:"class"`
	Enabled bool   `json:"enabled"`
	Type    string `json:"type"`
	Request string `json:"request"`
}

func Manifest() []Task {
	return []Task{
		{ID: "t1-api", Class: "T1", Enabled: true, Type: "FEATURE", Request: "Add a small API sum that validates two integers and returns an error for invalid input. Include a test."},
		{ID: "t2-crud", Class: "T2", Enabled: false, Type: "FEATURE"},
		{ID: "t3-debug", Class: "T3", Enabled: true, Type: "BUGFIX", Request: "Fix Sum. The test expects a+b. Change the return. Do not add files."},
		{ID: "t4-modify", Class: "T4", Enabled: false, Type: "FEATURE"},
		{ID: "t5-responsive", Class: "T5", Enabled: false, Type: "DESIGN"},
		{ID: "t6-product", Class: "T6", Enabled: false, Type: "DESIGN"},
		{ID: "t7-flagship", Class: "T7", Enabled: false, Type: "DESIGN"},
		{ID: "t8-research", Class: "T8", Enabled: false, Type: "FEATURE"},
		{ID: "t9-reverse", Class: "T9", Enabled: false, Type: "FEATURE"},
		{ID: "t10-multi", Class: "T10", Enabled: false, Type: "FEATURE"},
	}
}

func EnabledTasks() []Task {
	var out []Task
	for _, t := range Manifest() {
		if t.Enabled {
			out = append(out, t)
		}
	}
	return out
}

// ExpansionOpen reports whether a later pass may enable harder classes.
// This pilot does not flip those manifest flags.
func ExpansionOpen(tracesOK, isolated, scored, recovered bool) bool {
	return tracesOK && isolated && scored && recovered
}
