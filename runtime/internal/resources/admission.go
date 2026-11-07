package resources

import (
	"fmt"
	"strings"
)

// AdmissionOrder is catalog admission. It is not the run lifecycle in contract.go.
var AdmissionOrder = []string{
	"discovered",
	"candidate",
	"evaluated",
	"approved",
	"installed",
	"configured",
	"healthy",
	"verified",
	"monitored",
}

func Admit(from, to string) error {
	from = strings.TrimSpace(from)
	to = strings.TrimSpace(to)
	if from == "" || to == "" {
		return fmt.Errorf("admission states are required")
	}
	if to == "rejected" || to == "deprecated" {
		if !knownAdmission(from) && from != "rejected" && from != "deprecated" {
			return fmt.Errorf("unknown admission state %s", from)
		}
		return nil
	}
	fromIndex := admissionIndex(from)
	toIndex := admissionIndex(to)
	if fromIndex < 0 || toIndex < 0 {
		return fmt.Errorf("unknown admission transition %s -> %s", from, to)
	}
	if toIndex != fromIndex+1 {
		return fmt.Errorf("admission cannot jump %s -> %s", from, to)
	}
	return nil
}

func admissionIndex(state string) int {
	for i, name := range AdmissionOrder {
		if name == state {
			return i
		}
	}
	return -1
}

func knownAdmission(state string) bool {
	return admissionIndex(state) >= 0
}

type FitNotes struct {
	TriggerMatch bool
	LicenseOK    bool
	Maintained   bool
	Evidence     string
}

func ScoreFit(notes FitNotes) (int, error) {
	evidence := strings.ToLower(strings.TrimSpace(notes.Evidence))
	if evidence == "" || strings.Contains(evidence, "none recorded") {
		return 0, fmt.Errorf("fit score needs evidence")
	}
	score := 0
	if notes.TriggerMatch {
		score++
	}
	if notes.LicenseOK {
		score++
	}
	if notes.Maintained {
		score++
	}
	return score, nil
}
