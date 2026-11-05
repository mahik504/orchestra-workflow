package capability

import (
	"errors"
	"strings"

	"github.com/mahik504/orchestra-workflow/runtime/internal/mission"
	"github.com/mahik504/orchestra-workflow/runtime/internal/visual"
)

const (
	IDCodingWorker      = "coding.worker"
	IDResearchWorker    = "research.worker"
	IDDesignWorker      = "design.worker"
	IDBrowserPlaywright = "browser.playwright"
	IDVisualStructural  = "visual.structural"
	IDVisualScreenshot  = "visual.screenshot"
	IDVisualMultimodal  = "visual.multimodal"
	IDJudgeIndependent  = "judge.independent"

	CatCodingWorker      = "coding.worker"
	CatResearchWorker    = "research.worker"
	CatDesignWorker      = "design.worker"
	CatBrowserPlaywright = "browser.playwright"
	CatVisualStructural  = "visual.structural"
	CatVisualScreenshot  = "visual.screenshot"
	CatVisualMultimodal  = "visual.multimodal"
	CatJudgeIndependent  = "judge.independent"
	CatSkill             = "skill"
	CatResource          = "resource"
	CatMCPTool           = "mcp.tool"

	StrengthRequired      = "required"
	StrengthPreferred     = "preferred"
	StrengthOptional      = "optional"
	StrengthFallback      = "fallback"
	StrengthUnavailable   = "unavailable"
	StrengthHumanRequired = "human_required"

	MaxSkills    = 8
	MaxResources = 8
	MaxMCPTools  = 4
)

var ErrHumanStop = errors.New("human approval gate required before writing code")

func IsHumanStop(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, ErrHumanStop) || strings.Contains(strings.ToLower(err.Error()), "human approval gate required")
}

type HostSnapshot struct {
	ID        string `json:"id,omitempty"`
	Available bool   `json:"available,omitempty"`
}

type MCPCandidate struct {
	ID       string `json:"id,omitempty"`
	Health   string `json:"health,omitempty"`
	Required bool   `json:"required,omitempty"`
	Fallback string `json:"fallback,omitempty"`
}

type Input struct {
	Mission                 *mission.Mission
	Class                   visual.Class
	Skills                  []string
	Resources               []string
	MCP                     []MCPCandidate
	Hosts                   []HostSnapshot
	PlaywrightAvailable     *bool
	ScreenshotEvalAvailable *bool
	DiscoveredModels        []string
	Judge                   *visual.JudgeInput
	// Trusted IDs were named by the winning route. They skip the keyword filter.
	Trusted []string
}

type Item struct {
	ID        string `json:"id,omitempty"`
	Category  string `json:"category,omitempty"`
	Strength  string `json:"strength,omitempty"`
	Available bool   `json:"available,omitempty"`
	Reason    string `json:"reason,omitempty"`
}

type Fallback struct {
	Requested string `json:"requested,omitempty"`
	Selected  string `json:"selected,omitempty"`
	Reason    string `json:"reason,omitempty"`
}

type Decision struct {
	MissionID   string     `json:"mission_id,omitempty"`
	Selected    []Item     `json:"selected,omitempty"`
	Fallbacks   []Fallback `json:"fallbacks,omitempty"`
	Rejected    []Item     `json:"rejected,omitempty"`
	Unavailable []Item     `json:"unavailable,omitempty"`
	HostLabel   string     `json:"host_label,omitempty"`
	Skills      []string   `json:"skills,omitempty"`
	Resources   []string   `json:"resources,omitempty"`
	Tools       []string   `json:"tools,omitempty"`
	HumanDetail string     `json:"-"`
}

func (d Decision) IDs(category string) []string {
	var out []string
	for _, it := range d.Selected {
		if it.Category == category && it.ID != "" {
			out = append(out, it.ID)
		}
	}
	return out
}

func (d Decision) Has(id string) bool {
	for _, it := range d.Selected {
		if it.ID == id {
			return true
		}
	}
	return false
}

