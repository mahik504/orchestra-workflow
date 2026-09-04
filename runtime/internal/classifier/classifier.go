// Package classifier turns a raw human request into a structured brief.
//
// The brief is the thing the human approves at the re-brief step: archetype,
// quality bar, platform, and the hard constraints we think we heard. Getting it
// wrong here is cheap to fix and expensive to discover later, which is why the
// classifier reports its runners-up and asks one question when two routes fit.
package classifier

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/user/orchestra-v3/internal/resources"
)

// Task type constants.
const (
	TypeFeature  = "FEATURE"
	TypeBugfix   = "BUGFIX"
	TypeRefactor = "REFACTOR"
	TypeDesign   = "DESIGN"
	TypeResearch = "RESEARCH"
	TypeSecurity = "SECURITY"
)

// Quality bars, mirroring the contract in AGENTS.md.
const (
	BarStandard     = "STANDARD"
	BarPremium      = "PREMIUM"
	BarExperimental = "EXPERIMENTAL"
)

// Research and verification depth.
const (
	ResearchNone  = "NONE"
	ResearchLight = "LIGHT"
	ResearchDeep  = "DEEP"

	VerifyBuild            = "BUILD"
	VerifyBrowser          = "BROWSER"
	VerifyBrowserViewports = "BROWSER_MULTIVIEWPORT"
)

type UserOverride struct {
	ForceAgent      string `json:"force_agent"`
	SkipVisualGate  bool   `json:"skip_visual_gate"`
	ForceBypassGate bool   `json:"force_bypass_gate"`
}

// Task is the normalized request shape consumed by the router and research
// coordinator. It is a projection of Brief.
type Task struct {
	ID                 string        `json:"id"`
	RawRequest         string        `json:"raw_request"`
	Type               string        `json:"type"`
	RequiresVisual     bool          `json:"requires_visual"`
	RequiresSecurity   bool          `json:"requires_security"`
	ExtractedKeywords  []string      `json:"extracted_keywords"`
	SuggestedResources []string      `json:"suggested_resources"`
	UserOverride       *UserOverride `json:"user_override,omitempty"`
}

// Candidate is one capability row weighed against the request. Every row in the
// graph produces a Candidate, including the ones we decline, so that "why not
// that route?" always has an answer.
type Candidate struct {
	CapabilityID     string   `json:"capability_id"`
	Name             string   `json:"name"`
	PrimaryArchetype string   `json:"primary_archetype"`
	QualityBar       string   `json:"quality_bar"`
	RiskRank         int      `json:"risk_rank"`
	Platform         string   `json:"platform"`
	Score            float64  `json:"score"`
	MatchedTags      []string `json:"matched_tags,omitempty"`
	FiredTriggers    []string `json:"fired_triggers,omitempty"`
	FiredSkips       []string `json:"fired_skips,omitempty"`
	Declined         bool     `json:"declined"`
	DeclineReason    string   `json:"decline_reason,omitempty"`
}

// Brief is the structured re-brief. Ambiguous/ClarifyingQuestion carry the one
// question we are allowed to ask; Assumed records what we picked when nobody
// answered it.
type Brief struct {
	TaskID             string      `json:"task_id"`
	RawRequest         string      `json:"raw_request"`
	Type               string      `json:"type"`
	Archetype          string      `json:"archetype"`
	CapabilityID       string      `json:"capability_id"`
	ArchetypeReason    string      `json:"archetype_reason"`
	QualityBar         string      `json:"quality_bar"`
	QualityBarReason   string      `json:"quality_bar_reason"`
	Platform           string      `json:"platform"`
	RequiresVisual     bool        `json:"requires_visual"`
	RequiresSecurity   bool        `json:"requires_security"`
	ResearchDepth      string      `json:"research_depth"`
	VerifyDepth        string      `json:"verify_depth"`
	DesignLabRequired  bool        `json:"design_lab_required"`
	DesignLabReason    string      `json:"design_lab_reason"`
	HardConstraints    []string    `json:"hard_constraints,omitempty"`
	Selected           []Candidate `json:"selected"`
	Considered         []Candidate `json:"considered"`
	Ambiguous          bool        `json:"ambiguous"`
	ClarifyingQuestion string      `json:"clarifying_question,omitempty"`
	Assumed            bool        `json:"assumed"`
	AssumptionNote     string      `json:"assumption_note,omitempty"`
	UnknownTechnology  []string    `json:"unknown_technology,omitempty"`
}

