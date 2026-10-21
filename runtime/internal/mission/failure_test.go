package mission

import (
	"errors"
	"testing"
)

func TestClassifyRepresentativeFailures(t *testing.T) {
	cases := []struct {
		name  string
		in    ClassifyInput
		class FailureClass
		human bool
		retry bool
	}{
		{
			name:  "environment",
			in:    ClassifyInput{Message: "open go.mod: no such file or directory", Stage: "Discover"},
			class: ClassEnvironment, retry: true,
		},
		{
			name:  "dependency",
			in:    ClassifyInput{Message: "npm ERR! code ERESOLVE could not resolve", Stage: "Implement"},
			class: ClassDependency, retry: true,
		},
		{
			name:  "provider tool",
			in:    ClassifyInput{Message: "mcp playwright: connection refused", Stage: "VisualQA"},
			class: ClassProviderTool, retry: true,
		},
		{
			name:  "implementation",
			in:    ClassifyInput{Message: "build failed: does not compile", Stage: "Implement"},
			class: ClassImplementation, retry: true,
		},
		{
			name:  "test",
			in:    ClassifyInput{Message: "--- FAIL: TestHealth (0.01s) assertion", Stage: "VisualQA"},
			class: ClassTest, retry: true,
		},
		{
			name:  "specification",
			in:    ClassifyInput{Message: "acceptance criteria cannot be interpreted safely", Stage: "Classify"},
			class: ClassSpecification, human: true, retry: false,
		},
		{
			name:  "human sentinel",
			in:    ClassifyInput{Err: errors.New("human approval gate required before writing code"), Stage: "Synthesize"},
			class: ClassHuman, human: true, retry: false,
		},
		{
			name:  "design lab gate",
			in:    ClassifyInput{Message: "design lab gate is not cleared: refusing to write App.tsx (PREMIUM)", Stage: "Implement"},
			class: ClassHuman, human: true, retry: false,
		},
		{
			name:  "timeout",
			in:    ClassifyInput{Message: "context deadline exceeded: timeout", Stage: "Research"},
			class: ClassTimeout, retry: true,
		},
		{
			name:  "nonzero setup exit is not implementation",
			in:    ClassifyInput{Message: "setup failed", ExitCode: 1, Stage: "Discover"},
			class: ClassEnvironment, retry: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := Classify(tc.in)
			if got.Class != tc.class {
				t.Fatalf("class=%s want %s why=%s", got.Class, tc.class, got.Why)
			}
			if got.HumanRequired != tc.human {
				t.Fatalf("human=%v want %v", got.HumanRequired, tc.human)
			}
			if got.Retryable != tc.retry {
				t.Fatalf("retryable=%v want %v", got.Retryable, tc.retry)
			}
			if got.Signature == "" {
				t.Fatalf("missing signature")
			}
		})
	}
}

func TestClassifyContractDoesNotBecomeImplementation(t *testing.T) {
	got := Classify(ClassifyInput{
		Message: "specification/contract conflict on checkout fees",
		Stage:   "Implement",
	})
	if got.Class != ClassSpecification {
		t.Fatalf("got %s", got.Class)
	}
	if got.Retryable {
		t.Fatalf("contract failures must not be retryable")
	}
}

func TestClassifyHumanStaysHuman(t *testing.T) {
	got := Classify(ClassifyInput{Message: "missing critical credential for deploy"})
	if got.Class != ClassHuman || !got.HumanRequired || got.Retryable {
		t.Fatalf("%+v", got)
	}
}

func TestClassifyHumanNeededSituations(t *testing.T) {
	cases := []struct {
		msg   string
		class FailureClass
		human bool
		retry bool
	}{
		{"ambiguous product behavior in checkout fees", ClassSpecification, true, false},
		{"missing critical credential for deploy", ClassHuman, true, false},
		{"destructive action requires approval before drop table", ClassHuman, true, false},
		{"specification/contract conflict on tax rules", ClassSpecification, true, false},
		{"acceptance criteria cannot be interpreted safely", ClassSpecification, true, false},
		{"retries exhausted; human direction required", ClassHuman, true, false},
		{"material model / tool / host decision missing", ClassHuman, true, false},
	}
	for _, tc := range cases {
		t.Run(tc.msg, func(t *testing.T) {
			got := Classify(ClassifyInput{Message: tc.msg})
			if got.Class != tc.class {
				t.Fatalf("class=%s want %s why=%s", got.Class, tc.class, got.Why)
			}
			if got.HumanRequired != tc.human || got.Retryable != tc.retry {
				t.Fatalf("human=%v retry=%v want human=%v retry=%v (%+v)", got.HumanRequired, got.Retryable, tc.human, tc.retry, got)
			}
			dec := Decide(DecideInput{
				Failure: got,
				Budget:  Budget{RetryBudget: 4, Reasoning: ReasoningLow, Context: ContextScoped, Research: ResearchNone, VisualReview: VisualReviewNone, IndependentJudge: JudgeOff},
			})
			if dec.Retry {
				t.Fatalf("auto-retried: %+v", dec)
			}
			if !dec.Human {
				t.Fatalf("expected human stop: %+v", dec)
			}
		})
	}
}
