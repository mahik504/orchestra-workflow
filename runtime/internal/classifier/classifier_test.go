package classifier

import (
	"strings"
	"testing"
)

// Every capability row must be able to decline itself. A route with no skip
// condition is a default in disguise, and defaults are how everything we build
// ends up looking the same.
func TestGraph_EveryRowHasTriggerAndSkip(t *testing.T) {
	g := loadRealGraph(t)
	if len(g.Capabilities) == 0 {
		t.Fatal("graph has no capabilities")
	}
	for id, cap := range g.Capabilities {
		if len(cap.TriggerConditions) == 0 {
			t.Errorf("%s: no trigger_conditions", id)
		}
		if len(cap.SkipConditions) == 0 {
			t.Errorf("%s: no skip_conditions", id)
		}
		switch cap.QualityBar {
		case BarStandard, BarPremium, BarExperimental:
		default:
			t.Errorf("%s: quality_bar %q is not one of STANDARD/PREMIUM/EXPERIMENTAL", id, cap.QualityBar)
		}
		if cap.RiskRank < 1 || cap.RiskRank > 10 {
			t.Errorf("%s: risk_rank %d out of range", id, cap.RiskRank)
		}
	}
}

// The classifier must report on every row it saw, not just the winner, so that
// "why not that route?" always has an answer on the record.
func TestClassify_ConsidersEveryRow(t *testing.T) {
	g := loadRealGraph(t)
	b := NewClassifierWithGraph(g).ClassifyBrief("Build a landing page for a coffee roastery", Options{})
	if len(b.Considered) != len(g.Capabilities) {
		t.Fatalf("considered %d rows, graph has %d", len(b.Considered), len(g.Capabilities))
	}
	for _, c := range b.Considered {
		if c.Declined && c.DeclineReason == "" {
			t.Errorf("%s was declined with no reason", c.CapabilityID)
		}
	}
}

func TestClassify_RoutesAndBars(t *testing.T) {
	g := loadRealGraph(t)
	c := NewClassifierWithGraph(g)

	cases := []struct {
		name     string
		request  string
		wantCap  string
		wantType string
		wantBar  string
		wantLab  bool
	}{
		{
			name:    "public brand page is premium and owes a lab",
			request: "Build a landing page for a friend's coffee roastery, they want it to feel expensive",
			wantCap: "premium-website", wantType: TypeFeature, wantBar: BarPremium, wantLab: true,
		},
		{
			name:    "backend bug is standard and owes nothing",
			request: "The login form throws a 500 when the email has a plus sign",
			wantCap: "", wantType: TypeBugfix, wantBar: BarStandard, wantLab: false,
		},
		{
			name:    "reading surface routes to the reader",
			request: "A reading app for arXiv papers with proper footnotes and math",
			wantCap: "academic-reader", wantType: TypeFeature, wantBar: BarPremium, wantLab: true,
		},
		{
			name:    "charts route to the dashboard, not the portal",
			request: "Build a scheduling dashboard for a school with attendance charts",
			wantCap: "saas-dashboard", wantType: TypeFeature, wantBar: BarStandard, wantLab: false,
		},
		{
			name:    "injection routes to the security audit and is not visual work",
			request: "Audit our API for injection and check the dependencies",
			wantCap: "security-audit", wantType: TypeSecurity, wantBar: BarStandard, wantLab: false,
		},
		{
			name:    "a named URL routes to extraction",
			request: "Extract the design tokens from https://example.com and write them up",
			wantCap: "reverse-engineering", wantType: TypeDesign, wantBar: BarStandard, wantLab: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			b := c.ClassifyBrief(tc.request, Options{})
			if b.CapabilityID != tc.wantCap {
				t.Errorf("capability = %q, want %q (reason: %s)", b.CapabilityID, tc.wantCap, b.ArchetypeReason)
			}
			if b.Type != tc.wantType {
				t.Errorf("type = %q, want %q", b.Type, tc.wantType)
			}
			if b.QualityBar != tc.wantBar {
				t.Errorf("quality bar = %q, want %q (%s)", b.QualityBar, tc.wantBar, b.QualityBarReason)
			}
			if b.DesignLabRequired != tc.wantLab {
				t.Errorf("design lab = %v, want %v (%s)", b.DesignLabRequired, tc.wantLab, b.DesignLabReason)
			}
		})
	}
}

// The shop-portfolio is the brief that used to get silently forced into one
// archetype. It should produce exactly one question, and on silence should fall
// to the cheaper of the two routes.
func TestClassify_AmbiguousShopPortfolio(t *testing.T) {
	g := loadRealGraph(t)
	b := NewClassifierWithGraph(g).ClassifyBrief(
		"I want a portfolio site that also sells prints, with a checkout and an admin area to manage orders",
		Options{})

	if !b.Ambiguous {
		t.Fatalf("expected an ambiguous brief, got %s alone (%s)", b.CapabilityID, b.ArchetypeReason)
	}
	if b.ClarifyingQuestion == "" {
		t.Fatal("ambiguous brief asked no question")
	}
	if strings.Count(b.ClarifyingQuestion, "?") != 1 {
		t.Errorf("expected exactly one question, got: %s", b.ClarifyingQuestion)
	}
	if len(b.Selected) < 2 {
		t.Fatalf("expected two contenders, got %d", len(b.Selected))
	}

	first, second := b.Selected[0], b.Selected[1]
	b.ResolveSilence()

	if b.Ambiguous {
		t.Error("brief is still marked ambiguous after resolving silence")
	}
	if !b.Assumed {
		t.Error("silence fallback did not mark the brief as assumed")
	}
	want := "assumed " + b.CapabilityID + ", no response"
	if b.AssumptionNote != want {
		t.Errorf("assumption note = %q, want %q", b.AssumptionNote, want)
	}

	cheaper := first
	if second.RiskRank < first.RiskRank {
		cheaper = second
	}
	if b.CapabilityID != cheaper.CapabilityID {
		t.Errorf("silence picked %s (risk %d); lower-risk route was %s (risk %d)",
			b.CapabilityID, b.Selected[0].RiskRank, cheaper.CapabilityID, cheaper.RiskRank)
	}
}