// Options carries signals the caller already knows so the classifier does not
// have to guess them from prose.
type Options struct {
	TaskID        string
	ExtraTags     []string
	ArchetypeHint string
	SkipLab       bool
	DeclaredType  string
}

type Classifier struct {
	graph *resources.DesignResourceGraph
}

func NewClassifier() *Classifier { return &Classifier{} }

// NewClassifierWithGraph builds a classifier that can route against the
// capability graph. Without a graph the classifier can still read type, visual,
// and security intent, but it cannot choose an archetype.
func NewClassifierWithGraph(g *resources.DesignResourceGraph) *Classifier {
	return &Classifier{graph: g}
}

// ---------- lexical plumbing ----------

var wordRE = regexp.MustCompile(`[a-z0-9]+`)

var stopwords = map[string]bool{
	"the": true, "and": true, "for": true, "that": true, "with": true, "this": true,
	"than": true, "then": true, "them": true, "they": true, "there": true, "their": true,
	"from": true, "into": true, "onto": true, "over": true, "only": true, "also": true,
	"rather": true, "instead": true, "route": true, "should": true, "would": true,
	"which": true, "when": true, "what": true, "where": true, "while": true, "does": true,
	"have": true, "having": true, "been": true, "being": true, "will": true, "want": true,
	"wants": true, "need": true, "needs": true, "make": true, "made": true, "like": true,
	"just": true, "some": true, "same": true, "such": true, "each": true, "more": true,
	"most": true, "much": true, "many": true, "very": true, "real": true, "even": true,
	"here": true, "your": true, "yours": true, "ours": true, "human": true, "said": true,
	"asked": true, "brief": true, "names": true, "named": true, "request": true,
	"primary": true, "already": true, "exists": true, "point": true,
	"target": true, "audience": true, "content": true, "surface": true, "surfaces": true,
	"thing": true, "things": true, "work": true, "works": true, "part": true, "core": true,
	"other": true, "another": true, "someone": true, "else": true, "never": true,
	"without": true, "before": true, "after": true, "these": true, "those": true,
	"whole": true, "achieves": true, "result": true, "matters": true,
	"outrank": true, "expected": true, "explicit": true, "explicitly": true,
	"concrete": true, "specific": true, "plain": true, "later": true, "gate": true,
	"drives": true, "carry": true, "look": true, "wired": true, "merely": true,
	"planned": true, "acceptable": true, "define": true, "definition": true,
	"consequence": true, "review": true, "under": true, "yet": true, "attach": true,
	"unit": true, "product": true, "screen": true, "screens": true, "page": true,
	"pages": true, "brand": true,
}

func tokens(s string) []string {
	return wordRE.FindAllString(strings.ToLower(s), -1)
}

func contentWords(s string) map[string]bool {
	out := map[string]bool{}
	for _, w := range tokens(s) {
		if len(w) < 4 || stopwords[w] {
			continue
		}
		out[w] = true
	}
	return out
}

// signalIndex holds, per capability, the distinctive words of its trigger and
// skip prose. Distinctive means "used by few capabilities" — a word that every
// row uses cannot discriminate between rows.
type signalIndex struct {
	trigger map[string]map[string]string // capID -> word -> source condition
	skip    map[string]map[string]string
}

const maxDocFreq = 3