func Select(in Input) (Decision, error) {
	_ = in.DiscoveredModels

	d := Decision{}
	if in.Mission != nil {
		d.MissionID = in.Mission.ID
	}

	class := in.Class
	if class == "" && in.Mission != nil {
		class = visual.Classify(visual.ClassifyInput{
			Type:             string(in.Mission.Type),
			Objective:        in.Mission.Objective,
			Tags:             in.Mission.Capabilities,
			RequiresVisual:   in.Mission.VisualImportance != mission.VisualNone && in.Mission.VisualImportance != "",
			VisualImportance: string(in.Mission.VisualImportance),
		})
	}
	policy := visual.PolicyFor(class)

	var budget mission.Budget
	if in.Mission != nil {
		budget = in.Mission.Budget
	}
	review := string(budget.VisualReview)
	if review == "" {
		review = policy.VisualReview
	}

	d.HostLabel = pickHostLabel(in.Hosts)

	addSelected(&d, Item{
		ID:        IDCodingWorker,
		Category:  CatCodingWorker,
		Strength:  StrengthRequired,
		Available: true,
		Reason:    "implementation compiles a builder capsule",
	})

	if budget.Research == mission.ResearchTargeted || budget.Research == mission.ResearchDeep {
		addSelected(&d, Item{
			ID:        IDResearchWorker,
			Category:  CatResearchWorker,
			Strength:  StrengthRequired,
			Available: true,
			Reason:    "budget.research is " + string(budget.Research),
		})
	}

	if policy.Directions {
		addSelected(&d, Item{
			ID:        IDDesignWorker,
			Category:  CatDesignWorker,
			Strength:  StrengthRequired,
			Available: true,
			Reason:    "flagship policy requires one sealed design contract",
		})
	}

	browserWanted := review == string(mission.VisualReviewBrowser) || review == string(mission.VisualReviewDeep)
	pwAvail := playwrightAvail(in)
	if browserWanted {
		if !pwAvail {
			d.HumanDetail = "required: browser.playwright; why: visual_review=" + review +
				"; unavailable: harness Chromium Layer A; fallback prohibited"
			addUnavailable(&d, Item{
				ID:        IDBrowserPlaywright,
				Category:  CatBrowserPlaywright,
				Strength:  StrengthHumanRequired,
				Available: false,
				Reason:    d.HumanDetail,
			})
			return d, ErrHumanStop
		}
		addSelected(&d, Item{
			ID:        IDBrowserPlaywright,
			Category:  CatBrowserPlaywright,
			Strength:  StrengthRequired,
			Available: true,
			Reason:    "visual_review=" + review + " requires harness Playwright",
		})
	}

	shotAvail := screenshotAvail(in)
	if browserWanted && class != visual.ClassOrdinary {
		if class == visual.ClassFlagship && review == string(mission.VisualReviewDeep) {
			addUnavailable(&d, Item{
				ID:        IDVisualMultimodal,
				Category:  CatVisualMultimodal,
				Strength:  StrengthPreferred,
				Available: false,
				Reason:    "no multimodal adapter is registered",
			})
			if shotAvail {
				addSelected(&d, Item{
					ID:        IDVisualScreenshot,
					Category:  CatVisualScreenshot,
					Strength:  StrengthFallback,
					Available: true,
					Reason:    "screenshot evaluator is registered",
				})
				d.Fallbacks = append(d.Fallbacks, Fallback{
					Requested: IDVisualMultimodal,
					Selected:  IDVisualScreenshot,
					Reason:    "multimodal unavailable; screenshot adapter present",
				})
			} else {
				addUnavailable(&d, Item{
					ID:        IDVisualScreenshot,
					Category:  CatVisualScreenshot,
					Strength:  StrengthPreferred,
					Available: false,
					Reason:    "no screenshot adapter registered",
				})
				addSelected(&d, Item{
					ID:        IDVisualStructural,
					Category:  CatVisualStructural,
					Strength:  StrengthFallback,
					Available: true,
					Reason:    "structural critic; not a vision pass",
				})
				d.Fallbacks = append(d.Fallbacks, Fallback{
					Requested: IDVisualMultimodal,
					Selected:  IDVisualStructural,
					Reason:    "multimodal unavailable; screenshot adapter absent; structural critic runs",
				})
				d.Fallbacks = append(d.Fallbacks, Fallback{
					Requested: IDVisualScreenshot,
					Selected:  IDVisualStructural,
					Reason:    "screenshot adapter absent; structural critic runs",
				})
			}
		}
		if shotAvail && !d.Has(IDVisualScreenshot) {
			addSelected(&d, Item{
				ID:        IDVisualScreenshot,
				Category:  CatVisualScreenshot,
				Strength:  StrengthPreferred,
				Available: true,
				Reason:    "screenshot evaluator is registered",
			})
		} else if !shotAvail && !d.Has(IDVisualStructural) {
			addUnavailable(&d, Item{
				ID:        IDVisualScreenshot,
				Category:  CatVisualScreenshot,
				Strength:  StrengthPreferred,
				Available: false,
				Reason:    "no screenshot adapter registered",
			})
			d.Fallbacks = append(d.Fallbacks, Fallback{
				Requested: IDVisualScreenshot,
				Selected:  IDVisualStructural,
				Reason:    "screenshot adapter absent; structural critic runs",
			})
		}
		if !d.Has(IDVisualStructural) {
			addSelected(&d, Item{
				ID:        IDVisualStructural,
				Category:  CatVisualStructural,
				Strength:  StrengthRequired,
				Available: true,
				Reason:    "Layer B structural critic on Layer A captures",
			})
		}
	}

	jd := decideJudge(in, class, budget)
	if jd.Activated {
		addSelected(&d, Item{
			ID:        IDJudgeIndependent,
			Category:  CatJudgeIndependent,
			Strength:  StrengthRequired,
			Available: true,
			Reason:    jd.Reason,
		})
	}

	cats := selectedCategories(d)
	trust := trustedSet(in.Trusted)
	d.Skills = filterIDs(in.Skills, cats, MaxSkills, CatSkill, &d, trust)
	d.Resources = filterIDs(in.Resources, cats, MaxResources, CatResource, &d, trust)
	tools, err := filterMCP(in.MCP, cats, &d, trust)
	if err != nil {
		return d, err
	}
	d.Tools = tools
	return d, nil
}

