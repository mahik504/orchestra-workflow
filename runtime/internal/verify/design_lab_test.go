package verify

import (
	"errors"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/mahik504/orchestra-workflow/runtime/internal/classifier"
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

func oneDirection() Direction {
	return Direction{
		ID: "a", Concept: "Warm editorial",
		Typography: "Freight Display + Söhne", TypographySrc: "Klim type specimen",
		ColorWorld: "Roasted umber on bone", ColorSrc: "1970s coffee packaging archive",
		LayoutLanguage: "Asymmetric editorial grid", ComponentKit: "custom",
		MotionEngine: "CSS transitions", MotionWhy: "page is mostly static, no timeline needed",
		LogoMethod: "wordmark", IconSystem: "Phosphor", Stack: []string{"astro", "tailwind"},
		ThreeD: "no", Shader: "no",
	}
}

func altDirection() Direction {
	d := oneDirection()
	d.ID = "b"
	d.Concept = "Stark industrial"
	d.Typography = "Diatype + Diatype Mono"
	d.TypographySrc = "Dinamo specimen"
	d.ColorWorld = "Cold steel with ember accent"
	d.ColorSrc = "Braun product photography"
	d.LayoutLanguage = "Strict 12-column"
	d.MotionEngine = "GSAP"
	d.MotionWhy = "scroll-linked reveals need a timeline"
	d.LogoMethod = "monogram"
	d.IconSystem = "Lucide"
	d.Stack = []string{"next", "tailwind"}
	return d
}

func twentyThreeCards() []DirectionCard {
	cards := make([]DirectionCard, SurveyCardCount)
	for i := 0; i < SurveyCardCount; i++ {
		cards[i] = DirectionCard{
			ID:           fmt.Sprintf("c%d", i+1),
			Name:         fmt.Sprintf("Direction %d", i+1),
			OneLiner:     "Short survey card, not a full DESIGN.md",
			Typography:   "Named pairing",
			ColorWorld:   "Named palette",
			ThreeD:       "no",
			MotionEngine: "CSS",
		}
	}
	cards[0].MotionEngine = "GSAP"
	cards[1].ThreeD = "yes — R3F"
	return cards
}

func readyContract(t *testing.T, dir Direction) *DesignLab {
	t.Helper()
	lab := pendingLab(t)
	if err := lab.SkipSurvey("test fixture: skip survey to isolate the contract"); err != nil {
		t.Fatalf("SkipSurvey: %v", err)
	}
	if err := lab.OfferContract(dir); err != nil {
		t.Fatalf("OfferContract: %v", err)
	}
	return lab
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
	lab := readyContract(t, altDirection())
	if err := lab.GuardWrite("src/App.tsx"); err == nil {
		t.Fatal("offering a contract must not by itself unlock writes")
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
	lab := readyContract(t, oneDirection())

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

func TestGate_SurveyRequiresExactlyTwentyThreeCards(t *testing.T) {
	short := twentyThreeCards()[:10]
	if err := pendingLab(t).OfferSurvey(short); err == nil {
		t.Error("accepted a short survey")
	}
	long := append(twentyThreeCards(), twentyThreeCards()[0])
	if err := pendingLab(t).OfferSurvey(long); err == nil {
		t.Error("accepted more than 23 cards")
	}
	incomplete := twentyThreeCards()
	incomplete[3].OneLiner = ""
	if err := pendingLab(t).OfferSurvey(incomplete); err == nil {
		t.Error("accepted an incomplete survey card")
	}
	ok := pendingLab(t)
	if err := ok.OfferSurvey(twentyThreeCards()); err != nil {
		t.Fatalf("valid survey refused: %v", err)
	}
	if err := ok.GuardWrite("src/App.tsx"); err == nil {
		t.Fatal("a survey must not unlock frontend writes")
	}
}

func TestGate_ContractRequiresSurveyOrSkipAndSources(t *testing.T) {
	lab := pendingLab(t)
	if err := lab.OfferContract(oneDirection()); err == nil {
		t.Error("accepted a contract with no survey and no skip")
	}

	lab = pendingLab(t)
	if err := lab.SkipSurvey(""); err == nil {
		t.Error("accepted an empty skip reason")
	}
	if err := lab.SkipSurvey("prompt named DESIGN.md"); err != nil {
		t.Fatalf("SkipSurvey: %v", err)
	}

	unsourced := oneDirection()
	unsourced.TypographySrc = ""
	if err := lab.OfferContract(unsourced); err == nil {
		t.Error("accepted a contract whose typography has no named source")
	}

	noMotionWhy := oneDirection()
	noMotionWhy.MotionWhy = ""
	if err := lab.OfferContract(noMotionWhy); err == nil {
		t.Error("accepted a motion engine with no stated reason")
	}

	if err := lab.OfferContract(oneDirection()); err != nil {
		t.Fatalf("valid contract after skip refused: %v", err)
	}
}

func TestGate_SurveyThenOneContract(t *testing.T) {
	lab := pendingLab(t)
	if err := lab.OfferSurvey(twentyThreeCards()); err != nil {
		t.Fatalf("OfferSurvey: %v", err)
	}
	if err := lab.OfferContract(oneDirection()); err != nil {
		t.Fatalf("OfferContract after survey: %v", err)
	}
	if err := lab.Approve("a", "operator"); err != nil {
		t.Fatalf("Approve: %v", err)
	}
}

func TestGate_ApproveCustomSkipsSurvey(t *testing.T) {
	lab := pendingLab(t)
	if err := lab.ApproveCustom("operator", ""); err == nil {
		t.Error("accepted custom approval with no note")
	}
	if err := lab.ApproveCustom("operator", "pasted DESIGN.md for the tinted editorial system"); err != nil {
		t.Fatalf("ApproveCustom: %v", err)
	}
	if lab.State != GateApproved {
		t.Fatalf("state = %s, want APPROVED", lab.State)
	}
	if err := lab.GuardWrite("src/App.tsx"); err != nil {
		t.Errorf("custom DESIGN.md still blocking: %v", err)
	}
	if lab.Approved == nil || lab.Approved.DirectionID != "custom" {
		t.Errorf("custom approval not recorded: %+v", lab.Approved)
	}
}

func TestGate_NamedOverrideSkipsSurvey(t *testing.T) {
	lab := pendingLab(t)
	if err := lab.SkipSurvey("prompt named refero-design and a pasted DESIGN.md"); err != nil {
		t.Fatalf("SkipSurvey: %v", err)
	}
	if err := lab.OfferSurvey(twentyThreeCards()); err == nil {
		t.Error("offered a survey after a named override")
	}
	if err := lab.OfferContract(oneDirection()); err != nil {
		t.Fatalf("OfferContract after named override: %v", err)
	}
}

func TestGate_RejectionsArePersistedAndNotReOffered(t *testing.T) {
	workspace := t.TempDir()
	brief := &classifier.Brief{TaskID: "t3", DesignLabRequired: true, DesignLabReason: "PREMIUM"}

	first := NewDesignLab(brief, workspace)
	if err := first.SkipSurvey("test fixture"); err != nil {
		t.Fatalf("SkipSurvey: %v", err)
	}
	if err := first.OfferContract(oneDirection()); err != nil {
		t.Fatalf("OfferContract: %v", err)
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

	second := NewDesignLab(&classifier.Brief{TaskID: "t4", DesignLabRequired: true}, workspace)
	if err := second.SkipSurvey("second pass"); err != nil {
		t.Fatalf("SkipSurvey: %v", err)
	}
	renamed := oneDirection()
	renamed.ID = "a-again"
	renamed.Concept = "Warm editorial, take two"
	if err := second.OfferContract(renamed); err == nil {
		t.Error("re-offered a rejected stack combination under a new name")
	}

	changed := oneDirection()
	changed.ID = "c"
	changed.ColorWorld = "Bleached linen with ink"
	changed.ColorSrc = "Japanese stationery catalogues"
	if err := second.OfferContract(changed); err != nil {
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
	a := oneDirection()
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