func buildSignalIndex(caps map[string]*resources.CapabilityPhaseDefinition) *signalIndex {
	df := map[string]int{}
	perCapAll := map[string]map[string]bool{}

	for id, c := range caps {
		all := map[string]bool{}
		for _, cond := range append(append([]string{}, c.TriggerConditions...), c.SkipConditions...) {
			for w := range contentWords(cond) {
				all[w] = true
			}
		}
		perCapAll[id] = all
		for w := range all {
			df[w]++
		}
	}

	idx := &signalIndex{
		trigger: map[string]map[string]string{},
		skip:    map[string]map[string]string{},
	}
	for id, c := range caps {
		idx.trigger[id] = map[string]string{}
		idx.skip[id] = map[string]string{}
		for _, cond := range c.TriggerConditions {
			for w := range contentWords(cond) {
				if df[w] <= maxDocFreq {
					if _, seen := idx.trigger[id][w]; !seen {
						idx.trigger[id][w] = cond
					}
				}
			}
		}
		for _, cond := range c.SkipConditions {
			for w := range contentWords(cond) {
				if df[w] <= maxDocFreq {
					if _, seen := idx.skip[id][w]; !seen {
						idx.skip[id][w] = cond
					}
				}
			}
		}
	}
	return idx
}

// ---------- scoring weights ----------

const (
	weightExactTag     = 3.0
	weightSplitTag     = 1.75
	weightTriggerWord  = 1.0
	weightSkipWord     = -1.5
	weightHint         = 6.0
	triggerWordCap     = 4.0
	skipWordFloor      = -4.5
	activationFloor    = 2.0
	ambiguityRatio     = 0.70
	ambiguityAbsMargin = 2.5
)

// ---------- classification ----------

// Classify keeps the original narrow contract for callers that only need a Task.
func (c *Classifier) Classify(rawRequest string) (*Task, error) {
	b := c.ClassifyBrief(rawRequest, Options{})
	return b.ToTask(), nil
}