func playwrightAvail(in Input) bool {
	if in.PlaywrightAvailable != nil {
		return *in.PlaywrightAvailable
	}
	return visual.ProbePlaywright().Available
}

func screenshotAvail(in Input) bool {
	if in.ScreenshotEvalAvailable != nil {
		return *in.ScreenshotEvalAvailable
	}
	return visual.ActiveScreenshotEvaluator() != nil
}

func decideJudge(in Input, class visual.Class, budget mission.Budget) visual.JudgeDecision {
	ji := visual.JudgeInput{Class: class, Policy: string(budget.IndependentJudge)}
	if in.Judge != nil {
		ji = *in.Judge
		if ji.Class == "" {
			ji.Class = class
		}
	} else if in.Mission != nil {
		ji.VisualNovelty = in.Mission.Novelty
		ji.Confidence = in.Mission.Confidence
		ji.MultiSurface = len(in.Mission.AffectedSurfaces) > 1
		ji.SecurityRelease = in.Mission.RequiresSecurity
		ji.HumanRequested = budget.IndependentJudge == mission.JudgeRequired
	}
	return visual.DecideJudge(ji)
}

func pickHostLabel(hosts []HostSnapshot) string {
	if len(hosts) == 0 {
		hosts = []HostSnapshot{{ID: "cursor", Available: true}, {ID: "claude", Available: true}}
	}
	var cursorOK bool
	var other string
	for _, h := range hosts {
		id := strings.TrimSpace(h.ID)
		if id == "" || !h.Available {
			continue
		}
		if id == "cursor" {
			cursorOK = true
		} else if other == "" {
			other = id
		}
	}
	if cursorOK {
		return "cursor"
	}
	return other
}

func selectedCategories(d Decision) map[string]bool {
	out := map[string]bool{}
	for _, it := range d.Selected {
		out[it.Category] = true
		out[it.ID] = true
	}
	return out
}

func trustedSet(ids []string) map[string]bool {
	out := map[string]bool{}
	for _, raw := range ids {
		id := strings.TrimSpace(raw)
		if id != "" {
			out[id] = true
		}
	}
	return out
}

func filterIDs(ids []string, cats map[string]bool, capN int, category string, d *Decision, trust map[string]bool) []string {
	seen := map[string]bool{}
	var out []string
	for _, raw := range ids {
		id := strings.TrimSpace(raw)
		if id == "" || seen[id] {
			continue
		}
		seen[id] = true
		if !trust[id] && !idRelevant(id, cats) {
			addRejected(d, Item{
				ID:       id,
				Category: category,
				Strength: StrengthOptional,
				Reason:   "not relevant to selected categories",
			})
			continue
		}
		if len(out) >= capN {
			addRejected(d, Item{
				ID:       id,
				Category: category,
				Strength: StrengthOptional,
				Reason:   "over list cap",
			})
			continue
		}
		out = append(out, id)
		addSelected(d, Item{
			ID:        id,
			Category:  category,
			Strength:  StrengthOptional,
			Available: true,
			Reason:    "narrowed candidate",
		})
	}
	return out
}

