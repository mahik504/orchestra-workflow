package kernel

import (
	"fmt"
	"github.com/user/orchestra-v3/internal/classifier"
	"github.com/user/orchestra-v3/internal/resources"
	"github.com/user/orchestra-v3/internal/router"
)

type Kernel struct {
	Registry   *resources.Registry
	Classifier *classifier.Classifier
	Router     *router.Router
	Allocator  *router.Allocator
}

func NewKernel() *Kernel {
	reg := resources.NewRegistry()
	cls := classifier.NewClassifier()
	rtr := router.NewRouter(reg)
	alloc := router.NewAllocator()

	return &Kernel{
		Registry:   reg,
		Classifier: cls,
		Router:     rtr,
		Allocator:  alloc,
	}
}

// ProcessRequest executes the primary vertical slice:
// PRD -> Task -> Classify -> Capability Selection -> Plan -> Approval Gate
func (k *Kernel) ProcessRequest(raw string) error {
	fmt.Println("Orchestra Intake: Processing Request...")
	
	task, err := k.Classifier.Classify(raw)
	if err != nil {
		return fmt.Errorf("classification failed: %w", err)
	}
	
	fmt.Printf("Classified Task [%s]: Visual=%v, Security=%v\n", task.Type, task.RequiresVisual, task.RequiresSecurity)

	allocation := k.Allocator.Allocate(task)
	if task.UserOverride != nil && task.UserOverride.ForceAgent != "" {
		fmt.Printf("User Overrode Agent Selection. Forcing: %s (was %s)\n", task.UserOverride.ForceAgent, allocation.PrimaryAgent)
		allocation.PrimaryAgent = router.AgentType(task.UserOverride.ForceAgent)
	}
	fmt.Printf("Allocation: Agent=%s, Model=%s, Mode=%s\n", allocation.PrimaryAgent, allocation.Model, allocation.Mode)

	plan := k.Router.Compose(task)
	fmt.Printf("Composed Plan: %d capabilities selected (Est Cost: %.2f)\n", len(plan.SelectedCapabilities), plan.EstimatedTokenCost)
	for _, cap := range plan.SelectedCapabilities {
		fmt.Printf(" - %s (%s)\n", cap.Name, cap.Category)
	}

	if plan.RequiresHumanGate {
		fmt.Printf("[GATE] HUMAN APPROVAL REQUIRED: %s\n", plan.ApprovalReason)
		// Gate blocks here until user approves
	} else {
		fmt.Println("[GATE] No manual approval gates triggered. Proceeding autonomously.")
	}

	fmt.Println("Handoff state generation starting...")
	// call internal/handoff logic here
	
	return nil
}