// ClassifyBrief scores every capability row in the graph against the request
// and returns the structured brief.
func (c *Classifier) ClassifyBrief(rawRequest string, opts Options) *Brief {
	haystack := strings.ToLower(rawRequest)
	if len(opts.ExtraTags) > 0 {
		haystack += " " + strings.ToLower(strings.Join(opts.ExtraTags, " "))
	}
	if opts.ArchetypeHint != "" {
		haystack += " " + strings.ToLower(opts.ArchetypeHint)
	}
	present := contentWords(haystack)

	taskID := opts.TaskID
	if taskID == "" {
		taskID = "task-" + shortHash(rawRequest)
	}

	b := &Brief{
		TaskID:          taskID,
		RawRequest:      rawRequest,
		Type:            detectType(haystack, opts.DeclaredType),
		Platform:        "web",
		HardConstraints: extractConstraints(rawRequest),
	}

	if c.graph == nil || len(c.graph.Capabilities) == 0 {
		b.Archetype = "standard-feature"
		b.ArchetypeReason = "no capability graph loaded; defaulted to standard-feature"
		b.RequiresVisual = b.Type == TypeDesign
		b.RequiresSecurity = b.Type == TypeSecurity
		b.QualityBar = BarStandard
		b.QualityBarReason = "no graph; STANDARD is the safe default"
		b.ResearchDepth = ResearchNone
		b.VerifyDepth = VerifyBuild
		return b
	}

	idx := buildSignalIndex(c.graph.Capabilities)

	var considered []Candidate
	for id, cap := range c.graph.Capabilities {
		cand := Candidate{
			CapabilityID:     id,
			Name:             cap.Name,
			PrimaryArchetype: cap.PrimaryArchetype,
			QualityBar:       cap.QualityBar,
			RiskRank:         cap.RiskRank,
			Platform:         cap.Platform,
		}

		// Strong signal: the operator used the route's own vocabulary.
		for _, tag := range cap.TriggerTags {
			t := strings.ToLower(strings.TrimSpace(tag))
			if t == "" {
				continue
			}
			spaced := strings.ReplaceAll(t, "-", " ")
			if strings.Contains(haystack, t) || strings.Contains(haystack, spaced) {
				cand.Score += weightExactTag
				cand.MatchedTags = append(cand.MatchedTags, tag)
				continue
			}
			parts := strings.Split(spaced, " ")
			if len(parts) > 1 {
				all := true
				for _, p := range parts {
					if len(p) < 3 {
						continue
					}
					if !strings.Contains(haystack, p) {
						all = false
						break
					}
				}
				if all {
					cand.Score += weightSplitTag
					cand.MatchedTags = append(cand.MatchedTags, tag)
				}
			}
		}

		// Weak positive signal: distinctive words from this row's own triggers.
		trigSum := 0.0
		firedTrig := map[string]bool{}
		for w, cond := range idx.trigger[id] {
			if present[w] {
				trigSum += weightTriggerWord
				firedTrig[cond] = true
			}
		}
		if trigSum > triggerWordCap {
			trigSum = triggerWordCap
		}
		cand.Score += trigSum
		cand.FiredTriggers = sortedKeys(firedTrig)

		// Negative signal: the row itself says when not to take it.
		skipSum := 0.0
		firedSkip := map[string]bool{}
		for w, cond := range idx.skip[id] {
			if present[w] {
				skipSum += weightSkipWord
				firedSkip[cond] = true
			}
		}
		if skipSum < skipWordFloor {
			skipSum = skipWordFloor
		}
		cand.Score += skipSum
		cand.FiredSkips = sortedKeys(firedSkip)

		if opts.ArchetypeHint != "" &&
			(strings.EqualFold(opts.ArchetypeHint, id) || strings.EqualFold(opts.ArchetypeHint, cap.PrimaryArchetype)) {
			cand.Score += weightHint
		}

		if cand.Score < activationFloor {
			cand.Declined = true
			if len(cand.FiredSkips) > 0 {
				cand.DeclineReason = "skip condition fired: " + cand.FiredSkips[0]
			} else {
				cand.DeclineReason = fmt.Sprintf("no trigger condition met (score %.2f below floor %.2f)", cand.Score, activationFloor)
			}
		}
		considered = append(considered, cand)
	}

	sort.Slice(considered, func(i, j int) bool {
		if considered[i].Score != considered[j].Score {
			return considered[i].Score > considered[j].Score
		}
		if considered[i].RiskRank != considered[j].RiskRank {
			return considered[i].RiskRank < considered[j].RiskRank
		}
		return considered[i].CapabilityID < considered[j].CapabilityID
	})
	b.Considered = considered

	for _, cand := range considered {
		if !cand.Declined {
			b.Selected = append(b.Selected, cand)
		}
	}

	switch {
	case len(b.Selected) == 0:
		b.Archetype = "standard-feature"
		b.CapabilityID = ""
		b.ArchetypeReason = "no capability row cleared its trigger conditions; treating this as ordinary feature work"
		b.QualityBar = BarStandard
	default:
		top := b.Selected[0]
		b.Archetype = top.PrimaryArchetype
		b.CapabilityID = top.CapabilityID
		b.QualityBar = top.QualityBar
		if top.Platform != "" {
			b.Platform = top.Platform
		}
		reason := fmt.Sprintf("%s scored %.2f", top.CapabilityID, top.Score)
		if len(top.MatchedTags) > 0 {
			reason += " on tags: " + strings.Join(top.MatchedTags, ", ")
		} else if len(top.FiredTriggers) > 0 {
			reason += " on trigger: " + top.FiredTriggers[0]
		}
		b.ArchetypeReason = reason
	}

	// Two routes genuinely fit. Ask once, and only once.
	if len(b.Selected) >= 2 {
		a, bb := b.Selected[0], b.Selected[1]
		// Close either relatively or absolutely is close enough to ask about.
		if bb.Score >= ambiguityRatio*a.Score || (a.Score-bb.Score) < ambiguityAbsMargin {
			b.Ambiguous = true
			b.ClarifyingQuestion = fmt.Sprintf(
				"This reads as both %q and %q. Which is the primary job? "+
					"Answer either and I will treat the other as secondary; if you say nothing I will assume %s, the lower-risk of the two.",
				displayName(a), displayName(bb), lowerRisk(a, bb).CapabilityID)
		}
	}

	b.RequiresVisual = decideVisual(b, haystack)
	b.RequiresSecurity = decideSecurity(b, haystack)
	applyQualityBarOverrides(b, haystack)
	b.ResearchDepth = decideResearchDepth(b)
	b.VerifyDepth = decideVerifyDepth(b)
	applyDesignLab(b, haystack, opts.SkipLab)
	b.UnknownTechnology = detectUnknownTech(haystack, c.graph)

	return b
}

