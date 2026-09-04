package verify

// The Design Lab gate. Its whole job is to make "I'll design it as I code"
// impossible for work a stranger will see.
//
// The gate is a lock on frontend writes, not a printed warning. Until a named
// direction is approved, GuardWrite refuses any file the browser would render.
// Backend files, notes, and the design brief itself stay writable, because the
// human needs something to look at before they can approve anything.

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/user/orchestra-v3/internal/classifier"
)

// GateState is the lifecycle of one Design Lab gate.
type GateState string

const (
	// GateNotRequired means the bar or the work does not call for a lab.
	GateNotRequired GateState = "NOT_REQUIRED"
	// GatePending means directions must be shown and approved before frontend writes.
	GatePending GateState = "PENDING"
	// GateApproved means a named direction was approved by the human.
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

// Direction is one option offered at the gate. Every claim needs a named source;
// a direction with unattributed choices is not a direction, it is a vibe.
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

// Approval records who cleared the gate and for which direction.
type Approval struct {
	TaskID      string    `json:"task_id"`
	DirectionID string    `json:"direction_id"`
	Concept     string    `json:"concept"`
	ApprovedBy  string    `json:"approved_by"`
	ApprovedAt  time.Time `json:"approved_at"`
	Bypass      bool      `json:"bypass"`
	BypassNote  string    `json:"bypass_note,omitempty"`
}

// DesignLab holds the gate state for one task.
type DesignLab struct {
	TaskID        string
	WorkspaceRoot string
	State         GateState
	Reason        string
	Directions    []Direction
	Approved      *Approval
	rejections    []Rejection
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
	return &ErrGateNotCleared{
		Path:   path,
		Reason: firstNonEmpty(d.Reason, "no direction approved yet"),
	}
}

// Cleared reports whether frontend writes are permitted.
func (d *DesignLab) Cleared() bool {
	return d == nil || d.State == GateNotRequired || d.State == GateApproved || d.State == GateBypassed
}

// Offer records the directions shown to the human. The contract asks for two or
// three: one is a decree, four is a survey.
func (d *DesignLab) Offer(dirs []Direction) error {
	if d.State != GatePending {
		return fmt.Errorf("cannot offer directions: gate is %s", d.State)
	}
	if len(dirs) < 2 || len(dirs) > 3 {
		return fmt.Errorf("design lab requires 2 or 3 directions, got %d", len(dirs))
	}
	for _, dir := range dirs {
		if missing := dir.Unsourced(); len(missing) > 0 {
			return fmt.Errorf("direction %q has unattributed claims: %s", dir.ID, strings.Join(missing, ", "))
		}
	}
	// Do not re-offer a combination the human already rejected at this gate.
	prior, err := d.LoadRejections()
	if err != nil {
		return err
	}
	seen := map[string]string{}
	for _, r := range prior {
		seen[r.Fingerprint] = r.Reason
	}
	for _, dir := range dirs {
		if reason, hit := seen[Fingerprint(dir)]; hit {
			return fmt.Errorf("direction %q repeats a combination rejected earlier (%s)", dir.ID, reason)
		}
	}
	d.Directions = dirs
	return nil
}

// Approve clears the gate for a named direction.
func (d *DesignLab) Approve(directionID, approvedBy string) error {
	if d.State == GateNotRequired {
		return nil
	}
	if d.State != GatePending {
		return fmt.Errorf("cannot approve: gate is %s", d.State)
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
	}
	d.State = GateApproved
	d.Reason = "approved: " + chosen.Concept
	return d.persistApproval()
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
