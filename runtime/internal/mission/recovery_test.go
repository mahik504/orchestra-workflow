package mission

import (
	"testing"
)

func TestRecoveryRetryableWithinBound(t *testing.T) {
	fail := Classify(ClassifyInput{Message: "build failed: does not compile", Stage: "Implement"})
	dec := Decide(DecideInput{
		Failure: fail,
		Budget:  Budget{RetryBudget: 2, Reasoning: ReasoningLow, Context: ContextScoped, Research: ResearchNone, VisualReview: VisualReviewNone, IndependentJudge: JudgeOff},
		Attempt: 1,
	})
	if !dec.Retry {
		t.Fatalf("expected retry: %+v", dec)
	}
	if dec.Change == ChangeNone {
		t.Fatalf("retry must change something")
	}
	hist := RecordAttempt(nil, 1, fail, dec, StatusRetrying)
	if len(hist) != 1 || hist[0].FailureClass != ClassImplementation {
		t.Fatalf("history %+v", hist)
	}
}

func TestRecoveryNonRetryableStops(t *testing.T) {
	fail := Classify(ClassifyInput{Message: "acceptance criteria cannot be interpreted"})
	dec := Decide(DecideInput{
		Failure: fail,
		Budget:  Budget{RetryBudget: 3, Reasoning: ReasoningMedium, Context: ContextScoped, Research: ResearchNone, VisualReview: VisualReviewNone, IndependentJudge: JudgeOff},
	})
	if dec.Retry {
		t.Fatalf("contract failure retried")
	}
	if !dec.Human || dec.Action != ActionHuman {
		t.Fatalf("expected human stop: %+v", dec)
	}
}

func TestRecoveryHumanNeededStops(t *testing.T) {
	fail := Classify(ClassifyInput{Message: "human approval gate required before writing code"})
	dec := Decide(DecideInput{Failure: fail, Budget: Budget{RetryBudget: 4, Reasoning: ReasoningLow, Context: ContextMinimal, Research: ResearchNone, VisualReview: VisualReviewNone, IndependentJudge: JudgeOff}})
	if dec.Retry || !dec.Human {
		t.Fatalf("%+v", dec)
	}
}

func TestRecoveryIdenticalFailuresTerminate(t *testing.T) {
	fail := Classify(ClassifyInput{Message: "npm ERR! code ERESOLVE could not resolve", Stage: "Implement"})
	budget := Budget{RetryBudget: 4, Reasoning: ReasoningLow, Context: ContextScoped, Research: ResearchNone, VisualReview: VisualReviewNone, IndependentJudge: JudgeOff}
	first := Decide(DecideInput{Failure: fail, Budget: budget, Attempt: 1})
	if !first.Retry {
		t.Fatalf("first should retry")
	}
	hist := RecordAttempt(nil, 1, fail, first, StatusRetrying)
	second := Decide(DecideInput{Failure: fail, Budget: budget, History: hist, Attempt: 2})
	if second.Retry {
		t.Fatalf("identical replay retried: %+v", second)
	}
	if !second.Human {
		t.Fatalf("identical replay should go to human")
	}
	hist = RecordAttempt(hist, 2, fail, second, StatusWaitingHuman)
	if len(hist) != 2 {
		t.Fatalf("history len %d", len(hist))
	}
	if hist[1].Signature != fail.Signature {
		t.Fatalf("signature not recorded")
	}
}

func TestRecoveryBudgetExhaustedStops(t *testing.T) {
	fail := Classify(ClassifyInput{Message: "timeout deadline exceeded"})
	hist := []AttemptRecord{
		{Attempt: 1, RecoveryAction: ActionRetry, Signature: "other"},
	}
	dec := Decide(DecideInput{
		Failure: fail,
		Budget:  Budget{RetryBudget: 1, Reasoning: ReasoningLow, Context: ContextScoped, Research: ResearchNone, VisualReview: VisualReviewNone, IndependentJudge: JudgeOff},
		History: hist,
		Attempt: 2,
	})
	if dec.Retry {
		t.Fatalf("exhausted budget retried: %+v", dec)
	}
}

func TestRecoveryDifferentFailureAfterRetryStillRetries(t *testing.T) {
	budget := Budget{RetryBudget: 3, Reasoning: ReasoningLow, Context: ContextScoped, Research: ResearchNone, VisualReview: VisualReviewNone, IndependentJudge: JudgeOff}
	first := Classify(ClassifyInput{Message: "build failed: does not compile", Stage: "Implement"})
	d1 := Decide(DecideInput{Failure: first, Budget: budget, Attempt: 1})
	if !d1.Retry {
		t.Fatalf("first should retry: %+v", d1)
	}
	hist := RecordAttempt(nil, 1, first, d1, StatusRetrying)
	second := Classify(ClassifyInput{Message: "timeout deadline exceeded", Stage: "Research"})
	if second.Signature == first.Signature {
		t.Fatalf("signatures collided")
	}
	d2 := Decide(DecideInput{Failure: second, Budget: budget, History: hist, Attempt: 2})
	if !d2.Retry {
		t.Fatalf("different failure should still retry: %+v", d2)
	}
	if d2.Change == ChangeNone {
		t.Fatalf("retry must change something")
	}
}

func TestRecoveryDoesNotMutateAcceptanceOrBar(t *testing.T) {
	m := validMission()
	fp := m.AcceptanceFingerprint
	bar := m.QualityBar
	fail := Classify(ClassifyInput{Message: "build failed: does not compile", Stage: "Implement"})
	_ = Decide(DecideInput{
		Failure: fail,
		Budget:  m.Budget,
	})
	if m.AcceptanceFingerprint != fp {
		t.Fatalf("acceptance mutated")
	}
	if m.QualityBar != bar {
		t.Fatalf("quality bar mutated")
	}
	if err := m.AcceptancePreserved(); err != nil {
		t.Fatal(err)
	}
}
