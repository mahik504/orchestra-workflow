package benchmark

import "testing"

func TestPromoteSkill(t *testing.T) {
	if ok, _ := PromoteSkill(0.1, 1.05, 3); !ok {
		t.Fatal("a non-negative lift inside the cost cap should promote")
	}
	if ok, _ := PromoteSkill(-0.01, 1.0, 5); ok {
		t.Fatal("a negative success change must not promote")
	}
	if ok, _ := PromoteSkill(0.2, 1.2, 5); ok {
		t.Fatal("a cost rise over 10 percent must not promote")
	}
	if ok, reason := PromoteSkill(0, 1, 2); ok || reason != "fewer than 3 trials" {
		t.Fatalf("short run promoted: %v %s", ok, reason)
	}
}
