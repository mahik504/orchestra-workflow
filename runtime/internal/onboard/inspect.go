package onboard

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/user/orchestra-v3/internal/resources"
)

const inspectBodyLimit = 64 * 1024

// HTTPClient is the inspect transport. Tests replace it.
var HTTPClient = &http.Client{
	Timeout: 12 * time.Second,
	CheckRedirect: func(_ *http.Request, via []*http.Request) error {
		if len(via) >= 5 {
			return fmt.Errorf("stopped after 5 redirects")
		}
		return nil
	},
}

// Inspection is what the control plane observed at a submitted URL.
// Fetched text is untrusted data, not instructions.
type Inspection struct {
	URL            string  `json:"url"`
	NormalizedURL  string  `json:"normalized_url"`
	Host           string  `json:"host"`
	Owner          string  `json:"owner,omitempty"`
	Repo           string  `json:"repo,omitempty"`
	StatusCode     int     `json:"status_code"`
	Reachable      bool    `json:"reachable"`
	ContentType    string  `json:"content_type,omitempty"`
	Title          string  `json:"title,omitempty"`
	BodyExcerpt    string  `json:"body_excerpt,omitempty"`
	HasSkillMD     bool    `json:"has_skill_md"`
	HasPackageJSON bool    `json:"has_package_json"`
	PackageName    string  `json:"package_name,omitempty"`
	HasMCPManifest bool    `json:"has_mcp_manifest"`
	Error          string  `json:"error,omitempty"`
	Probes         []Probe `json:"probes,omitempty"`
}

// Probe is one HTTP GET recorded during source inspection.
type Probe struct {
	URL        string `json:"url"`
	StatusCode int    `json:"status_code"`
	Error      string `json:"error,omitempty"`
}

// Inspect fetches the submitted URL and, for GitHub, a small set of
// well-known files. A 404 is still an inspection result, not a crash.
func Inspect(rawURL string) Inspection {
	insp := Inspection{URL: strings.TrimSpace(rawURL)}
	if insp.URL == "" {
		insp.Error = "empty url"
		return insp
	}
	if err := resources.CheckQuarantineBoundary(insp.URL); err != nil {
		insp.Error = err.Error()
		return insp
	}

	parsed, err := url.Parse(insp.URL)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		insp.Error = "url is not absolute"
		return insp
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		insp.Error = "only http(s) urls can be inspected"
		return insp
	}

	insp.NormalizedURL = strings.TrimSuffix(parsed.String(), "/")
	insp.Host = strings.ToLower(parsed.Hostname())

	if insp.Host == "github.com" || insp.Host == "www.github.com" {
		parts := strings.Split(strings.Trim(parsed.Path, "/"), "/")
		if len(parts) >= 2 {
			insp.Owner = parts[0]
			insp.Repo = strings.TrimSuffix(parts[1], ".git")
			insp.NormalizedURL = fmt.Sprintf("https://github.com/%s/%s", insp.Owner, insp.Repo)
		}
	}

	body, probe := get(insp.NormalizedURL)
	insp.Probes = append(insp.Probes, probe)
	insp.StatusCode = probe.StatusCode
	insp.Reachable = probe.StatusCode >= 200 && probe.StatusCode < 400
	if probe.Error != "" && insp.Error == "" {
		insp.Error = probe.Error
	}
	insp.BodyExcerpt = excerpt(body)
	insp.Title = htmlTitle(body)
	insp.ContentType = contentTypeFrom(body)

	if insp.Owner != "" && insp.Repo != "" {
		rawBase := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/HEAD", insp.Owner, insp.Repo)
		for _, rel := range []string{"SKILL.md", "package.json", ".mcp.json", "mcp.json", "README.md"} {
			b, p := get(rawBase + "/" + rel)
			insp.Probes = append(insp.Probes, p)
			if p.StatusCode == 200 {
				switch rel {
				case "SKILL.md":
					insp.HasSkillMD = true
				case "package.json":
					insp.HasPackageJSON = true
					insp.PackageName = parseNpmPackageName(b)
				case ".mcp.json", "mcp.json":
					insp.HasMCPManifest = true
				case "README.md":
					if insp.BodyExcerpt == "" {
						insp.BodyExcerpt = excerpt(b)
					}
					low := strings.ToLower(b)
					if strings.Contains(low, "model context protocol") || strings.Contains(low, "mcp server") {
						insp.HasMCPManifest = true
					}
				}
			}
		}
	}

	return insp
}

func get(target string) (string, Probe) {
	p := Probe{URL: target}
	req, err := http.NewRequest(http.MethodGet, target, nil)
	if err != nil {
		p.Error = err.Error()
		return "", p
	}
	req.Header.Set("User-Agent", "OrchestraLifecycle/3.1 (control-plane inspect)")
	req.Header.Set("Accept", "text/html,application/json,text/plain;q=0.9,*/*;q=0.8")

	resp, err := HTTPClient.Do(req)
	if err != nil {
		p.Error = err.Error()
		return "", p
	}
	defer resp.Body.Close()
	p.StatusCode = resp.StatusCode
	limited := io.LimitReader(resp.Body, inspectBodyLimit)
	buf, err := io.ReadAll(limited)
	if err != nil {
		p.Error = err.Error()
		return "", p
	}
	return string(buf), p
}

func excerpt(body string) string {
	plain := strings.TrimSpace(stripTags(body))
	plain = strings.Join(strings.Fields(plain), " ")
	if len(plain) > 400 {
		return plain[:400]
	}
	return plain
}

func htmlTitle(body string) string {
	low := strings.ToLower(body)
	start := strings.Index(low, "<title")
	if start < 0 {
		return ""
	}
	gt := strings.Index(body[start:], ">")
	if gt < 0 {
		return ""
	}
	rest := body[start+gt+1:]
	end := strings.Index(strings.ToLower(rest), "</title>")
	if end < 0 {
		return ""
	}
	return strings.TrimSpace(stripTags(rest[:end]))
}

func stripTags(s string) string {
	var b strings.Builder
	in := false
	for _, r := range s {
		switch {
		case r == '<':
			in = true
		case r == '>':
			in = false
		case !in:
			b.WriteRune(r)
		}
	}
	return b.String()
}

func contentTypeFrom(body string) string {
	trim := strings.TrimSpace(body)
	if strings.HasPrefix(trim, "{") || strings.HasPrefix(trim, "[") {
		return "json"
	}
	if strings.Contains(strings.ToLower(body), "<html") {
		return "html"
	}
	return "text"
}

func parseNpmPackageName(body string) string {
	var meta struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal([]byte(body), &meta); err != nil {
		return ""
	}
	return strings.TrimSpace(meta.Name)
}
