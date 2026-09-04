package verify

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/user/orchestra-v3/internal/classifier"
)

func pendingLab(t *testing.T) *DesignLab {
	t.Helper()
	return NewDesignLab(&classifier.Brief{
		TaskID:            "t1",
		DesignLabRequired: true,
		DesignLabReason:   "PREMIUM visual work",
		QualityBar:        classifier.BarPremium,
	}, t.TempDir())
}

func twoDirections() []Direction {
	return []Direction{
		{
			ID: "a", Concept: "Warm editorial",
			Typography: "Freight Display + Söhne", TypographySrc: "Klim type specimen",
			ColorWorld: "Roasted umber on bone", ColorSrc: "1970s coffee packaging archive",
			LayoutLanguage: "Asymmetric editorial grid", ComponentKit: "custom",
			MotionEngine: "CSS transitions", MotionWhy: "page is mostly static, no timeline needed",
			LogoMethod: "wordmark", IconSystem: "Phosphor", Stack: []string{"astro", "tailwind"},
		},
		{
			ID: "b", Concept: "Stark industrial",
			Typography: "Diatype + Diatype Mono", TypographySrc: "Dinamo specimen",
			ColorWorld: "Cold steel with ember accent", ColorSrc: "Braun product photography",
			LayoutLanguage: "Strict 12-column", ComponentKit: "custom",
			MotionEngine: "GSAP", MotionWhy: "scroll-linked reveals need a timeline",
			LogoMethod: "monogram", IconSystem: "Lucide", Stack: []string{"next", "tailwind"},
		},
	}
}

// The gate is a lock. While it is pending, nothing the browser renders gets written.
func TestGate_BlocksFrontendWritesWhilePending(t *testing.T) {
	lab := pendingLab(t)
	if lab.State != GatePending {
		t.Fatalf("state = %s, want PENDING", lab.State)
	}
	if lab.Cleared() {
		t.Fatal("a pending gate reported itself cleared")
	}

	blocked := []string{
		"src/App.tsx", "styles/main.css", "index.html", "components/Hero.jsx",
		"app.vue", "page.svelte", "shaders/water.frag", "tailwind.config.js",
		"src/tokens.ts", "styles/globals.js",
	}
	for _, p := range blocked {
		err := lab.GuardWrite(p)
		var gateErr *ErrGateNotCleared
		if !errors.As(err, &gateErr) {
			t.Errorf("GuardWrite(%q) = %v, want a gate error", p, err)
		}
	}

	allowed := []string{
		"main.go", "server/api.py", "README.md", "DESIGN.md",
		"package.json", "notes/brief.txt", "migrations/001.sql",
	}
	for _, p := range allowed {
		if err := lab.GuardWrite(p); err != nil {
			t.Errorf("GuardWrite(%q) blocked non-frontend work: %v", p, err)
		}
	}
}

func TestGate_NotRequiredNeverBlocks(t *testing.T) {
	lab := NewDesignLab(&classifier.Brief{TaskID: "t2", DesignLabRequired: false}, t.TempDir())
	if lab.State != GateNotRequired {
		t.Fatalf("state = %s, want NOT_REQUIRED", lab.State)
	}
	if err := lab.GuardWrite("src/App.tsx"); err != nil {
		t.Errorf("unrequired gate blocked a write: %v", err)
	}
}

func TestGate_ApprovalUnlocksWrites(t *testing.T) {
	lab := pendingLab(t)
	if err := lab.Offer(twoDirections()); err != nil {
		t.Fatalf("Offer: %v", err)
	}
	if err := lab.GuardWrite("src/App.tsx"); err == nil {
		t.Fatal("offering directions must not by itself unlock writes")
	}
	if err := lab.Approve("b", "operator"); err != nil {
		t.Fatalf("Approve: %v", err)
	}
	if lab.State != GateApproved {
		t.Fatalf("state = %s, want APPROVED", lab.State)
	}
	if err := lab.GuardWrite("src/App.tsx"); err != nil {
		t.Errorf("approved gate still blocking: %v", err)
	}
	if lab.Approved == nil || lab.Approved.DirectionID != "b" || lab.Approved.ApprovedBy != "operator" {
		t.Errorf("approval not recorded: %+v", lab.Approved)
	}
	if _, err := filepath.Abs(lab.ApprovalPath()); err != nil {
		t.Errorf("approval path: %v", err)
	}
}

