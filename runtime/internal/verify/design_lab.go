package verify

// The Design Lab gate. Its whole job is to make "I'll design it as I code"
// impossible for work a stranger will see.
//
// Two locks for PREMIUM/EXPERIMENTAL visual work (3.3.1):
//  1. Contract — a sourced DESIGN.md (LOCKED vs OPEN). A DESIGN.md is not visual evidence.
//  2. Visual evidence — golden stills the human approved. Product UI stays refused until then.
//
// Backend files, notes, DESIGN.md, REFERENCE_BENCHMARK.md, and files under
// .orchestra/design-lab/golden-stills/ stay writable after the contract is
// approved so a throwaway still can be rendered. Product src/app stays locked.

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/mahik504/orchestra-workflow/runtime/internal/classifier"
)

// GateState is the lifecycle of one Design Lab gate.
type GateState string

const (
	// GateNotRequired means the bar or the work does not call for a lab.
	GateNotRequired GateState = "NOT_REQUIRED"
	// GatePending means translations/survey plus a contract are owed.
	GatePending GateState = "PENDING"
	// GateContractApproved means DESIGN.md LOCKED sections were accepted.
	// Product frontend writes stay refused until golden stills pass.
	GateContractApproved GateState = "CONTRACT_APPROVED"
	// GateApproved means golden stills were human-approved. Product UI may be written.
	GateApproved GateState = "APPROVED"
	// GateBypassed means the human explicitly waived the lab. Recorded, never silent.
	GateBypassed GateState = "BYPASSED"
)

// ErrGateNotCleared is returned when a frontend write is attempted before approval.
type ErrGateNotCleared struct {
	Path   string
	Reason string
}

func (e *ErrGateNotCleared) Error() string {
	return fmt.Sprintf("design lab gate is not cleared: refusing to write %s (%s)", e.Path, e.Reason)
}

// SurveyCardCount is the open-ended cheap survey size. Twenty-three short cards,
// then one full DESIGN.md. Not twenty-three full contracts.
const SurveyCardCount = 23

// NamedReferenceCardCount is the named-reference survey size. Substantial
// references produce three evidence-backed translations, not 23 random cards
// and not a skip straight to one DESIGN.md.
const NamedReferenceCardCount = 3

// DirectionCard is one cheap survey option. It is not a DESIGN.md.
type DirectionCard struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	OneLiner     string `json:"one_liner"`
	Typography   string `json:"typography"`
	ColorWorld   string `json:"color_world"`
	ThreeD       string `json:"three_d"`
	MotionEngine string `json:"motion_engine"`
}

// Incomplete reports missing survey fields.
func (c DirectionCard) Incomplete() []string {
	var missing []string
	if strings.TrimSpace(c.ID) == "" {
		missing = append(missing, "id")
	}
	if strings.TrimSpace(c.Name) == "" {
		missing = append(missing, "name")
	}
	if strings.TrimSpace(c.OneLiner) == "" {
		missing = append(missing, "one_liner")
	}
	if strings.TrimSpace(c.Typography) == "" {
		missing = append(missing, "typography")
	}
	if strings.TrimSpace(c.ColorWorld) == "" {
		missing = append(missing, "color_world")
	}
	if strings.TrimSpace(c.ThreeD) == "" {
		missing = append(missing, "three_d")
	}
	if strings.TrimSpace(c.MotionEngine) == "" {
		missing = append(missing, "motion_engine")
	}
	return missing
}

// Direction is the one full contract after a survey pick, a named override, or
// a pasted DESIGN.md. Every claim needs a named source; a direction with
// unattributed choices is not a direction, it is a vibe.
type Direction struct {
	ID             string   `json:"id"`
	Concept        string   `json:"concept"`
	Typography     string   `json:"typography"`
	TypographySrc  string   `json:"typography_source"`
	ColorWorld     string   `json:"color_world"`
	ColorSrc       string   `json:"color_source"`
	LayoutLanguage string   `json:"layout_language"`
	ComponentKit   string   `json:"component_kit"`
	MotionEngine   string   `json:"motion_engine"`
	MotionWhy      string   `json:"motion_why"`
	ThreeD         string   `json:"three_d"`
	Shader         string   `json:"shader"`
	LogoMethod     string   `json:"logo_method"`
	IconSystem     string   `json:"icon_system"`
	Stack          []string `json:"implementation_stack"`
}