// ResolveSilence applies the no-answer rule: take the lower-risk of the two
// contenders and say out loud that it was assumed.
func (b *Brief) ResolveSilence() {
	if !b.Ambiguous || len(b.Selected) < 2 {
		return
	}
	pick := lowerRisk(b.Selected[0], b.Selected[1])
	b.Archetype = pick.PrimaryArchetype
	b.CapabilityID = pick.CapabilityID
	b.QualityBar = pick.QualityBar
	if pick.Platform != "" {
		b.Platform = pick.Platform
	}
	b.Assumed = true
	b.Ambiguous = false
	b.AssumptionNote = fmt.Sprintf("assumed %s, no response", pick.CapabilityID)
	b.ArchetypeReason = fmt.Sprintf("%s (lower risk_rank %d of the two contenders)", pick.CapabilityID, pick.RiskRank)
	b.ResearchDepth = decideResearchDepth(b)
	b.VerifyDepth = decideVerifyDepth(b)
}

// ToTask projects the brief onto the narrower Task shape.
func (b *Brief) ToTask() *Task {
	var kw []string
	for _, c := range b.Selected {
		kw = append(kw, c.MatchedTags...)
	}
	return &Task{
		ID:                b.TaskID,
		RawRequest:        b.RawRequest,
		Type:              b.Type,
		RequiresVisual:    b.RequiresVisual,
		RequiresSecurity:  b.RequiresSecurity,
		ExtractedKeywords: kw,
	}
}

// ---------- helpers ----------

func displayName(c Candidate) string {
	if c.Name != "" {
		return c.Name
	}
	return c.CapabilityID
}

func lowerRisk(a, b Candidate) Candidate {
	if b.RiskRank != 0 && (a.RiskRank == 0 || b.RiskRank < a.RiskRank) {
		return b
	}
	return a
}