func TestGate_ApprovalRequiresAnOfferedDirectionAndANamedApprover(t *testing.T) {
	lab := pendingLab(t)
	_ = lab.Offer(twoDirections())

	if err := lab.Approve("does-not-exist", "operator"); err == nil {
		t.Error("approved a direction that was never offered")
	}
	if err := lab.Approve("a", "  "); err == nil {
		t.Error("approved with no named approver")
	}
	if lab.State != GatePending {
		t.Errorf("failed approvals changed state to %s", lab.State)
	}
}

// Two or three directions. One is a decree, four is a survey.
func TestGate_OfferRequiresTwoOrThreeSourcedDirections(t *testing.T) {
	dirs := twoDirections()

	if err := pendingLab(t).Offer(dirs[:1]); err == nil {
		t.Error("accepted a single direction")
	}
	if err := pendingLab(t).Offer(append(append([]Direction{}, dirs...), dirs[0], dirs[1])); err == nil {
		t.Error("accepted four directions")
	}

	unsourced := twoDirections()
	unsourced[0].TypographySrc = ""
	if err := pendingLab(t).Offer(unsourced); err == nil {
		t.Error("accepted a direction whose typography has no named source")
	}

	noMotionWhy := twoDirections()
	noMotionWhy[1].MotionWhy = ""
	if err := pendingLab(t).Offer(noMotionWhy); err == nil {
		t.Error("accepted a motion engine with no stated reason")
	}
}

// A rejection is only useful if it is remembered and it has a reason.
func TestGate_RejectionsArePersistedAndNotReOffered(t *testing.T) {
	workspace := t.TempDir()
	brief := &classifier.Brief{TaskID: "t3", DesignLabRequired: true, DesignLabReason: "PREMIUM"}

	first := NewDesignLab(brief, workspace)
	if err := first.Offer(twoDirections()); err != nil {
		t.Fatalf("Offer: %v", err)
	}
	if err := first.Reject("a", ""); err == nil {
		t.Error("accepted a rejection with no stated reason")
	}
	if err := first.Reject("a", "the umber reads as mud on my monitor"); err != nil {
		t.Fatalf("Reject: %v", err)
	}

	logged, err := first.LoadRejections()
	if err != nil {
		t.Fatalf("LoadRejections: %v", err)
	}
	if len(logged) != 1 {
		t.Fatalf("rejection log has %d entries, want 1", len(logged))
	}
	if logged[0].Reason == "" || logged[0].Fingerprint == "" {
		t.Errorf("rejection is missing reason or fingerprint: %+v", logged[0])
	}

	// A later pass at the same gate must not re-offer the rejected combination,
	// even under a different name.
	second := NewDesignLab(&classifier.Brief{TaskID: "t4", DesignLabRequired: true}, workspace)
	renamed := twoDirections()
	renamed[0].ID = "a-again"
	renamed[0].Concept = "Warm editorial, take two"
	if err := second.Offer(renamed); err == nil {
		t.Error("re-offered a rejected stack combination under a new name")
	}

	// Changing the actual stack is allowed.
	changed := twoDirections()
	changed[0].ID = "c"
	changed[0].ColorWorld = "Bleached linen with ink"
	changed[0].ColorSrc = "Japanese stationery catalogues"
	if err := second.Offer(changed); err != nil {
		t.Errorf("blocked a genuinely different direction: %v", err)
	}
}

// Bypass is allowed. Silent bypass is not.
func TestGate_BypassIsRecorded(t *testing.T) {
	lab := pendingLab(t)
	if err := lab.Bypass(""); err == nil {
		t.Error("accepted a bypass with no note")
	}
	if lab.State != GatePending {
		t.Fatalf("failed bypass changed state to %s", lab.State)
	}
	if err := lab.Bypass("operator waived the lab for a one-off internal demo"); err != nil {
		t.Fatalf("Bypass: %v", err)
	}
	if lab.State != GateBypassed {
		t.Errorf("state = %s, want BYPASSED", lab.State)
	}
	if err := lab.GuardWrite("src/App.tsx"); err != nil {
		t.Errorf("bypassed gate still blocking: %v", err)
	}
	if lab.Approved == nil || !lab.Approved.Bypass || lab.Approved.BypassNote == "" {
		t.Errorf("bypass not recorded: %+v", lab.Approved)
	}
}

func TestFingerprint_IgnoresNamesAndStackOrder(t *testing.T) {
	a := twoDirections()[0]
	b := a
	b.ID = "different"
	b.Concept = "Different name entirely"
	b.Stack = []string{"tailwind", "astro"}
	if Fingerprint(a) != Fingerprint(b) {
		t.Error("renaming a direction changed its fingerprint")
	}

	c := a
	c.MotionEngine = "Motion One"
	if Fingerprint(a) == Fingerprint(c) {
		t.Error("swapping the motion engine did not change the fingerprint")
	}
}
