package t1api

import "testing"

func TestSumValid(t *testing.T) {
	got, err := Sum(1, 2)
	if err != nil {
		t.Fatal(err)
	}
	if got != 3 {
		t.Fatalf("got %d", got)
	}
}