func sortedKeys(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func shortHash(s string) string {
	var h uint32 = 2166136261
	for i := 0; i < len(s); i++ {
		h ^= uint32(s[i])
		h *= 16777619
	}
	return fmt.Sprintf("%08x", h)
}

func detectType(hay, declared string) string {
	if d := strings.ToUpper(strings.TrimSpace(declared)); d != "" {
		return d
	}
	switch {
	case containsAny(hay, "pentest", "vulnerability", "sast", "dast", "security audit", "harden",
		"injection", "xss", "csrf", "owasp", "exploit", "attack surface"):
		return TypeSecurity
	case containsAny(hay, "bug", "broken", "crash", "regression", "not working", "fails", "error",
		"throws", "exception", "stack trace", "traceback", "doesn't work", "does not work",
		"stopped working", "500", "502", "504", "undefined is not"):
		return TypeBugfix
	case containsAny(hay, "refactor", "clean up", "cleanup", "simplify", "rename", "dedupe"):
		return TypeRefactor
	case containsAny(hay, "research", "compare", "evaluate", "investigate", "which library", "options for",
		"citations", "thesis", "manuscript", "methods paper"):
		return TypeResearch
	case containsAny(hay, "redesign", "design", "restyle", "visual", "look and feel"):
		return TypeDesign
	}
	return TypeFeature
}

func containsAny(hay string, needles ...string) bool {
	for _, n := range needles {
		if strings.Contains(hay, n) {
			return true
		}
	}
	return false
}

func decideVisual(b *Brief, hay string) bool {
	if b.Type == TypeDesign {
		return true
	}
	switch b.CapabilityID {
	case "security-audit", "research-paper":
		return false
	case "":
		return containsAny(hay, "ui", "frontend", "css", "layout", "screen", "page", "component", "styling")
	}
	return true
}

func decideSecurity(b *Brief, hay string) bool {
	if b.Type == TypeSecurity || b.CapabilityID == "security-audit" {
		return true
	}
	return containsAny(hay, "auth", "login", "token", "secret", "session", "permission", "sanitiz", "injection")
}

func applyQualityBarOverrides(b *Brief, hay string) {
	if b.QualityBar == "" {
		b.QualityBar = BarStandard
	}
	b.QualityBarReason = "default for " + firstNonEmpty(b.CapabilityID, "standard-feature")

	switch b.Type {
	case TypeBugfix, TypeRefactor:
		b.QualityBar = BarStandard
		b.QualityBarReason = strings.ToLower(b.Type) + " work runs at STANDARD"
		return
	}

	if containsAny(hay, "internal tool", "internal dashboard", "just for me", "just for us",
		"scratch", "throwaway", "prototype", "glue script", "only staff", "internal use") {
		b.QualityBar = BarStandard
		b.QualityBarReason = "request describes internal or throwaway work"
		return
	}

	if containsAny(hay, "client", "recruiter", "public launch", "investor", "customers will see",
		"pitch", "showcase", "portfolio") && b.QualityBar == BarStandard {
		b.QualityBar = BarPremium
		b.QualityBarReason = "a stranger will see this"
	}
}

func decideResearchDepth(b *Brief) string {
	if !b.RequiresVisual {
		if len(b.UnknownTechnology) > 0 || b.Type == TypeResearch {
			return ResearchLight
		}
		return ResearchNone
	}
	switch b.QualityBar {
	case BarExperimental, BarPremium:
		return ResearchDeep
	default:
		return ResearchLight
	}
}

func decideVerifyDepth(b *Brief) string {
	if !b.RequiresVisual {
		return VerifyBuild
	}
	if b.QualityBar == BarStandard {
		return VerifyBrowser
	}
	return VerifyBrowserViewports
}

func applyDesignLab(b *Brief, hay string, optSkip bool) {
	if !b.RequiresVisual {
		b.DesignLabRequired = false
		b.DesignLabReason = "not visual work"
		return
	}
	if optSkip || strings.Contains(hay, "skip the lab") {
		b.DesignLabRequired = false
		b.DesignLabReason = "operator said skip the lab"
		return
	}
	if b.QualityBar == BarStandard {
		b.DesignLabRequired = false
		b.DesignLabReason = "STANDARD bar: lab is off by default, ask to opt in"
		return
	}
	b.DesignLabRequired = true
	b.DesignLabReason = b.QualityBar + " visual work: directions must be approved before frontend files are written"
}

var constraintMarkers = []string{"must ", "must not", "cannot", "can't", "no ", "without ",
	"deadline", "due ", "by tomorrow", "offline", "free tier", "budget", "only use", "do not"}

func extractConstraints(raw string) []string {
	var out []string
	for _, s := range regexp.MustCompile(`[.;\n]`).Split(raw, -1) {
		t := strings.TrimSpace(s)
		if t == "" {
			continue
		}
		l := strings.ToLower(t)
		for _, m := range constraintMarkers {
			if strings.Contains(l, m) {
				out = append(out, t)
				break
			}
		}
	}
	return out
}

// detectUnknownTech surfaces named libraries that no capability row mentions, so
// the loop can research and register them instead of forcing a wrong archetype.
func detectUnknownTech(hay string, g *resources.DesignResourceGraph) []string {
	known := map[string]bool{}
	for _, cap := range g.Capabilities {
		for _, t := range cap.TriggerTags {
			for _, w := range tokens(t) {
				known[w] = true
			}
		}
		for _, w := range tokens(cap.Name + " " + cap.Description) {
			known[w] = true
		}
	}
	for dom, res := range g.Domains {
		known[strings.ToLower(dom)] = true
		for _, r := range res {
			for _, w := range tokens(r) {
				known[w] = true
			}
		}
	}

	var out []string
	seen := map[string]bool{}
	// A capitalised or dotted token in the raw text that nothing in the graph
	// knows about is a candidate unknown library.
	for _, m := range regexp.MustCompile(`\b[a-z][a-z0-9]*(?:\.[a-z0-9]+)+\b|\b[a-z][a-z0-9]{4,}\b`).FindAllString(hay, -1) {
		base := strings.Split(m, ".")[0]
		if known[m] || known[base] || stopwords[m] || seen[m] {
			continue
		}
		if strings.Contains(m, ".") {
			seen[m] = true
			out = append(out, m)
		}
	}
	sort.Strings(out)
	return out
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}
