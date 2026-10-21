package t3debug

import "testing"

func TestSum(t *testing.T) {
	if got := Sum(1, 2); got != 3 {
		t.Fatalf("test failed: Sum(1,2) = 3, got %d", got)
	}
}
