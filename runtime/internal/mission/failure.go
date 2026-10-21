package mission

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
)

const (
	ClassEnvironment    FailureClass = "environment"
	ClassDependency     FailureClass = "dependency"
	ClassProviderTool   FailureClass = "provider_tool"
	ClassImplementation FailureClass = "implementation"
	ClassTest           FailureClass = "test"
	ClassSpecification  FailureClass = "specification"
	ClassHuman          FailureClass = "human_needed"
	ClassTimeout        FailureClass = "timeout_resource"
)

type FailureClass string

type Failure struct {
	Class             FailureClass `json:"class"`
	What              string       `json:"what"`
	Why               string       `json:"why"`
	Retryable         bool         `json:"retryable"`
	ChangeBeforeRetry string       `json:"change_before_retry,omitempty"`
	HumanRequired     bool         `json:"human_required"`
	Signature         string       `json:"signature"`
	Stage             string       `json:"stage,omitempty"`
}

type FailureRecord struct {
	Class FailureClass `json:"class"`
	Note  string       `json:"note,omitempty"`
	Stage string       `json:"stage,omitempty"`
}

type ClassifyInput struct {
	Err      error
	Message  string
	ExitCode int
	Stage    string
	Signal   string
}

func Classify(in ClassifyInput) Failure {
	msg := strings.TrimSpace(in.Message)
	if msg == "" && in.Err != nil {
		msg = in.Err.Error()
	}
	low := strings.ToLower(msg)
	stage := strings.ToLower(strings.TrimSpace(in.Stage))
	signal := strings.ToLower(strings.TrimSpace(in.Signal))

	f := Failure{
		What:  firstNonEmpty(msg, "unspecified failure"),
		Stage: in.Stage,
	}

	switch {
	case signal == string(ClassHuman) || isHuman(low, in.Err):
		f.Class = ClassHuman
		f.Why = "a human decision, credential, or approval is required"
		f.Retryable = false
		f.HumanRequired = true
		f.ChangeBeforeRetry = string(ChangeNone)
	case signal == string(ClassSpecification) || isSpecification(low):
		f.Class = ClassSpecification
		f.Why = "the mission or contract cannot be executed without clarification"
		f.Retryable = false
		f.HumanRequired = true
		f.ChangeBeforeRetry = string(ChangeNone)
	case signal == string(ClassTimeout) || isTimeout(low):
		f.Class = ClassTimeout
		f.Why = "time or resource budget was exceeded"
		f.Retryable = !strings.Contains(low, "maximum qa iteration")
		f.HumanRequired = strings.Contains(low, "maximum qa iteration")
		f.ChangeBeforeRetry = string(ChangeContext)
	case signal == string(ClassEnvironment) || isEnvironment(low):
		f.Class = ClassEnvironment
		f.Why = "the execution environment is missing or misconfigured"
		f.Retryable = true
		f.ChangeBeforeRetry = string(ChangeEnvironment)
	case signal == string(ClassDependency) || isDependency(low):
		f.Class = ClassDependency
		f.Why = "a declared dependency is missing or failed to resolve"
		f.Retryable = true
		f.ChangeBeforeRetry = string(ChangeTool)
	case signal == string(ClassProviderTool) || isProviderTool(low):
		f.Class = ClassProviderTool
		f.Why = "an authorized tool or provider route failed"
		f.Retryable = true
		f.ChangeBeforeRetry = string(ChangeTool)
	case signal == string(ClassTest) || isTest(low, stage):
		f.Class = ClassTest
		f.Why = "verification evidence failed"
		f.Retryable = true
		f.ChangeBeforeRetry = string(ChangeInput)
	case signal == string(ClassImplementation) || isImplementation(low, stage):
		f.Class = ClassImplementation
		f.Why = "the implementation did not satisfy the assignment"
		f.Retryable = true
		f.ChangeBeforeRetry = string(ChangeApproach)
	default:
		if in.ExitCode != 0 && !isImplementation(low, stage) && !isTest(low, stage) {
			if stage == "discover" || strings.Contains(low, "setup") {
				f.Class = ClassEnvironment
				f.Why = "non-zero exit during setup is treated as environment, not implementation"
				f.Retryable = true
				f.ChangeBeforeRetry = string(ChangeEnvironment)
			} else {
				f.Class = ClassProviderTool
				f.Why = "non-zero exit without implementation evidence is treated as a tool failure"
				f.Retryable = true
				f.ChangeBeforeRetry = string(ChangeTool)
			}
		} else {
			f.Class = ClassImplementation
			f.Why = "unclassified failure defaulted to implementation"
			f.Retryable = true
			f.ChangeBeforeRetry = string(ChangeApproach)
		}
	}

	if f.HumanRequired {
		f.Retryable = false
		f.ChangeBeforeRetry = string(ChangeNone)
	}
	f.Signature = FailureSignature(f.Class, msg)
	return f
}

func FailureSignature(class FailureClass, msg string) string {
	norm := strings.ToLower(strings.Join(strings.Fields(msg), " "))
	sum := sha256.Sum256([]byte(string(class) + "|" + norm))
	return hex.EncodeToString(sum[:8])
}

func isHuman(low string, err error) bool {
	if err != nil && (errors.Is(err, errHumanSentinel) || strings.Contains(err.Error(), "human approval gate required")) {
		return true
	}
	needles := []string{
		"human approval gate required",
		"design lab gate is not cleared",
		"waiting for human",
		"credential",
		"authorization required",
		"destructive action",
		"irreversible deploy",
		"identity / brand",
		"human_needed",
		"retries exhausted",
		"material model",
		"tool / host",
		"host decision missing",
		"unauthorized host",
	}
	return containsAny(low, needles)
}

func isSpecification(low string) bool {
	return containsAny(low, []string{
		"acceptance criteria",
		"ambiguous requirement",
		"ambiguous product",
		"contract conflict",
		"specification",
		"cannot be interpreted",
		"mission/spec",
		"clarification required",
	})
}

func isTimeout(low string) bool {
	return containsAny(low, []string{
		"timeout",
		"deadline exceeded",
		"context canceled",
		"out of memory",
		"killed",
		"maximum qa iteration",
		"resource exhausted",
	})
}

func isEnvironment(low string) bool {
	return containsAny(low, []string{
		"no such file or directory",
		"enoent",
		"command not found",
		"not a git repository",
		"workspace root",
		"environment",
		"setup failed",
	})
}

func isDependency(low string) bool {
	return containsAny(low, []string{
		"cannot find package",
		"module not found",
		"npm err",
		"could not resolve",
		"missing go.sum",
		"dependency",
		"package.json",
	})
}

func isProviderTool(low string) bool {
	return containsAny(low, []string{
		"mcp",
		"tool failed",
		"provider",
		"connection refused",
		"status 5",
		"rpc error",
	})
}

func isTest(low, stage string) bool {
	if strings.Contains(stage, "visualqa") || strings.Contains(stage, "verify") {
		return true
	}
	return containsAny(low, []string{
		"test failed",
		"tests failed",
		"assertion",
		"fail:",
		"expected",
		"visual qa violations",
	})
}

func isImplementation(low, stage string) bool {
	if stage == "implement" {
		return !isHuman(low, nil) && !isSpecification(low)
	}
	return containsAny(low, []string{
		"does not compile",
		"build failed",
		"layout/code",
		"token/style",
		"implementation",
	})
}

func containsAny(low string, needles []string) bool {
	for _, n := range needles {
		if strings.Contains(low, n) {
			return true
		}
	}
	return false
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

var errHumanSentinel = errors.New("human approval gate required before writing code")