// Unsourced lists the claims in this direction that have no named source.
func (d Direction) Unsourced() []string {
	var missing []string
	if strings.TrimSpace(d.Typography) != "" && strings.TrimSpace(d.TypographySrc) == "" {
		missing = append(missing, "typography")
	}
	if strings.TrimSpace(d.ColorWorld) != "" && strings.TrimSpace(d.ColorSrc) == "" {
		missing = append(missing, "color_world")
	}
	if strings.TrimSpace(d.MotionEngine) != "" && strings.TrimSpace(d.MotionWhy) == "" {
		missing = append(missing, "motion_engine")
	}
	return missing
}

// Rejection records a direction the human turned down and why. The reason is the
// point: it is what stops the next pass from re-offering the same combination.
type Rejection struct {
	TaskID      string    `json:"task_id"`
	DirectionID string    `json:"direction_id"`
	Concept     string    `json:"concept"`
	Fingerprint string    `json:"fingerprint"`
	Reason      string    `json:"reason"`
	RejectedAt  time.Time `json:"rejected_at"`
}

// Approval records who cleared a gate stage and for which direction.
type Approval struct {
	TaskID      string    `json:"task_id"`
	DirectionID string    `json:"direction_id"`
	Concept     string    `json:"concept"`
	ApprovedBy  string    `json:"approved_by"`
	ApprovedAt  time.Time `json:"approved_at"`
	Bypass      bool      `json:"bypass"`
	BypassNote  string    `json:"bypass_note,omitempty"`
	Stage       string    `json:"stage,omitempty"` // "contract" or "stills"
}

// StillsApproval is the visual-evidence pass. Missing files refuse.
type StillsApproval struct {
	TaskID      string    `json:"task_id"`
	ApprovedBy  string    `json:"approved_by"`
	DesktopPath string    `json:"desktop_path"`
	MobilePath  string    `json:"mobile_path"`
	Note        string    `json:"note"`
	ApprovedAt  time.Time `json:"approved_at"`
}

// DesignLab holds the gate state for one task.
type DesignLab struct {
	TaskID              string
	WorkspaceRoot       string
	State               GateState
	Reason              string
	Cards               []DirectionCard
	SurveySkipped       bool
	SkipReason          string
	NamedReferenceMode  bool
	Directions          []Direction
	Approved            *Approval
	Stills              *StillsApproval
	rejections          []Rejection
}

// NewDesignLab derives the gate from the brief. This is the only place that
// decides whether a lab is required, so the rule cannot drift between hosts.
func NewDesignLab(b *classifier.Brief, workspaceRoot string) *DesignLab {
	lab := &DesignLab{TaskID: b.TaskID, WorkspaceRoot: workspaceRoot}
	if b.DesignLabRequired {
		lab.State = GatePending
		lab.Reason = b.DesignLabReason
	} else {
		lab.State = GateNotRequired
		lab.Reason = b.DesignLabReason
	}
	return lab
}

// frontendExts are the files a browser renders. Writing these is what the gate blocks.
var frontendExts = map[string]bool{
	".css": true, ".scss": true, ".sass": true, ".less": true,
	".html": true, ".htm": true, ".vue": true, ".svelte": true,
	".jsx": true, ".tsx": true, ".astro": true,
	".glsl": true, ".frag": true, ".vert": true,
}

// frontendNames catch token and style files that carry no telling extension.
var frontendNames = []string{"tailwind.config", "theme.", "tokens.", "globals.", "design-system."}

// IsFrontendPath reports whether a path is one the gate protects.
func IsFrontendPath(p string) bool {
	base := strings.ToLower(filepath.Base(p))
	if frontendExts[strings.ToLower(filepath.Ext(p))] {
		return true
	}
	for _, n := range frontendNames {
		if strings.HasPrefix(base, n) {
			return true
		}
	}
	return false
}

// GoldenStillsDir is the only product-adjacent path writable after the contract
// and before visual evidence. Throwaway still renderers live here, not in src/.
func (d *DesignLab) GoldenStillsDir() string {
	if d == nil || d.WorkspaceRoot == "" {
		return ""
	}
	return filepath.Join(d.labDir(), "golden-stills")
}

