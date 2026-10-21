package mission

import (
	"fmt"
	"time"
)

const (
	ActionRetry     RecoveryAction = "retry"
	ActionStop      RecoveryAction = "stop"
	ActionHuman     RecoveryAction = "human_gate"
	ActionRepairEnv RecoveryAction = "repair_environment"
)

const (
	ChangeNone        RecoveryChange = "none"
	ChangeInput       RecoveryChange = "input"
	ChangeContext     RecoveryChange = "context"
	ChangeStrategy    RecoveryChange = "strategy"
	ChangeTool        RecoveryChange = "tool"
	ChangeApproach    RecoveryChange = "approach"
	ChangeEnvironment RecoveryChange = "environment"
)

type RecoveryAction string
type RecoveryChange string

type AttemptRecord struct {
	Attempt         int            `json:"attempt"`
	FailureClass    FailureClass   `json:"failure_class"`
	PreviousFailure string         `json:"previous_failure,omitempty"`
	RecoveryAction  RecoveryAction `json:"recovery_action"`
	Change          RecoveryChange `json:"change"`
	ResultingState  MissionStatus  `json:"resulting_state"`
	At              string         `json:"at"`
	Signature       string         `json:"signature,omitempty"`
	Reason          string         `json:"reason,omitempty"`
}

type Decision struct {
	Action         RecoveryAction `json:"action"`
	Retry          bool           `json:"retry"`
	Human          bool           `json:"human"`
	Change         RecoveryChange `json:"change"`
	Reason         string         `json:"reason"`
	AdjustedBudget *Budget        `json:"adjusted_budget,omitempty"`
	NextAction     string         `json:"next_action,omitempty"`
}

type DecideInput struct {
	Failure Failure
	Budget  Budget
	History []AttemptRecord
	Attempt int
}

func Decide(in DecideInput) Decision {
	used := retryCount(in.History)
	remaining := in.Budget.RetryBudget - used
	if remaining < 0 {
		remaining = 0
	}

	if identicalReplay(in.Failure, in.History) {
		return Decision{
			Action:     ActionHuman,
			Retry:      false,
			Human:      true,
			Change:     ChangeNone,
			Reason:     "identical failure repeated; automatic retry stopped",
			NextAction: "await_human",
		}
	}

	switch in.Failure.Class {
	case ClassHuman:
		return Decision{
			Action:     ActionHuman,
			Human:      true,
			Change:     ChangeNone,
			Reason:     "human-needed failure terminates automatic execution",
			NextAction: "await_human",
		}
	case ClassSpecification:
		return Decision{
			Action:     ActionHuman,
			Human:      true,
			Change:     ChangeNone,
			Reason:     "specification/contract failure requires mission clarification, not a blind implementation change",
			NextAction: "await_human",
		}
	}

	if !in.Failure.Retryable || remaining <= 0 {
		action := ActionStop
		human := in.Failure.HumanRequired
		next := "stop"
		if human {
			action = ActionHuman
			next = "await_human"
		}
		return Decision{
			Action:     action,
			Human:      human,
			Change:     ChangeNone,
			Reason:     fmt.Sprintf("not retryable or retry budget exhausted (remaining=%d)", remaining),
			NextAction: next,
		}
	}

	change := RecoveryChange(in.Failure.ChangeBeforeRetry)
	if change == "" {
		change = defaultChange(in.Failure.Class)
	}
	if change == ChangeNone {
		return Decision{
			Action:     ActionStop,
			Change:     ChangeNone,
			Reason:     "retry refused because no meaningful change was identified",
			NextAction: "stop",
		}
	}

	adjusted := in.Budget
	if in.Failure.Class == ClassTimeout {
		if adjusted.Context == ContextExtended {
			adjusted.Context = ContextScoped
		}
		adjusted.Reasoning = bumpReasoning(adjusted.Reasoning)
	}

	action := ActionRetry
	next := "retry_with_" + string(change)
	if in.Failure.Class == ClassEnvironment {
		action = ActionRepairEnv
		next = "repair_environment"
	}

	return Decision{
		Action:         action,
		Retry:          true,
		Change:         change,
		Reason:         fmt.Sprintf("retry %s after %s (remaining=%d)", in.Failure.Class, change, remaining),
		AdjustedBudget: &adjusted,
		NextAction:     next,
	}
}

func RecordAttempt(history []AttemptRecord, attempt int, fail Failure, dec Decision, state MissionStatus) []AttemptRecord {
	rec := AttemptRecord{
		Attempt:         attempt,
		FailureClass:    fail.Class,
		PreviousFailure: fail.What,
		RecoveryAction:  dec.Action,
		Change:          dec.Change,
		ResultingState:  state,
		At:              time.Now().UTC().Format(time.RFC3339),
		Signature:       fail.Signature,
		Reason:          dec.Reason,
	}
	return append(append([]AttemptRecord{}, history...), rec)
}

func retryCount(history []AttemptRecord) int {
	n := 0
	for _, h := range history {
		if h.RecoveryAction == ActionRetry || h.RecoveryAction == ActionRepairEnv {
			n++
		}
	}
	return n
}

func identicalReplay(fail Failure, history []AttemptRecord) bool {
	if fail.Signature == "" || len(history) == 0 {
		return false
	}
	same := 0
	for i := len(history) - 1; i >= 0; i-- {
		if history[i].Signature == fail.Signature {
			same++
			continue
		}
		break
	}
	return same >= 1
}

func defaultChange(class FailureClass) RecoveryChange {
	switch class {
	case ClassEnvironment:
		return ChangeEnvironment
	case ClassDependency, ClassProviderTool:
		return ChangeTool
	case ClassImplementation:
		return ChangeApproach
	case ClassTest:
		return ChangeInput
	case ClassTimeout:
		return ChangeContext
	default:
		return ChangeNone
	}
}