func filterMCP(cands []MCPCandidate, cats map[string]bool, d *Decision, trust map[string]bool) ([]string, error) {
	seen := map[string]bool{}
	var out []string
	for _, c := range cands {
		id := strings.TrimSpace(c.ID)
		if id == "" || seen[id] {
			continue
		}
		seen[id] = true
		health := strings.ToUpper(strings.TrimSpace(c.Health))
		// Health does not decide whether the mission requires the tool.
		// Only HEALTHY satisfies availability. OPTIONAL and blank are unprobed.
		if !mcpSatisfies(health) {
			item := Item{
				ID:        id,
				Category:  CatMCPTool,
				Available: false,
				Reason:    "mcp health " + health,
			}
			if c.Required {
				item.Strength = StrengthHumanRequired
				d.HumanDetail = "required: mcp.tool " + id + "; why: required tool; unavailable: health=" + health + "; fallback prohibited"
				d.Unavailable = append(d.Unavailable, item)
				return out, ErrHumanStop
			}
			item.Strength = StrengthOptional
			d.Unavailable = append(d.Unavailable, item)
			reason := "optional MCP skipped: " + health
			if c.Fallback != "" {
				reason = "optional MCP skipped: " + health + "; fallback: " + c.Fallback
				d.Fallbacks = append(d.Fallbacks, Fallback{Requested: id, Reason: c.Fallback})
			}
			addRejected(d, Item{
				ID: id, Category: CatMCPTool, Strength: StrengthOptional,
				Reason: reason,
			})
			continue
		}
		if !trust[id] && !idRelevant(id, cats) {
			addRejected(d, Item{
				ID: id, Category: CatMCPTool, Strength: StrengthOptional,
				Reason: "not relevant to selected categories",
			})
			continue
		}
		if len(out) >= MaxMCPTools {
			addRejected(d, Item{
				ID: id, Category: CatMCPTool, Strength: StrengthOptional,
				Reason: "over list cap",
			})
			continue
		}
		out = append(out, id)
		addSelected(d, Item{
			ID: id, Category: CatMCPTool, Strength: StrengthOptional, Available: true,
			Reason: "narrowed candidate",
		})
	}
	return out, nil
}

// mcpSatisfies is availability, not mission requirement.
// HEALTHY is the only probed-good state. OPTIONAL, blank, BROKEN,
// AUTH_REQUIRED, DISABLED, and PURCHASE_PENDING do not satisfy a tool.
func mcpSatisfies(health string) bool {
	return health == "HEALTHY"
}

func idRelevant(id string, cats map[string]bool) bool {
	low := strings.ToLower(id)
	visualish := containsAny(low, []string{
		"gsap", "impeccable", "taste", "hallmark", "visual", "animate",
		"stitch", "golden", "design-lab", "refero", "mobbin", "screenshot",
		"playwright", "browser",
	})
	researchish := containsAny(low, []string{"research", "scrapling", "context7", "web-search"})
	switch {
	case visualish:
		return cats[CatDesignWorker] || cats[CatVisualStructural] || cats[CatVisualScreenshot] ||
			cats[CatVisualMultimodal] || cats[CatBrowserPlaywright]
	case researchish:
		return cats[CatResearchWorker]
	default:
		return cats[CatCodingWorker]
	}
}

func containsAny(low string, parts []string) bool {
	for _, p := range parts {
		if strings.Contains(low, p) {
			return true
		}
	}
	return false
}

func addSelected(d *Decision, it Item) {
	if d.Has(it.ID) {
		return
	}
	d.Selected = append(d.Selected, it)
}

const maxRejectedRows = 16

func addRejected(d *Decision, it Item) {
	if d == nil || len(d.Rejected) >= maxRejectedRows {
		return
	}
	d.Rejected = append(d.Rejected, it)
}

func addUnavailable(d *Decision, it Item) {
	d.Unavailable = append(d.Unavailable, it)
}