func (d *DesignLab) isGoldenStillPath(p string) bool {
	root := d.GoldenStillsDir()
	if root == "" {
		return false
	}
	if !filepath.IsAbs(p) && d.WorkspaceRoot != "" {
		p = filepath.Join(d.WorkspaceRoot, p)
	}
	abs, err := filepath.Abs(p)
	if err != nil {
		return false
	}
	base, err := filepath.Abs(root)
	if err != nil {
		return false
	}
	rel, err := filepath.Rel(base, abs)
	if err != nil {
		return false
	}
	return rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

// GuardWrite is the lock. Call it before writing any file during implementation.
func (d *DesignLab) GuardWrite(path string) error {
	if d == nil {
		return nil
	}
	switch d.State {
	case GateNotRequired, GateApproved, GateBypassed:
		return nil
	}
	if !IsFrontendPath(path) {
		return nil
	}
	if d.State == GateContractApproved && d.isGoldenStillPath(path) {
		return nil
	}
	reason := firstNonEmpty(d.Reason, "no direction approved yet")
	if d.State == GateContractApproved {
		reason = firstNonEmpty(d.Reason, "DESIGN.md is not visual evidence; golden stills are unpaid")
	}
	return &ErrGateNotCleared{
		Path:   path,
		Reason: reason,
	}
}

// Cleared reports whether product frontend writes are permitted.
func (d *DesignLab) Cleared() bool {
	return d == nil || d.State == GateNotRequired || d.State == GateApproved || d.State == GateBypassed
}

func validateSurveyCards(cards []DirectionCard) error {
	seen := map[string]bool{}
	for _, c := range cards {
		if missing := c.Incomplete(); len(missing) > 0 {
			return fmt.Errorf("survey card %q is incomplete: %s", firstNonEmpty(c.ID, c.Name), strings.Join(missing, ", "))
		}
		if seen[c.ID] {
			return fmt.Errorf("survey card id %q is duplicated", c.ID)
		}
		seen[c.ID] = true
	}
	return nil
}

// OfferSurvey records either 23 open-ended cards or 3 named-reference translations.
func (d *DesignLab) OfferSurvey(cards []DirectionCard) error {
	if d.State != GatePending {
		return fmt.Errorf("cannot offer a survey: gate is %s", d.State)
	}
	if d.SurveySkipped {
		return fmt.Errorf("cannot offer a survey: it was already skipped (%s)", d.SkipReason)
	}
	n := len(cards)
	switch n {
	case SurveyCardCount:
		d.NamedReferenceMode = false
	case NamedReferenceCardCount:
		d.NamedReferenceMode = true
	default:
		return fmt.Errorf("design lab survey requires exactly %d cards (open) or %d translations (named reference), got %d", SurveyCardCount, NamedReferenceCardCount, n)
	}
	if err := validateSurveyCards(cards); err != nil {
		return err
	}
	d.Cards = cards
	return nil
}

// SkipSurvey is for a pasted DESIGN.md or a named skill/MCP/pack that already
// is the contract. Named visual sites go through visual-forensics and three
// evidence-backed translations, not this skip. The reason is required.
func (d *DesignLab) SkipSurvey(reason string) error {
	if d.State != GatePending {
		return fmt.Errorf("cannot skip survey: gate is %s", d.State)
	}
	if strings.TrimSpace(reason) == "" {
		return fmt.Errorf("skipping the survey needs a reason (pasted DESIGN.md, or a named skill/MCP/pack that is the contract)")
	}
	d.SurveySkipped = true
	d.SkipReason = reason
	d.Cards = nil
	return nil
}

func (d *DesignLab) surveySatisfied() bool {
	if d.SurveySkipped {
		return true
	}
	if d.NamedReferenceMode && len(d.Cards) == NamedReferenceCardCount {
		return true
	}
	return len(d.Cards) == SurveyCardCount
}

// OfferContract records the one full sourced direction after a survey pick or a skip.
func (d *DesignLab) OfferContract(dir Direction) error {
	if d.State != GatePending {
		return fmt.Errorf("cannot offer a contract: gate is %s", d.State)
	}
	if !d.surveySatisfied() {
		return fmt.Errorf("offer a %d-card survey, %d named-reference translations, or skip the survey before the contract", SurveyCardCount, NamedReferenceCardCount)
	}
	if missing := dir.Unsourced(); len(missing) > 0 {
		return fmt.Errorf("direction %q has unattributed claims: %s", dir.ID, strings.Join(missing, ", "))
	}
	prior, err := d.LoadRejections()
	if err != nil {
		return err
	}
	for _, r := range prior {
		if r.Fingerprint == Fingerprint(dir) {
			return fmt.Errorf("direction %q repeats a combination rejected earlier (%s)", dir.ID, r.Reason)
		}
	}
	d.Directions = []Direction{dir}
	return nil
}

// ApproveCustom clears the gate for an operator-provided DESIGN.md. No survey.
func (d *DesignLab) ApproveCustom(approvedBy, note string) error {
	if d.State == GateNotRequired {
		return nil
	}
	if d.State != GatePending {
		return fmt.Errorf("cannot approve custom DESIGN.md: gate is %s", d.State)
	}
	if strings.TrimSpace(approvedBy) == "" {
		return fmt.Errorf("approval requires a named approver")
	}
	if strings.TrimSpace(note) == "" {
		return fmt.Errorf("custom DESIGN.md approval needs a note (what was provided)")
	}
	d.SurveySkipped = true
	d.SkipReason = "operator-provided DESIGN.md"
	d.Approved = &Approval{
		TaskID:      d.TaskID,
		DirectionID: "custom",
		Concept:     note,
		ApprovedBy:  approvedBy,
		ApprovedAt:  time.Now().UTC(),
		Stage:       "contract",
	}
	d.State = GateContractApproved
	d.Reason = "contract approved (custom DESIGN.md); golden stills unpaid: " + note
	return d.persistApproval()
}

// ApproveContract accepts LOCKED DESIGN.md. Product UI stays locked.
func (d *DesignLab) ApproveContract(directionID, approvedBy string) error {
	if d.State == GateNotRequired {
		return nil
	}
	if d.State != GatePending {
		return fmt.Errorf("cannot approve contract: gate is %s", d.State)
	}
	var chosen *Direction
	for i := range d.Directions {
		if d.Directions[i].ID == directionID {
			chosen = &d.Directions[i]
			break
		}
	}
	if chosen == nil {
		return fmt.Errorf("no direction %q was offered at this gate", directionID)
	}
	if strings.TrimSpace(approvedBy) == "" {
		return fmt.Errorf("approval requires a named approver")
	}
	d.Approved = &Approval{
		TaskID:      d.TaskID,
		DirectionID: chosen.ID,
		Concept:     chosen.Concept,
		ApprovedBy:  approvedBy,
		ApprovedAt:  time.Now().UTC(),
		Stage:       "contract",
	}
	d.State = GateContractApproved
	d.Reason = "contract approved: " + chosen.Concept + "; DESIGN.md is not visual evidence"
	return d.persistApproval()
}

// Approve is the contract stage. Golden stills remain unpaid. Prefer ApproveContract.
func (d *DesignLab) Approve(directionID, approvedBy string) error {
	return d.ApproveContract(directionID, approvedBy)
}

func (d *DesignLab) resolveExisting(p string) (string, error) {
	p = strings.TrimSpace(p)
	if p == "" {
		return "", fmt.Errorf("still path is empty")
	}
	cands := []string{p}
	if d.WorkspaceRoot != "" && !filepath.IsAbs(p) {
		cands = append(cands, filepath.Join(d.WorkspaceRoot, p), filepath.Join(d.GoldenStillsDir(), p))
	}
	var last error
	for _, c := range cands {
		st, err := os.Stat(c)
		if err != nil {
			last = err
			continue
		}
		if st.IsDir() {
			last = fmt.Errorf("%s is a directory", c)
			continue
		}
		abs, err := filepath.Abs(c)
		if err != nil {
			return "", err
		}
		return abs, nil
	}
	if last != nil {
		return "", last
	}
	return "", fmt.Errorf("still file not found: %s", p)
}

// ApproveStills is the visual-evidence gate. Missing files or an empty note refuse.
func (d *DesignLab) ApproveStills(approvedBy, desktopPath, mobilePath, note string) error {
	if d.State == GateNotRequired {
		return nil
	}
	if d.State != GateContractApproved {
		return fmt.Errorf("cannot approve stills: gate is %s (contract first)", d.State)
	}
	if strings.TrimSpace(approvedBy) == "" {
		return fmt.Errorf("stills approval requires a named approver")
	}
	if strings.TrimSpace(note) == "" {
		return fmt.Errorf("stills approval needs a note (what passed)")
	}
	desk, err := d.resolveExisting(desktopPath)
	if err != nil {
		return fmt.Errorf("desktop still: %w", err)
	}
	mob, err := d.resolveExisting(mobilePath)
	if err != nil {
		return fmt.Errorf("mobile still: %w", err)
	}
	d.Stills = &StillsApproval{
		TaskID:      d.TaskID,
		ApprovedBy:  approvedBy,
		DesktopPath: desk,
		MobilePath:  mob,
		Note:        note,
		ApprovedAt:  time.Now().UTC(),
	}
	d.State = GateApproved
	d.Reason = "golden stills approved: " + note
	return d.persistStills()
}

// Reject records a turned-down direction and its stated reason.
func (d *DesignLab) Reject(directionID, reason string) error {
	if strings.TrimSpace(reason) == "" {
		return fmt.Errorf("a rejection needs a stated reason, otherwise the next pass learns nothing")
	}
	var chosen *Direction
	for i := range d.Directions {
		if d.Directions[i].ID == directionID {
			chosen = &d.Directions[i]
			break
		}
	}
	if chosen == nil {
		return fmt.Errorf("no direction %q was offered at this gate", directionID)
	}
	rej := Rejection{
		TaskID:      d.TaskID,
		DirectionID: chosen.ID,
		Concept:     chosen.Concept,
		Fingerprint: Fingerprint(*chosen),
		Reason:      reason,
		RejectedAt:  time.Now().UTC(),
	}
	d.rejections = append(d.rejections, rej)
	return d.appendRejection(rej)
}

// Bypass waives the lab. It is allowed, but it is written down.
func (d *DesignLab) Bypass(note string) error {
	if strings.TrimSpace(note) == "" {
		return fmt.Errorf("a bypass needs a note saying who waived the lab and why")
	}
	d.State = GateBypassed
	d.Reason = "bypassed: " + note
	d.Approved = &Approval{
		TaskID:     d.TaskID,
		ApprovedAt: time.Now().UTC(),
		Bypass:     true,
		BypassNote: note,
	}
	return d.persistApproval()
}

// Fingerprint identifies a stack combination so a rejected one is recognisable
// even when it comes back wearing a different name.
func Fingerprint(d Direction) string {
	parts := []string{
		norm(d.Typography), norm(d.ColorWorld), norm(d.LayoutLanguage),
		norm(d.ComponentKit), norm(d.MotionEngine), norm(d.ThreeD), norm(d.Shader),
	}
	stack := append([]string{}, d.Stack...)
	sort.Strings(stack)
	for _, s := range stack {
		parts = append(parts, norm(s))
	}
	return strings.Join(parts, "|")
}

func norm(s string) string {
	return strings.Join(strings.Fields(strings.ToLower(strings.TrimSpace(s))), " ")
}

// ---------- persistence ----------

func (d *DesignLab) labDir() string {
	return filepath.Join(d.WorkspaceRoot, ".orchestra", "design-lab")
}

// RejectionLogPath is where turned-down directions accumulate for this workspace.
func (d *DesignLab) RejectionLogPath() string {
	return filepath.Join(d.labDir(), "rejected-directions.json")
}

// ApprovalPath is where the cleared gate is recorded for this task.
func (d *DesignLab) ApprovalPath() string {
	return filepath.Join(d.labDir(), "approved-"+d.TaskID+".json")
}

// LoadRejections reads every direction previously turned down in this workspace.
func (d *DesignLab) LoadRejections() ([]Rejection, error) {
	data, err := os.ReadFile(d.RejectionLogPath())
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var out []Rejection
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, fmt.Errorf("rejection log is corrupt at %s: %w", d.RejectionLogPath(), err)
	}
	return out, nil
}

func (d *DesignLab) appendRejection(r Rejection) error {
	existing, err := d.LoadRejections()
	if err != nil {
		return err
	}
	existing = append(existing, r)
	return writeJSON(d.RejectionLogPath(), existing)
}

func (d *DesignLab) persistApproval() error {
	return writeJSON(d.ApprovalPath(), d.Approved)
}

func (d *DesignLab) persistStills() error {
	if err := d.persistApproval(); err != nil {
		return err
	}
	return writeJSON(filepath.Join(d.labDir(), "stills-"+d.TaskID+".json"), d.Stills)
}

func writeJSON(path string, v any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o644)
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}
