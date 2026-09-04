package planner

import (
	"fmt"
	"github.com/user/orchestra-v3/internal/classifier"
	"github.com/user/orchestra-v3/internal/resources"
	"sync"
)

type ExecutionNode struct {
	ID           string
	Task         *classifier.Task
	Capabilities []*resources.Capability
	Dependencies []string
	Status       string // PENDING, RUNNING, COMPLETED, FAILED
}

type DAGPlanner struct {
	Nodes map[string]*ExecutionNode
	mu    sync.Mutex
}

func NewDAGPlanner() *DAGPlanner {
	return &DAGPlanner{
		Nodes: make(map[string]*ExecutionNode),
	}
}

func (p *DAGPlanner) AddNode(node *ExecutionNode) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Nodes[node.ID] = node
}

func (p *DAGPlanner) Execute() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	fmt.Println("Executing DAG Plan in parallel...")

	var wg sync.WaitGroup
	errCh := make(chan error, len(p.Nodes))

	// In a real DAG, we'd topologically sort and run nodes whose dependencies are met.
	// For now, we simulate parallel execution of independent nodes.
	for _, node := range p.Nodes {
		if len(node.Dependencies) == 0 {
			wg.Add(1)
			go func(n *ExecutionNode) {
				defer wg.Done()
				fmt.Printf("Starting execution of node: %s\n", n.ID)

				// Lazy load capabilities just before execution
				for _, cap := range n.Capabilities {
					if err := cap.LoadDetails(); err != nil {
						fmt.Printf("Warning: failed to lazy load capability %s: %v\n", cap.ID, err)
					}
				}

				n.Status = "COMPLETED"
				fmt.Printf("Finished execution of node: %s\n", n.ID)
			}(node)
		}
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			return err
		}
	}

	return nil
}
