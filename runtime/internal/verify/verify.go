package verify

import (
	"fmt"
)

type VerificationResult struct {
	Passed  bool
	Metrics map[string]float64
	Report  string
}

type VerificationAdapter interface {
	Name() string
	Run(target string) (*VerificationResult, error)
}

type Engine struct {
	Adapters map[string]VerificationAdapter
}

func NewEngine() *Engine {
	return &Engine{
		Adapters: make(map[string]VerificationAdapter),
	}
}

func (e *Engine) Register(adapter VerificationAdapter) {
	e.Adapters[adapter.Name()] = adapter
}

func (e *Engine) RunAll(target string) map[string]*VerificationResult {
	results := make(map[string]*VerificationResult)
	for name, adapter := range e.Adapters {
		fmt.Printf("[Verification] Running %s adapter...\n", name)
		res, err := adapter.Run(target)
		if err != nil {
			fmt.Printf("[Verification] %s failed to run: %v\n", name, err)
			continue
		}
		results[name] = res
	}
	return results
}
