package visual

import "testing"

func TestHealRequiresConcreteChange(t *testing.T) {
	_, ok := HealDirective(nil, 1)
	if ok {
		t.Fatal("no findings must not retry")
	}
	dir, ok := HealDirective([]Finding{{
		Type: "overflow", Layer: "A", Viewport: "mobile", Trivial: true,
		Explanation: "overflow", Action: "Constrain the overflowing row on mobile.",
	}}, 1)
	if !ok || dir == "" {
		t.Fatal("actionable finding must produce a change")
	}
	if MapFailureClass([]Finding{{Type: "overflow", Layer: "A"}}) != "LAYOUT_CODE" {
		t.Fatal("overflow maps to layout")
	}
	if !AllTrivial([]Finding{{Type: "overflow", Layer: "A", Trivial: true}}) {
		t.Fatal("overflow is trivial QA")
	}
	if AllTrivial([]Finding{{Type: "brand", Layer: "B"}}) {
		t.Fatal("brand is not trivial")
	}
}