// A skip condition has to actually cost a route the job, otherwise it is decoration.
func TestClassify_SkipConditionDeclinesRoute(t *testing.T) {
	g := loadRealGraph(t)
	b := NewClassifierWithGraph(g).ClassifyBrief(
		"An admin console behind a login for staff to manage tenants, roles and billing", Options{})

	if b.CapabilityID != "b2b-portal" {
		t.Fatalf("expected b2b-portal, got %q", b.CapabilityID)
	}
	for _, c := range b.Considered {
		if c.CapabilityID == "premium-website" {
			if !c.Declined {
				t.Errorf("premium-website should have been declined for a staff-only console, scored %.2f", c.Score)
			}
			if len(c.FiredSkips) == 0 {
				t.Error("premium-website was declined but no skip condition is on the record")
			}
			return
		}
	}
	t.Fatal("premium-website was never considered")
}

func TestClassify_SkipTheLabIsHonoured(t *testing.T) {
	g := loadRealGraph(t)
	c := NewClassifierWithGraph(g)

	req := "Build a landing page for a friend's coffee roastery, they want it to feel expensive"
	if b := c.ClassifyBrief(req, Options{}); !b.DesignLabRequired {
		t.Fatal("premium visual work should require the lab by default")
	}
	if b := c.ClassifyBrief(req+", skip the lab", Options{}); b.DesignLabRequired {
		t.Errorf("operator said skip the lab, gate still armed: %s", b.DesignLabReason)
	}
	if b := c.ClassifyBrief(req, Options{SkipLab: true}); b.DesignLabRequired {
		t.Error("SkipLab option ignored")
	}
}

// STANDARD is off by default; PREMIUM is on. That asymmetry is the whole policy.
func TestClassify_LabDefaultsFollowTheBar(t *testing.T) {
	g := loadRealGraph(t)
	c := NewClassifierWithGraph(g)

	internal := c.ClassifyBrief("Internal tool for our team to see which servers are down right now, live", Options{})
	if internal.QualityBar != BarStandard {
		t.Errorf("internal tool should run at STANDARD, got %s (%s)", internal.QualityBar, internal.QualityBarReason)
	}
	if internal.DesignLabRequired {
		t.Errorf("STANDARD work should not force a lab: %s", internal.DesignLabReason)
	}

	client := c.ClassifyBrief("A mission control cockpit with live telemetry for a client demo", Options{})
	if client.QualityBar == BarStandard {
		t.Errorf("client-facing work should not sit at STANDARD (%s)", client.QualityBarReason)
	}
	if !client.DesignLabRequired {
		t.Errorf("premium visual work should require the lab: %s", client.DesignLabReason)
	}
}

func TestClassify_DepthFollowsTheBar(t *testing.T) {
	g := loadRealGraph(t)
	c := NewClassifierWithGraph(g)

	visual := c.ClassifyBrief("Build a landing page for a client's studio", Options{})
	if visual.ResearchDepth != ResearchDeep {
		t.Errorf("research depth = %s, want DEEP", visual.ResearchDepth)
	}
	if visual.VerifyDepth != VerifyBrowserViewports {
		t.Errorf("verify depth = %s, want %s", visual.VerifyDepth, VerifyBrowserViewports)
	}

	backend := c.ClassifyBrief("The login form throws a 500 when the email has a plus sign", Options{})
	if backend.ResearchDepth != ResearchNone {
		t.Errorf("research depth = %s, want NONE", backend.ResearchDepth)
	}
	if backend.VerifyDepth != VerifyBuild {
		t.Errorf("verify depth = %s, want BUILD", backend.VerifyDepth)
	}
}

func TestClassify_NoGraphStillProducesUsableBrief(t *testing.T) {
	b := NewClassifier().ClassifyBrief("Build something", Options{})
	if b.Archetype != "standard-feature" {
		t.Errorf("archetype = %q, want standard-feature", b.Archetype)
	}
	if b.QualityBar != BarStandard {
		t.Errorf("quality bar = %q, want STANDARD", b.QualityBar)
	}
	if b.DesignLabRequired {
		t.Error("no graph should never arm the gate")
	}
}

func TestClassify_ExtractsHardConstraints(t *testing.T) {
	g := loadRealGraph(t)
	b := NewClassifierWithGraph(g).ClassifyBrief(
		"Build the landing page. It must work offline. Do not add any paid services.", Options{})
	if len(b.HardConstraints) < 2 {
		t.Errorf("expected the must/do-not sentences to be captured, got %v", b.HardConstraints)
	}
}
