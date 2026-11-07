package resources

import "testing"

func TestAdmissionNextStepOnly(t *testing.T) {
	if err := Admit("discovered", "candidate"); err != nil {
		t.Fatal(err)
	}
	if err := Admit("discovered", "installed"); err == nil {
		t.Fatal("jump was allowed")
	}
	if err := Admit("approved", "rejected"); err != nil {
		t.Fatal(err)
	}
	if err := Admit("monitored", "deprecated"); err != nil {
		t.Fatal(err)
	}
}

func TestScoreFitNeedsEvidence(t *testing.T) {
	if _, err := ScoreFit(FitNotes{TriggerMatch: true}); err == nil {
		t.Fatal("missing evidence scored a fit")
	}
	score, err := ScoreFit(FitNotes{TriggerMatch: true, LicenseOK: true, Evidence: "MIT LICENSE in the repo"})
	if err != nil {
		t.Fatal(err)
	}
	if score != 2 {
		t.Fatalf("score %d", score)
	}
	if _, err := ScoreFit(FitNotes{
		TriggerMatch: true,
		LicenseOK:    true,
		Maintained:   true,
		Evidence:     "none recorded in this registry",
	}); err == nil {
		t.Fatal("sentinel evidence scored a fit")
	}
}
